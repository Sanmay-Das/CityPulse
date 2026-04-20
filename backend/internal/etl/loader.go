// Package etl handles loading Chicago Crimes CSV and ZIP code GeoJSON,
// performing the spatial join in-memory, and writing aggregated results
// to PostgreSQL. Geometry is saved as a static GeoJSON file for the
// frontend — NOT stored in the DB (avoids slow large-text inserts).
package etl

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/tidwall/rtree"

	"chicago-safety-index/internal/db"
)

// ── Spatial index ─────────────────────────────────────────────────────────────

type zipEntry struct {
	zipCode  string
	city     string
	geometry orb.Geometry
}

type spatialIndex struct {
	tree rtree.RTree
	zips []zipEntry
}

func containsPoint(g orb.Geometry, pt orb.Point) bool {
	switch geom := g.(type) {
	case orb.Polygon:
		return planar.PolygonContains(geom, pt)
	case orb.MultiPolygon:
		for _, poly := range geom {
			if planar.PolygonContains(poly, pt) {
				return true
			}
		}
	}
	return false
}

func (idx *spatialIndex) findZip(lon, lat float64) string {
	pt := orb.Point{lon, lat}
	found := ""
	idx.tree.Search(
		[2]float64{lon, lat},
		[2]float64{lon, lat},
		func(_, _ [2]float64, value interface{}) bool {
			i := value.(int)
			if containsPoint(idx.zips[i].geometry, pt) {
				found = idx.zips[i].zipCode
				return false
			}
			return true
		},
	)
	return found
}

// ── Geometry simplification ───────────────────────────────────────────────────

func roundCoord(f float64) float64 { return math.Round(f*1e5) / 1e5 }

func simplifyRing(r orb.Ring) orb.Ring {
	out := make(orb.Ring, len(r))
	for i, p := range r {
		out[i] = orb.Point{roundCoord(p[0]), roundCoord(p[1])}
	}
	return out
}

func simplifyPolygon(poly orb.Polygon) orb.Polygon {
	out := make(orb.Polygon, len(poly))
	for i, ring := range poly {
		out[i] = simplifyRing(ring)
	}
	return out
}

func simplifyGeom(g orb.Geometry) orb.Geometry {
	switch geom := g.(type) {
	case orb.Polygon:
		return simplifyPolygon(geom)
	case orb.MultiPolygon:
		out := make(orb.MultiPolygon, len(geom))
		for i, poly := range geom {
			out[i] = simplifyPolygon(poly)
		}
		return out
	}
	return g
}

// ── ZIP code loader ───────────────────────────────────────────────────────────

// LoadZIPBoundaries reads a GeoJSON FeatureCollection, builds a spatial index,
// and writes a simplified GeoJSON FeatureCollection to outGeoJSONPath for the
// React frontend to load as a static asset. All entries are tagged with city.
func LoadZIPBoundaries(path string, outGeoJSONPath string, city string) (*spatialIndex, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open ZIP GeoJSON: %w", err)
	}
	defer f.Close()

	data, _ := io.ReadAll(f)
	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, nil, fmt.Errorf("parse ZIP GeoJSON: %w", err)
	}

	idx := &spatialIndex{}
	var zipCodes []string
	outFeatures := make([]*geojson.Feature, 0, len(fc.Features))

	for _, feat := range fc.Features {
		zipCode := stringProp(feat, "GEOID10", "ZCTA5CE10", "ZCTA5", "zip")
		if zipCode == "" {
			continue
		}

		g := feat.Geometry
		bound := g.Bound()

		i := len(idx.zips)
		idx.zips = append(idx.zips, zipEntry{zipCode: zipCode, city: city, geometry: g})
		idx.tree.Insert(
			[2]float64{bound.Min[0], bound.Min[1]},
			[2]float64{bound.Max[0], bound.Max[1]},
			i,
		)
		zipCodes = append(zipCodes, zipCode)

		// Build simplified feature for the output GeoJSON file.
		simplified := geojson.NewFeature(simplifyGeom(g))
		simplified.Properties = geojson.Properties{"zipCode": zipCode}
		outFeatures = append(outFeatures, simplified)
	}

	// Write simplified GeoJSON for frontend.
	outFC := geojson.NewFeatureCollection()
	outFC.Features = outFeatures
	outBytes, err := json.Marshal(outFC)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal output GeoJSON: %w", err)
	}
	if err := os.WriteFile(outGeoJSONPath, outBytes, 0644); err != nil {
		return nil, nil, fmt.Errorf("write output GeoJSON: %w", err)
	}
	log.Printf("Wrote simplified GeoJSON → %s (%d KB)", outGeoJSONPath, len(outBytes)/1024)

	log.Printf("Loaded %d ZIP code boundaries for city %q", len(idx.zips), city)
	return idx, zipCodes, nil
}

func stringProp(feat *geojson.Feature, keys ...string) string {
	for _, k := range keys {
		if v, ok := feat.Properties[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// ── CSV format descriptors ────────────────────────────────────────────────────

// CSVFormat describes the column layout and parsing rules for a crimes CSV.
type CSVFormat struct {
	// HasHeader indicates whether the first row is a header (to be skipped).
	HasHeader bool
	// LonCol / LatCol are zero-based column indices for longitude and latitude.
	LonCol, LatCol int
	// DateCol is the zero-based column index for the date/time of the crime.
	DateCol int
	// DateFormats is an ordered list of Go time layouts to try when parsing DateCol.
	DateFormats []string
	// TypeCol is the zero-based column index for the crime-type string.
	TypeCol int
	// YearCol, if >= 0, is a fallback column containing just the year integer.
	// Set to -1 when there is no dedicated year column.
	YearCol int
	// MinCols is the minimum number of columns a valid row must have.
	MinCols int
}

// ChicagoFormat matches the STAR UCR export (no header row).
var ChicagoFormat = CSVFormat{
	HasHeader:   false,
	LonCol:      0,
	LatCol:      1,
	DateCol:     4,
	DateFormats: []string{"01/02/2006 03:04:05 PM", "2006-01-02T15:04:05"},
	TypeCol:     7,
	YearCol:     18,
	MinCols:     19,
}

// LAFormat matches the LA Open Data crimes export (has header row).
// Columns: DR_NO(0), Date Rptd(1), DATE OCC(2), TIME OCC(3), AREA(4),
//
//	AREA NAME(5), Rpt Dist No(6), Part 1-2(7), Crm Cd(8), Crm Cd Desc(9),
//	Mocodes(10), Vict Age(11), Vict Sex(12), Vict Descent(13), Premis Cd(14),
//	Premis Desc(15), Weapon Used Cd(16), Weapon Desc(17), Status(18),
//	Status Desc(19), Crm Cd 1(20), Crm Cd 2(21), Crm Cd 3(22), Crm Cd 4(23),
//	LOCATION(24), Cross Street(25), LAT(26), LON(27)
var LAFormat = CSVFormat{
	HasHeader:   true,
	LonCol:      27,
	LatCol:      26,
	DateCol:     2,
	DateFormats: []string{"01/02/2006 03:04:05 PM", "01/02/2006 15:04:05", "01/02/2006"},
	TypeCol:     9,
	YearCol:     -1,
	MinCols:     28,
}

// ── Crime CSV loader ──────────────────────────────────────────────────────────

type aggKey struct {
	zipCode   string
	year      int
	month     int
	crimeType string
}

// LoadAndJoin reads a Crimes CSV using the supplied format descriptor,
// spatially joins each crime to a ZIP code, aggregates counts, and
// bulk-inserts into PostgreSQL.
func LoadAndJoin(ctx context.Context, database *db.Database, crimesPath string, idx *spatialIndex, zipCodes []string, city string, csvFmt CSVFormat) error {
	// ── 1. Insert ZIP codes (just the codes, no geometry) ─────────────────────
	log.Println("Inserting ZIP codes…")
	if err := insertZIPCodes(ctx, database, zipCodes, city); err != nil {
		return err
	}

	// ── 2. Stream and aggregate crimes CSV ────────────────────────────────────
	log.Println("Reading crimes CSV…")
	f, err := os.Open(crimesPath)
	if err != nil {
		return fmt.Errorf("open crimes CSV: %w", err)
	}
	defer f.Close()

	agg := make(map[aggKey]int)
	reader := csv.NewReader(f)
	reader.LazyQuotes = true // tolerate imperfect quoting in LA data

	// Skip header row if the format has one.
	if csvFmt.HasHeader {
		if _, err := reader.Read(); err != nil {
			return fmt.Errorf("read header: %w", err)
		}
	}

	const logEvery = 100_000
	total, skipped := 0, 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			skipped++
			continue
		}
		if len(record) < csvFmt.MinCols {
			skipped++
			continue
		}

		total++
		if total%logEvery == 0 {
			log.Printf("  processed %d rows…", total)
		}

		latStr := strings.TrimSpace(record[csvFmt.LatCol])
		lonStr := strings.TrimSpace(record[csvFmt.LonCol])
		if latStr == "" || lonStr == "" || latStr == "0" || lonStr == "0" {
			skipped++
			continue
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			skipped++
			continue
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			skipped++
			continue
		}

		zipCode := idx.findZip(lon, lat)
		if zipCode == "" {
			skipped++
			continue
		}

		var t time.Time
		dateStr := strings.TrimSpace(record[csvFmt.DateCol])
		for _, layout := range csvFmt.DateFormats {
			if parsed, e := time.Parse(layout, dateStr); e == nil {
				t = parsed
				break
			}
		}
		if t.IsZero() && csvFmt.YearCol >= 0 && len(record) > csvFmt.YearCol {
			yr, _ := strconv.Atoi(strings.TrimSpace(record[csvFmt.YearCol]))
			if yr > 0 {
				t = time.Date(yr, 1, 1, 0, 0, 0, 0, time.UTC)
			}
		}
		if t.IsZero() {
			skipped++
			continue
		}

		crimeType := "OTHER"
		if ct := strings.TrimSpace(record[csvFmt.TypeCol]); ct != "" {
			crimeType = ct
		}

		key := aggKey{
			zipCode:   zipCode,
			year:      t.Year(),
			month:     int(t.Month()),
			crimeType: crimeType,
		}
		agg[key]++
	}

	log.Printf("Finished reading: %d total, %d skipped, %d buckets", total, skipped, len(agg))

	return insertAggregated(ctx, database, agg, city)
}

// insertZIPCodes deletes existing rows for this city and inserts fresh ones.
func insertZIPCodes(ctx context.Context, database *db.Database, zipCodes []string, city string) error {
	if _, err := database.Pool.Exec(ctx, `DELETE FROM zip_codes WHERE city = $1`, city); err != nil {
		return fmt.Errorf("delete zip_codes for city %q: %w", city, err)
	}

	rows := make([][]any, len(zipCodes))
	for i, z := range zipCodes {
		rows[i] = []any{z, city}
	}

	_, err := database.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"zip_codes"},
		[]string{"zip_code", "city"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("copy zip_codes: %w", err)
	}
	log.Printf("Inserted %d ZIP codes for city %q", len(zipCodes), city)
	return nil
}

func insertAggregated(ctx context.Context, database *db.Database, agg map[aggKey]int, city string) error {
	log.Printf("Bulk-inserting %d aggregation rows…", len(agg))

	if _, err := database.Pool.Exec(ctx, `DELETE FROM zip_crime_stats WHERE city = $1`, city); err != nil {
		return fmt.Errorf("delete zip_crime_stats for city %q: %w", city, err)
	}

	rows := make([][]any, 0, len(agg))
	for k, count := range agg {
		rows = append(rows, []any{k.zipCode, city, k.year, k.month, k.crimeType, count})
	}

	_, err := database.Pool.CopyFrom(
		ctx,
		pgx.Identifier{"zip_crime_stats"},
		[]string{"zip_code", "city", "year", "month", "crime_type", "crime_count"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("bulk insert: %w", err)
	}
	log.Println("ETL complete.")
	return nil
}
