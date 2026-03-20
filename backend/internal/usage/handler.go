package usage

import (
	"encoding/json"
	"net/http"

	"ai-api-saas/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetUsage(w http.ResponseWriter, r *http.Request) {

	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)
	if !ok {
		http.Error(w, "missing api key", http.StatusUnauthorized)
		return
	}

	summary, err := h.service.UsageSummary(apiKeyID)
	if err != nil {
		http.Error(w, "failed to fetch usage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
