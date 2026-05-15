// Package ai provides a Groq LLM client for natural language query processing.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const groqAPIURL = "https://api.groq.com/openai/v1/chat/completions"
const model = "llama-3.3-70b-versatile"

// dbSchema is sent to Groq so it knows the structure of our database.
const dbSchema = `
Database schema (PostgreSQL):

Table: zip_crime_stats
  - zip_code  VARCHAR  (e.g. "60601")
  - city      VARCHAR  (e.g. "Chicago", "Los Angeles", "San Francisco")
  - year      INT      (e.g. 2023)
  - month     INT      (1-12)
  - crime_type VARCHAR (e.g. "THEFT", "BATTERY", "ROBBERY")
  - crime_count INT    (number of crimes)

Table: zip_codes
  - zip_code VARCHAR
  - city     VARCHAR

Cities available: Chicago, Los Angeles, San Francisco
`

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func callGroq(ctx context.Context, messages []message) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY not set")
	}

	reqBody := groqRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.1,
		MaxTokens:   512,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq API returned status %d", resp.StatusCode)
	}

	var groqResp groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return "", err
	}
	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no response from groq")
	}

	return strings.TrimSpace(groqResp.Choices[0].Message.Content), nil
}

// GenerateSQL converts a natural language question into a SQL query.
func GenerateSQL(ctx context.Context, question string) (string, error) {
	messages := []message{
		{
			Role: "system",
			Content: dbSchema + `
You are a PostgreSQL expert. Convert the user's question into a valid SQL SELECT query.
Rules:
- Only generate SELECT queries, never INSERT/UPDATE/DELETE
- Always use SUM(crime_count) for totals
- Always LIMIT results to 10 unless the user asks for more
- Return ONLY the raw SQL query with no explanation, no markdown, no backticks
- If the question is unclear, generate a reasonable query based on the schema
`,
		},
		{Role: "user", Content: question},
	}
	sql, err := callGroq(ctx, messages)
	if err != nil {
		return "", err
	}
	// Strip markdown code fences if present
	sql = strings.TrimPrefix(sql, "```sql")
	sql = strings.TrimPrefix(sql, "```")
	sql = strings.TrimSuffix(sql, "```")
	return strings.TrimSpace(sql), nil
}

// FormatAnswer converts SQL query results (and optional RAG context) into a natural language answer.
// ragContext may be empty string if no relevant neighborhood docs were found.
func FormatAnswer(ctx context.Context, question, results, ragContext string) (string, error) {
	userContent := fmt.Sprintf("Question: %s\n\nData from database:\n%s", question, results)
	if ragContext != "" {
		userContent += fmt.Sprintf("\n\nNeighborhood context:\n%s", ragContext)
	}
	userContent += "\n\nPlease answer the question based on this data."

	messages := []message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that explains US city crime statistics in plain English. Be concise, clear, and friendly. Provide specific numbers when available.",
		},
		{Role: "user", Content: userContent},
	}
	return callGroq(ctx, messages)
}
