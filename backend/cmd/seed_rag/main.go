// cmd/seed_rag generates neighborhood descriptions from crime data,
// embeds them with Nomic AI, and stores them in the neighborhood_docs table.
//
// Usage:
//
//	go run ./cmd/seed_rag
//	NOMIC_API_KEY=... DATABASE_URL=... go run ./cmd/seed_rag
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"chicago-safety-index/internal/ai"
	"chicago-safety-index/internal/db"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("NOMIC_API_KEY") == "" {
		log.Fatal("NOMIC_API_KEY is not set in .env or environment")
	}

	ctx := context.Background()

	database, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer database.Pool.Close()
	log.Println("Connected to database.")

	// Ensure the pgvector table exists.
	if err := database.EnsureRAGTable(ctx); err != nil {
		log.Fatalf("EnsureRAGTable: %v", err)
	}
	log.Println("neighborhood_docs table ready.")

	// Fetch all (city, zip_code) pairs.
	rows, err := database.Pool.Query(ctx, `
		SELECT city, zip_code FROM zip_codes ORDER BY city, zip_code
	`)
	if err != nil {
		log.Fatalf("fetch zip_codes: %v", err)
	}
	type pair struct{ city, zip string }
	var pairs []pair
	for rows.Next() {
		var p pair
		if err := rows.Scan(&p.city, &p.zip); err != nil {
			log.Fatalf("scan: %v", err)
		}
		pairs = append(pairs, p)
	}
	rows.Close()
	log.Printf("Found %d ZIP codes to process.", len(pairs))

	ok, skipped, failed := 0, 0, 0

	for i, p := range pairs {
		log.Printf("[%d/%d] %s %s ...", i+1, len(pairs), p.city, p.zip)

		desc, err := buildDescription(ctx, database, p.city, p.zip)
		if err != nil {
			log.Printf("  WARN: build description: %v", err)
			skipped++
			continue
		}

		embedding, err := ai.EmbedDocument(ctx, desc)
		if err != nil {
			log.Printf("  WARN: embed: %v", err)
			failed++
			continue
		}

		if err := database.UpsertNeighborhoodDoc(ctx, p.zip, p.city, desc, embedding); err != nil {
			log.Printf("  WARN: upsert: %v", err)
			failed++
			continue
		}

		ok++
		log.Printf("  OK — %d chars, %d-dim embedding", len(desc), len(embedding))

		// Nomic free tier: gentle rate limit.
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("\nDone. success=%d  skipped=%d  failed=%d", ok, skipped, failed)
}

// buildDescription generates a factual neighborhood description from crime stats.
// No LLM call needed — just plain data from the DB.
func buildDescription(ctx context.Context, database *db.Database, city, zip string) (string, error) {
	// Total crimes + year range.
	var totalCrimes int
	var minYear, maxYear int
	err := database.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(crime_count),0), COALESCE(MIN(year),0), COALESCE(MAX(year),0)
		FROM zip_crime_stats
		WHERE city = $1 AND zip_code = $2
	`, city, zip).Scan(&totalCrimes, &minYear, &maxYear)
	if err != nil {
		return "", err
	}
	if totalCrimes == 0 {
		return "", fmt.Errorf("no crime data for %s %s", city, zip)
	}

	// Top 5 crime types.
	typeRows, err := database.Pool.Query(ctx, `
		SELECT crime_type, SUM(crime_count) AS cnt
		FROM zip_crime_stats
		WHERE city = $1 AND zip_code = $2
		GROUP BY crime_type
		ORDER BY cnt DESC
		LIMIT 5
	`, city, zip)
	if err != nil {
		return "", err
	}
	defer typeRows.Close()

	type crimeType struct {
		name  string
		count int
	}
	var topTypes []crimeType
	for typeRows.Next() {
		var ct crimeType
		if err := typeRows.Scan(&ct.name, &ct.count); err != nil {
			return "", err
		}
		topTypes = append(topTypes, ct)
	}

	// Build description.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ZIP code %s in %s. ", zip, city))
	sb.WriteString(fmt.Sprintf("Total recorded crimes: %d between %d and %d. ", totalCrimes, minYear, maxYear))

	if len(topTypes) > 0 {
		sb.WriteString("Most common crime types: ")
		for i, ct := range topTypes {
			pct := float64(ct.count) / float64(totalCrimes) * 100
			sb.WriteString(fmt.Sprintf("%s (%.0f%%)", strings.Title(strings.ToLower(ct.name)), pct))
			if i < len(topTypes)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(". ")
	}

	monthlyAvg := float64(totalCrimes) / float64((maxYear-minYear+1)*12)
	sb.WriteString(fmt.Sprintf("Average monthly crime count: %.1f.", monthlyAvg))

	return sb.String(), nil
}
