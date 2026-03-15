package usage

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type usageResponse struct {
	Model     string    `json:"model"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) GetUsage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKeyID := r.Header.Get("X-API-Key-ID")

	if apiKeyID == "" {
		http.Error(w, "missing api key", http.StatusUnauthorized)
		return
	}

	rows, err := h.service.db.Query(
		`SELECT model, tokens, created_at
		 FROM usage_logs
		 WHERE api_key_id = $1
		 ORDER BY created_at DESC
		 LIMIT 100`,
		apiKeyID,
	)

	if err != nil {
		http.Error(w, "failed to fetch usage", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var results []usageResponse

	for rows.Next() {

		var u usageResponse

		err := rows.Scan(&u.Model, &u.Tokens, &u.CreatedAt)
		if err != nil {
			http.Error(w, "failed to parse usage", http.StatusInternalServerError)
			return
		}

		results = append(results, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
