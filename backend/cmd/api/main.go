package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"ai-api-saas/internal/ai"
	"ai-api-saas/internal/apikey"
	"ai-api-saas/internal/middleware"
	"ai-api-saas/pkg/config"
	"ai-api-saas/pkg/database"
)

func main() {

	godotenv.Load()

	cfg := config.Load()

	db := database.NewPostgres(cfg.DatabaseURL)

	apiKeyService := apikey.NewService(db)
	apiKeyHandler := apikey.NewHandler(apiKeyService)

	auth := middleware.NewAPIKeyAuth(db)

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	// public route
	http.HandleFunc("/api/keys", apiKeyHandler.CreateAPIKey)

	// AI handler
	aiHandler := ai.NewHandler(cfg.GroqAPIKey)

	// protected route
	aiRoute := auth.Middleware(http.HandlerFunc(aiHandler.Generate))

	http.Handle("/ai/generate", aiRoute)

	http.ListenAndServe(":"+cfg.Port, nil)
}
