package ai

import (
	"encoding/json"
	"log"
	"net/http"

	"ai-api-saas/internal/middleware"
	"ai-api-saas/internal/usage"
)

type Handler struct {
	router       *Router
	usageService *usage.Service
}

func NewHandler(router *Router, usageService *usage.Service) *Handler {
	return &Handler{
		router:       router,
		usageService: usageService,
	}
}

type generateRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
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

	// call provider
	resp, err := h.router.Generate(GenerateRequest{
		Prompt: req.Prompt,
		Model:  req.Model,
	})

	if err != nil {
		http.Error(w, "ai provider error", http.StatusInternalServerError)
		return
	}

	// Get API key ID from middleware
	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)

	if ok && resp.Tokens > 0 {

		err := h.usageService.LogUsage(
			apiKeyID,
			resp.Model,
			resp.Tokens,
		)

		if err != nil {
			log.Println("failed to log usage:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Raw)
}
