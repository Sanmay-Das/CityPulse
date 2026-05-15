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
	"strings"

	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"chicago-safety-index/internal/ai"
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
	// Load .env file for local development (ignored in production).
	_ = godotenv.Load()

	ctx := context.Background()

	// ── Database ──────────────────────────────────────────────────────────────
	database, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer database.Pool.Close()
	log.Println("Connected to database.")

	// ── pgvector RAG table ────────────────────────────────────────────────────
	if err := database.EnsureRAGTable(ctx); err != nil {
		log.Printf("WARN: EnsureRAGTable: %v (RAG will be unavailable)", err)
	} else {
		log.Println("RAG table ready.")
	}

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

	// AI natural language query endpoint.
	mux.HandleFunc("/api/ai/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Question string `json:"question"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Question == "" {
			http.Error(w, "missing question", http.StatusBadRequest)
			return
		}

		// Step 1 (RAG): embed the question and retrieve similar neighborhood docs.
		ragContext := ""
		if queryVec, embErr := ai.EmbedQuery(r.Context(), body.Question); embErr == nil {
			if docs, searchErr := database.SearchSimilarDocs(r.Context(), queryVec, 3); searchErr == nil && len(docs) > 0 {
				var ctxParts []string
				for _, doc := range docs {
					ctxParts = append(ctxParts, doc.Description)
				}
				ragContext = strings.Join(ctxParts, "\n")
			}
		}

		// Step 2: Generate SQL from question.
		sqlQuery, err := ai.GenerateSQL(r.Context(), body.Question)
		if err != nil {
			http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Safety: only allow SELECT queries.
		if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sqlQuery)), "SELECT") {
			http.Error(w, "only SELECT queries allowed", http.StatusBadRequest)
			return
		}

		// Step 3: Execute SQL against database.
		rows, err := database.Pool.Query(r.Context(), sqlQuery)
		if err != nil {
			http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Collect results as a readable string.
		fieldDescs := rows.FieldDescriptions()
		var resultLines []string
		for rows.Next() {
			vals, err := rows.Values()
			if err != nil {
				continue
			}
			var parts []string
			for i, v := range vals {
				parts = append(parts, fmt.Sprintf("%s: %v", fieldDescs[i].Name, v))
			}
			resultLines = append(resultLines, strings.Join(parts, ", "))
		}
		resultStr := strings.Join(resultLines, "\n")
		if resultStr == "" {
			resultStr = "No results found."
		}

		// Step 4: Format results as natural language (with optional RAG context).
		answer, err := ai.FormatAnswer(r.Context(), body.Question, resultStr, ragContext)
		if err != nil {
			http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"answer": answer,
			"sql":    sqlQuery,
		})
	})

	// ── CORS + server ─────────────────────────────────────────────────────────
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	if o := os.Getenv("ALLOWED_ORIGIN"); o != "" {
		allowedOrigins = append(allowedOrigins, o)
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
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
