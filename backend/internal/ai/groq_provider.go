package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GroqProvider struct {
	apiKey string
	client *http.Client
}

func NewGroqProvider(key string) *GroqProvider {
	return &GroqProvider{
		apiKey: key,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type groqUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type groqResp struct {
	Model string    `json:"model"`
	Usage groqUsage `json:"usage"`
}

func (g *GroqProvider) Generate(req GenerateRequest) (*GenerateResponse, error) {

	if req.Model == "" {
		req.Model = "llama-3.1-8b-instant"
	}

	groqReq := groqRequest{
		Model: req.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
	}

	body, err := json.Marshal(groqReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// attach context with timeout (extra safety)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// handle non-200 properly
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"groq API error: status=%d body=%s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var parsed groqResp
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &GenerateResponse{
		Model:  parsed.Model,
		Tokens: parsed.Usage.TotalTokens,
		Raw:    responseBody,
	}, nil
}
