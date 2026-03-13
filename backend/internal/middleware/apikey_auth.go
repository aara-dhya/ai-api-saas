package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

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

		var userID string

		err := a.db.QueryRow(
			"SELECT user_id FROM api_keys WHERE key = $1",
			apiKey,
		).Scan(&userID)

		if err != nil {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
