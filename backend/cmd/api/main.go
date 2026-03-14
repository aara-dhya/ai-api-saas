package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"ai-api-saas/internal/ai"
	"ai-api-saas/internal/apikey"
	"ai-api-saas/internal/middleware"
	"ai-api-saas/internal/usage"
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

	rateLimiter := middleware.NewRateLimiter(10, time.Minute)

	usageService := usage.NewService(db)

	// create provider
	groqProvider := ai.NewGroqProvider(cfg.GroqAPIKey)

	// create router
	router := ai.NewRouter()

	// register models
	router.Register("llama-3.1-8b-instant", groqProvider)

	// create handler
	aiHandler := ai.NewHandler(router, usageService)

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	// public route
	http.HandleFunc("/api/keys", apiKeyHandler.CreateAPIKey)

	// protected route
	aiRoute := auth.Middleware(
		rateLimiter.Middleware(
			http.HandlerFunc(aiHandler.Generate),
		),
	)

	http.Handle("/ai/generate", aiRoute)

	http.ListenAndServe(":"+cfg.Port, nil)
}
