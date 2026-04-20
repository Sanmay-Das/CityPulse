// cmd/etl runs the one-time ETL pipeline:
//  1. Loads ZIP code GeoJSON → builds spatial index
//  2. Writes simplified GeoJSON to frontend/public/ (static asset, per-city file)
//  3. Loads Crimes CSV → spatial join → aggregate
//  4. Bulk-inserts aggregated counts into PostgreSQL
//
// Usage (Chicago):
//
//	go run ./cmd/etl \
//	  -crimes  ../data/Chicago_Crimes.csv \
//	  -zipgeo  ../data/TIGER2018_ZCTA5.geojson \
//	  -city    Chicago \
//	  -format  chicago
//
// Usage (Los Angeles):
//
//	go run ./cmd/etl \
//	  -crimes  ../data/LA_Crimes.csv \
//	  -zipgeo  ../data/TIGER2018_ZCTA5_LA.geojson \
//	  -city    "Los Angeles" \
//	  -format  la
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"chicago-safety-index/internal/db"
	"chicago-safety-index/internal/etl"
)

// cityToFilename converts a city name to a lowercase, underscore-separated
// filename fragment (e.g. "Los Angeles" → "los_angeles").
func cityToFilename(city string) string {
	return strings.ToLower(strings.ReplaceAll(city, " ", "_"))
}

func main() {
	crimesPath := flag.String("crimes", "", "Path to Crimes CSV (required)")
	zipPath := flag.String("zipgeo", "", "Path to ZIP code boundaries GeoJSON (required)")
	outPath := flag.String("out", "", "Output path for simplified GeoJSON (auto-derived from city if empty)")
	city := flag.String("city", "Chicago", "City name to tag loaded data with (e.g. Chicago, \"Los Angeles\")")
	format := flag.String("format", "chicago", "CSV format: chicago | la")
	flag.Parse()

	if *crimesPath == "" || *zipPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Derive per-city output path when not specified.
	if *outPath == "" {
		*outPath = fmt.Sprintf("../frontend/public/zip_boundaries_%s.geojson", cityToFilename(*city))
	}

	var csvFmt etl.CSVFormat
	switch strings.ToLower(*format) {
	case "chicago":
		csvFmt = etl.ChicagoFormat
	case "la", "los_angeles", "losangeles":
		csvFmt = etl.LAFormat
	default:
		log.Fatalf("Unknown -format %q. Valid values: chicago, la", *format)
	}

	ctx := context.Background()

	database, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer database.Pool.Close()

	log.Printf("Loading ZIP code boundaries for city %q → %s", *city, *outPath)
	idx, zipCodes, err := etl.LoadZIPBoundaries(*zipPath, *outPath, *city)
	if err != nil {
		log.Fatalf("LoadZIPBoundaries: %v", err)
	}

	if err := etl.LoadAndJoin(ctx, database, *crimesPath, idx, zipCodes, *city, csvFmt); err != nil {
		log.Fatalf("LoadAndJoin: %v", err)
	}
}
