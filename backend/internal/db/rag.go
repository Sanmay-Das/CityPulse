package db

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

// NeighborhoodDoc is one embedded neighborhood description.
type NeighborhoodDoc struct {
	ZipCode     string
	City        string
	Description string
}

// EnsureRAGTable creates the neighborhood_docs table if it doesn't exist.
// Uses TEXT to store embeddings as JSON — no pgvector extension needed.
func (d *Database) EnsureRAGTable(ctx context.Context) error {
	_, err := d.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS neighborhood_docs (
			id          SERIAL PRIMARY KEY,
			zip_code    VARCHAR NOT NULL,
			city        VARCHAR NOT NULL,
			description TEXT    NOT NULL,
			embedding   TEXT    NOT NULL DEFAULT '[]',
			UNIQUE (zip_code, city)
		)
	`)
	if err != nil {
		return fmt.Errorf("create neighborhood_docs: %w", err)
	}
	return nil
}

// UpsertNeighborhoodDoc inserts or updates one neighborhood document.
func (d *Database) UpsertNeighborhoodDoc(ctx context.Context, zipCode, city, description string, embedding []float32) error {
	embJSON, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("marshal embedding: %w", err)
	}
	_, err = d.Pool.Exec(ctx, `
		INSERT INTO neighborhood_docs (zip_code, city, description, embedding)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (zip_code, city) DO UPDATE
			SET description = EXCLUDED.description,
			    embedding   = EXCLUDED.embedding
	`, zipCode, city, description, string(embJSON))
	if err != nil {
		return fmt.Errorf("upsert neighborhood_docs %s/%s: %w", city, zipCode, err)
	}
	return nil
}

// SearchSimilarDocs returns the top-k neighborhood docs closest to the query embedding.
// Cosine similarity is computed in Go after fetching all rows (fast enough for ~600 ZIPs).
func (d *Database) SearchSimilarDocs(ctx context.Context, queryEmbedding []float32, limit int) ([]NeighborhoodDoc, error) {
	rows, err := d.Pool.Query(ctx, `
		SELECT zip_code, city, description, embedding FROM neighborhood_docs
	`)
	if err != nil {
		return nil, fmt.Errorf("SearchSimilarDocs: %w", err)
	}
	defer rows.Close()

	type scored struct {
		doc   NeighborhoodDoc
		score float64
	}
	var results []scored

	for rows.Next() {
		var doc NeighborhoodDoc
		var embJSON string
		if err := rows.Scan(&doc.ZipCode, &doc.City, &doc.Description, &embJSON); err != nil {
			continue
		}
		var emb []float32
		if err := json.Unmarshal([]byte(embJSON), &emb); err != nil {
			continue
		}
		score := cosineSimilarity(queryEmbedding, emb)
		results = append(results, scored{doc, score})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Sort by highest similarity first.
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if limit > len(results) {
		limit = len(results)
	}
	docs := make([]NeighborhoodDoc, limit)
	for i := 0; i < limit; i++ {
		docs[i] = results[i].doc
	}
	return docs, nil
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
