package middleware

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
)

type contextKey string

const APIKeyIDKey contextKey = "api_key_id"

type APIKeyAuth struct {
	db *sql.DB
}

func NewAPIKeyAuth(db *sql.DB) *APIKeyAuth {
	return &APIKeyAuth{db: db}
}

func (a *APIKeyAuth) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		apiKey := parts[1]

		sum := sha256.Sum256([]byte(apiKey))
		hashedKey := hex.EncodeToString(sum[:])

		var apiKeyID string

		err := a.db.QueryRow(
			"SELECT id FROM api_keys WHERE key_hash = $1 AND revoked = false",
			hashedKey,
		).Scan(&apiKeyID)

		if err != nil {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), APIKeyIDKey, apiKeyID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
