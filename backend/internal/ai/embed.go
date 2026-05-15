package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const nomicAPIURL = "https://api-atlas.nomic.ai/v1/embedding/text"
const nomicModel = "nomic-embed-text-v1.5"

// EmbedDimension is the vector size for nomic-embed-text-v1.5.
const EmbedDimension = 768

type nomicRequest struct {
	Model    string   `json:"model"`
	Texts    []string `json:"texts"`
	TaskType string   `json:"task_type"`
}

type nomicResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

// EmbedDocument embeds text for storage (use when inserting neighborhood docs).
func EmbedDocument(ctx context.Context, text string) ([]float32, error) {
	return embed(ctx, text, "search_document")
}

// EmbedQuery embeds a user query for similarity search.
func EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	return embed(ctx, text, "search_query")
}

func embed(ctx context.Context, text, taskType string) ([]float32, error) {
	apiKey := os.Getenv("NOMIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("NOMIC_API_KEY not set")
	}

	reqBody := nomicRequest{
		Model:    nomicModel,
		Texts:    []string{text},
		TaskType: taskType,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", nomicAPIURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nomic API returned status %d", resp.StatusCode)
	}

	var nomicResp nomicResponse
	if err := json.NewDecoder(resp.Body).Decode(&nomicResp); err != nil {
		return nil, err
	}
	if len(nomicResp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned from nomic")
	}
	return nomicResp.Embeddings[0], nil
}
