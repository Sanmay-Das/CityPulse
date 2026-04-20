// cmd/server starts the HTTP server exposing:
//   - POST /graphql   — GraphQL API (Apollo-compatible)
//   - GET  /graphql   — GraphiQL playground (dev only)
//   - GET  /api/export/csv  — CSV export for PowerBI (?city=Chicago&year=2023)
//   - GET  /api/export/json — JSON export              (?city=Chicago&year=2023)
//   - GET  /api/health      — health check
//
// Usage:
//
//	go run ./cmd/server
//	DATABASE_URL=postgres://... PORT=8080 go run ./cmd/server
package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/graphql-go/handler"
	"github.com/rs/cors"

	"chicago-safety-index/internal/db"
	"chicago-safety-index/internal/graph"
)

// cityParam extracts the ?city= query parameter, defaulting to "Chicago".
func cityParam(r *http.Request) string {
	if c := r.URL.Query().Get("city"); c != "" {
		return c
	}
	return "Chicago"
}

func main() {
	ctx := context.Background()

	// ── Database ──────────────────────────────────────────────────────────────
	database, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer database.Pool.Close()
	log.Println("Connected to database.")

	// ── GraphQL schema ────────────────────────────────────────────────────────
	resolver := &graph.Resolver{DB: database}
	schema, err := graph.NewSchema(resolver)
	if err != nil {
		log.Fatalf("Build GraphQL schema: %v", err)
	}

	// ── HTTP handlers ─────────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// GraphQL endpoint (GraphiQL enabled for local dev).
	gqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	mux.Handle("/graphql", gqlHandler)

	// Health check.
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	// CSV export for PowerBI.
	mux.HandleFunc("/api/export/csv", func(w http.ResponseWriter, r *http.Request) {
		city := cityParam(r)
		year := 0
		if ys := r.URL.Query().Get("year"); ys != "" {
			year, _ = strconv.Atoi(ys)
		}

		rows, err := database.GetExportData(r.Context(), city, year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="crimes_export.csv"`)

		cw := csv.NewWriter(w)
		_ = cw.Write([]string{"zip_code", "year", "month", "crime_type", "crime_count"})
		for _, row := range rows {
			_ = cw.Write([]string{
				fmt.Sprint(row["zip_code"]),
				fmt.Sprint(row["year"]),
				fmt.Sprint(row["month"]),
				fmt.Sprint(row["crime_type"]),
				fmt.Sprint(row["crime_count"]),
			})
		}
		cw.Flush()
	})

	// JSON export for PowerBI (alternative to CSV).
	mux.HandleFunc("/api/export/json", func(w http.ResponseWriter, r *http.Request) {
		city := cityParam(r)
		year := 0
		if ys := r.URL.Query().Get("year"); ys != "" {
			year, _ = strconv.Atoi(ys)
		}

		rows, err := database.GetExportData(r.Context(), city, year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rows)
	})

	// ── CORS + server ─────────────────────────────────────────────────────────
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}).Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("Server listening on http://localhost%s", addr)
	log.Printf("GraphiQL playground: http://localhost%s/graphql", addr)
	log.Printf("CSV export:          http://localhost%s/api/export/csv?city=Chicago", addr)

	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
