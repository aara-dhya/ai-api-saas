package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"ai-api-saas/internal/apikey"
	"ai-api-saas/internal/middleware"
	"ai-api-saas/pkg/config"
	"ai-api-saas/pkg/database"
)

func main() {

	godotenv.Load()

	cfg := config.Load()

	db := database.NewPostgres(cfg.DatabaseURL)

	// create service
	apiKeyService := apikey.NewService(db)

	// create handler
	apiKeyHandler := apikey.NewHandler(apiKeyService)

	// initialize middleware
	auth := middleware.NewAPIKeyAuth(db)

	// health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	// register API key endpoint (public)
	http.HandleFunc("/api/keys", apiKeyHandler.CreateAPIKey)

	// protected test endpoint
	protected := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(middleware.UserIDKey)

		fmt.Fprintf(w, "Authenticated user: %s", userID)
	}))

	http.Handle("/protected", protected)

	http.ListenAndServe(":"+cfg.Port, nil)
}
