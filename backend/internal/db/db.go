package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context) (*Database, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://chicago:chicago@localhost:5432/chicago_safety"
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return &Database{Pool: pool}, nil
}

// ── Types ─────────────────────────────────────────────────────────────────────

type ZipStats struct {
	City         string
	ZipCode      string
	TotalCrimes  int
	TopCrimeType string
}

type MonthlyCount struct {
	Year  int
	Month int
	Count int
}

type CrimeTypeCount struct {
	CrimeType string
	Count     int
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (d *Database) GetAllZipStats(ctx context.Context, city string, year int) ([]ZipStats, error) {
	// args always starts with city ($1); year ($2) is added conditionally.
	args := []any{city}

	yearMainFilter := "" // applied on the outer LEFT JOIN e.g. "AND s.year = $2"
	yearSubFilter := ""  // applied inside the subquery    e.g. "AND s2.year = $2"
	if year > 0 {
		yearMainFilter = "AND s.year = $2"
		yearSubFilter = "AND s2.year = $2"
		args = append(args, year)
	}

	query := fmt.Sprintf(`
		SELECT
			z.city,
			z.zip_code,
			COALESCE(SUM(s.crime_count), 0) AS total_crimes,
			COALESCE(
				(SELECT crime_type
				 FROM zip_crime_stats s2
				 WHERE s2.zip_code = z.zip_code
				   AND s2.city = z.city
				   %s
				 GROUP BY crime_type
				 ORDER BY SUM(crime_count) DESC
				 LIMIT 1
				), 'N/A') AS top_crime_type
		FROM zip_codes z
		LEFT JOIN zip_crime_stats s
			ON s.zip_code = z.zip_code
			AND s.city = z.city
			%s
		WHERE z.city = $1
		GROUP BY z.city, z.zip_code
		ORDER BY total_crimes DESC
	`, yearSubFilter, yearMainFilter)

	pgRows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAllZipStats: %w", err)
	}
	defer pgRows.Close()

	var rows []ZipStats
	for pgRows.Next() {
		var r ZipStats
		if err := pgRows.Scan(&r.City, &r.ZipCode, &r.TotalCrimes, &r.TopCrimeType); err != nil {
			return nil, err
		}
		rows = append(rows, r)
	}
	return rows, pgRows.Err()
}

func (d *Database) GetZipStats(ctx context.Context, city string, zipCode string) (*ZipStats, error) {
	query := `
		SELECT
			z.city,
			z.zip_code,
			COALESCE(SUM(s.crime_count), 0) AS total_crimes,
			COALESCE(
				(SELECT crime_type
				 FROM zip_crime_stats s2
				 WHERE s2.zip_code = z.zip_code
				   AND s2.city = z.city
				 GROUP BY s2.crime_type
				 ORDER BY SUM(s2.crime_count) DESC
				 LIMIT 1
				), 'N/A') AS top_crime_type
		FROM zip_codes z
		LEFT JOIN zip_crime_stats s
			ON s.zip_code = z.zip_code
			AND s.city = z.city
		WHERE z.city = $1
		  AND z.zip_code = $2
		GROUP BY z.city, z.zip_code
	`
	var r ZipStats
	err := d.Pool.QueryRow(ctx, query, city, zipCode).Scan(&r.City, &r.ZipCode, &r.TotalCrimes, &r.TopCrimeType)
	if err != nil {
		return nil, fmt.Errorf("GetZipStats %s/%s: %w", city, zipCode, err)
	}
	return &r, nil
}

func (d *Database) GetCrimeTimeline(ctx context.Context, city string, zipCode string) ([]MonthlyCount, error) {
	pgRows, err := d.Pool.Query(ctx, `
		SELECT year, month, SUM(crime_count) AS count
		FROM zip_crime_stats
		WHERE city = $1
		  AND zip_code = $2
		GROUP BY year, month
		ORDER BY year, month
	`, city, zipCode)
	if err != nil {
		return nil, fmt.Errorf("GetCrimeTimeline: %w", err)
	}
	defer pgRows.Close()

	var result []MonthlyCount
	for pgRows.Next() {
		var m MonthlyCount
		if err := pgRows.Scan(&m.Year, &m.Month, &m.Count); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, pgRows.Err()
}

func (d *Database) GetCrimeTypes(ctx context.Context, city string, zipCode string, year int) ([]CrimeTypeCount, error) {
	args := []any{city, zipCode}
	yearFilter := ""
	if year > 0 {
		yearFilter = "AND year = $3"
		args = append(args, year)
	}

	query := fmt.Sprintf(`
		SELECT crime_type, SUM(crime_count) AS count
		FROM zip_crime_stats
		WHERE city = $1
		  AND zip_code = $2 %s
		GROUP BY crime_type
		ORDER BY count DESC
		LIMIT 15
	`, yearFilter)

	pgRows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetCrimeTypes: %w", err)
	}
	defer pgRows.Close()

	var result []CrimeTypeCount
	for pgRows.Next() {
		var c CrimeTypeCount
		if err := pgRows.Scan(&c.CrimeType, &c.Count); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, pgRows.Err()
}

func (d *Database) GetAvailableYears(ctx context.Context, city string) ([]int, error) {
	pgRows, err := d.Pool.Query(ctx, `SELECT DISTINCT year FROM zip_crime_stats WHERE city = $1 ORDER BY year`, city)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableYears: %w", err)
	}
	defer pgRows.Close()

	var years []int
	for pgRows.Next() {
		var y int
		if err := pgRows.Scan(&y); err != nil {
			return nil, err
		}
		years = append(years, y)
	}
	return years, pgRows.Err()
}

type YearRange struct {
	MinYear int
	MaxYear int
}

func (d *Database) GetYearRange(ctx context.Context, city string) (*YearRange, error) {
	var r YearRange
	err := d.Pool.QueryRow(ctx, `
		SELECT COALESCE(MIN(year), 0), COALESCE(MAX(year), 0)
		FROM zip_crime_stats
		WHERE city = $1
	`, city).Scan(&r.MinYear, &r.MaxYear)
	if err != nil {
		return nil, fmt.Errorf("GetYearRange: %w", err)
	}
	return &r, nil
}

func (d *Database) GetAvailableCities(ctx context.Context) ([]string, error) {
	pgRows, err := d.Pool.Query(ctx, `SELECT DISTINCT city FROM zip_codes ORDER BY city`)
	if err != nil {
		return nil, fmt.Errorf("GetAvailableCities: %w", err)
	}
	defer pgRows.Close()

	var cities []string
	for pgRows.Next() {
		var c string
		if err := pgRows.Scan(&c); err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}
	return cities, pgRows.Err()
}

func (d *Database) GetExportData(ctx context.Context, city string, year int) ([]map[string]any, error) {
	args := []any{city}
	yearFilter := ""
	if year > 0 {
		yearFilter = "AND year = $2"
		args = append(args, year)
	}

	query := fmt.Sprintf(`
		SELECT zip_code, year, month, crime_type, crime_count
		FROM zip_crime_stats
		WHERE city = $1 %s
		ORDER BY zip_code, year, month, crime_type
	`, yearFilter)

	pgRows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetExportData: %w", err)
	}
	defer pgRows.Close()

	var rows []map[string]any
	for pgRows.Next() {
		var zipCode, crimeType string
		var yr, month, count int
		if err := pgRows.Scan(&zipCode, &yr, &month, &crimeType, &count); err != nil {
			return nil, err
		}
		rows = append(rows, map[string]any{
			"zip_code":    zipCode,
			"year":        yr,
			"month":       month,
			"crime_type":  crimeType,
			"crime_count": count,
		})
	}
	return rows, pgRows.Err()
}
