package usage

import (
	"encoding/json"
	"net/http"

	"ai-api-saas/internal/middleware"
)

type Handler struct {
	service      *Service
	queryService *QueryService
}

func NewHandler(s *Service, qs *QueryService) *Handler {
	return &Handler{
		service:      s,
		queryService: qs,
	}
}

func (h *Handler) GetUsage(w http.ResponseWriter, r *http.Request) {

	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)
	if !ok {
		http.Error(w, "missing api key", http.StatusUnauthorized)
		return
	}

	summary, err := h.queryService.UsageSummary(apiKeyID)
	if err != nil {
		http.Error(w, "failed to fetch usage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
