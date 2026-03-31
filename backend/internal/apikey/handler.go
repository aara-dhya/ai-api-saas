package apikey

import (
	"ai-api-saas/internal/middleware"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createKeyRequest struct {
	Name string `json:"name"`
}

type createKeyResponse struct {
	APIKey string `json:"api_key"`
}

func (h *Handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)
	if !ok {
		http.Error(w, "missing user identity", http.StatusUnauthorized)
		return
	}

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

	key, err := h.service.CreateAPIKey(userID, req.Name)
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

	// ✅ validate plan exists
	var planID int
	err = h.service.db.QueryRow(
		`SELECT id FROM plans WHERE name = $1`,
		req.Plan,
	).Scan(&planID)

	if err != nil {
		http.Error(w, "invalid plan", http.StatusBadRequest)
		return
	}

	// ✅ update safely
	_, err = h.service.db.Exec(
		`UPDATE api_keys SET plan_id = $1 WHERE id = $2`,
		planID,
		apiKeyID,
	)

	if err != nil {
		http.Error(w, "failed to upgrade plan", http.StatusInternalServerError)
		return
	}

	// ✅ proper response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "plan updated",
	})
}

func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {

	apiKeyID, ok := r.Context().Value(middleware.APIKeyIDKey).(string)
	if !ok {
		http.Error(w, "missing api key", http.StatusUnauthorized)
		return
	}

	rows, err := h.service.db.Query(
		`SELECT ak.id, ak.key, p.name
		FROM api_keys ak
		JOIN plans p ON ak.plan_id = p.id
		WHERE ak.user_id = (
			SELECT user_id FROM api_keys WHERE id = $1
		)`,
		apiKeyID,
	)

	if err != nil {
		http.Error(w, "failed to fetch api keys", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type response struct {
		ID   string `json:"id"`
		Key  string `json:"key"`
		Plan string `json:"plan"`
	}

	var result []response

	for rows.Next() {
		var r response
		if err := rows.Scan(&r.ID, &r.Key, &r.Plan); err != nil {
			http.Error(w, "failed to parse data", http.StatusInternalServerError)
			return
		}

		// 🔥 MASK THE KEY
		r.Key = maskAPIKey(r.Key)

		result = append(result, r)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "failed during iteration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
