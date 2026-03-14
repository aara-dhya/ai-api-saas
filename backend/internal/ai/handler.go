package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"ai-api-saas/internal/middleware"
	"ai-api-saas/internal/usage"
)

type Handler struct {
	apiKey       string
	usageService *usage.Service
}

func NewHandler(key string, usageService *usage.Service) *Handler {
	return &Handler{
		apiKey:       key,
		usageService: usageService,
	}
}

type generateRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type groqResponse struct {
	Model string `json:"model"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req generateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "failed to encode request", http.StatusInternalServerError)
		return
	}

	httpReq, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "failed to call ai provider", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read provider response", http.StatusInternalServerError)
		return
	}

	// Parse Groq response to extract usage
	var parsed groqResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		log.Println("failed to parse groq response:", err)
	}

	// Get API key ID from middleware context
	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)
	if ok && parsed.Usage.TotalTokens > 0 {

		err := h.usageService.LogUsage(
			apiKeyID,
			parsed.Model,
			parsed.Usage.TotalTokens,
		)

		if err != nil {
			log.Println("failed to log usage:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
