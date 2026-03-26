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

	// services
	apiKeyService := apikey.NewService(db)
	usageService := usage.NewService(db)
	// quotaService := usage.NewQuotaService(db)

	// handlers
	apiKeyHandler := apikey.NewHandler(apiKeyService)
	usageHandler := usage.NewHandler(usageService)

	// middleware
	auth := middleware.NewAPIKeyAuth(db)
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)
	quotaService := usage.NewQuotaService(db)
	quotaMiddleware := middleware.NewQuotaMiddleware(quotaService.CheckQuota)

	// AI provider
	groqProvider := ai.NewGroqProvider(cfg.GroqAPIKey)

	// AI router
	router := ai.NewRouter()
	router.Register("llama-3.1-8b-instant", groqProvider)

	// AI handler
	aiHandler := ai.NewHandler(router, usageService)

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	// public route
	http.HandleFunc("/api/keys", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			apiKeyHandler.CreateAPIKey(w, r)
		case http.MethodGet:
			apiKeyHandler.ListAPIKeys(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// usage endpoint (protected)
	http.Handle(
		"/api/usage",
		auth.Middleware(
			http.HandlerFunc(usageHandler.GetUsage),
		),
	)

	http.Handle(
		"/api/keys/upgrade",
		auth.Middleware(http.HandlerFunc(apiKeyHandler.UpgradePlan)),
	)

	// AI generation endpoint (protected)
	aiRoute := auth.Middleware(
		rateLimiter.Middleware(
			quotaMiddleware.Middleware(
				http.HandlerFunc(aiHandler.Generate),
			),
		),
	)

	http.Handle("/v1/generate", aiRoute)

	fmt.Println("Server running on port", cfg.Port)

	http.ListenAndServe(":"+cfg.Port, nil)
}
