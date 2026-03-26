package apikey

import (
	"ai-api-saas/internal/middleware"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createKeyRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type createKeyResponse struct {
	APIKey string `json:"api_key"`
}

func (h *Handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createKeyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	key, err := h.service.CreateAPIKey(req.UserID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := createKeyResponse{
		APIKey: key,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) UpgradePlan(w http.ResponseWriter, r *http.Request) {

	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)
	if !ok {
		http.Error(w, "missing api key", http.StatusUnauthorized)
		return
	}

	var req struct {
		Plan string `json:"plan"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Plan == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err = h.service.db.Exec(
		`UPDATE api_keys
		 SET plan_id = (SELECT id FROM plans WHERE name = $1)
		 WHERE id = $2`,
		req.Plan,
		apiKeyID,
	)

	if err != nil {
		http.Error(w, "failed to upgrade plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("plan updated"))
}
