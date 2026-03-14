package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type GroqProvider struct {
	apiKey string
}

func NewGroqProvider(key string) *GroqProvider {
	return &GroqProvider{
		apiKey: key,
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
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("groq api error")
	}

	var parsed groqResp
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, err
	}

	return &GenerateResponse{
		Model:  parsed.Model,
		Tokens: parsed.Usage.TotalTokens,
		Raw:    responseBody,
	}, nil
}
