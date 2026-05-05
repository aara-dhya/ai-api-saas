package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"ai-api-saas/internal/ai"
	"ai-api-saas/internal/apikey"
	"ai-api-saas/internal/auth"
	"ai-api-saas/internal/middleware"
	"ai-api-saas/internal/usage"
	"ai-api-saas/pkg/config"
	"ai-api-saas/pkg/database"
)

func main() {

	godotenv.Load()

	cfg := config.Load()

	db := database.NewPostgres(cfg.DatabaseURL)

	// ----------------------
	// SERVICES
	// ----------------------

	apiKeyService := apikey.NewService(db)

	usageService := usage.NewService(db)      // writes
	queryService := usage.NewQueryService(db) // reads
	quotaService := usage.NewQuotaService(db) // quota

	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService)

	// ----------------------
	// HANDLERS
	// ----------------------

	apiKeyHandler := apikey.NewHandler(apiKeyService)
	usageHandler := usage.NewHandler(usageService, queryService)

	// ----------------------
	// MIDDLEWARE
	// ----------------------

	auth := middleware.NewAPIKeyAuth(db)
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)
	quotaMiddleware := middleware.NewQuotaMiddleware(quotaService.CheckQuota)
	// ----------------------
	// AI PROVIDER
	// ----------------------

	groqProvider := ai.NewGroqProvider(cfg.GroqAPIKey)

	router := ai.NewRouter()
	router.Register("llama-3.1-8b-instant", groqProvider)

	aiHandler := ai.NewHandler(router, usageService)

	// ----------------------
	// ROUTES
	// ----------------------

	// public auth routes
	http.HandleFunc("/auth/signup", authHandler.Signup)
	http.HandleFunc("/auth/login", authHandler.Login)

	http.Handle("/api/keys", middleware.AuthMiddleware(http.HandlerFunc(apiKeyHandler.CreateAPIKey)))
	http.Handle("/api/keys/list", middleware.AuthMiddleware(http.HandlerFunc(apiKeyHandler.ListAPIKeys)))
	http.Handle("/api/keys/upgrade", middleware.AuthMiddleware(http.HandlerFunc(apiKeyHandler.UpgradePlan)))

	// health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	// 🔐 API KEY ROUTES (protected)
	http.Handle(
		"/api/keys",
		auth.Middleware(http.HandlerFunc(apiKeyHandler.CreateAPIKey)),
	)

	http.Handle(
		"/api/keys/list",
		auth.Middleware(http.HandlerFunc(apiKeyHandler.ListAPIKeys)),
	)

	http.Handle(
		"/api/keys/upgrade",
		auth.Middleware(http.HandlerFunc(apiKeyHandler.UpgradePlan)),
	)

	// 🔐 USAGE
	http.Handle(
		"/api/usage",
		auth.Middleware(http.HandlerFunc(usageHandler.GetUsage)),
	)

	// 🔐 AI GENERATION
	aiRoute := auth.Middleware(
		rateLimiter.Middleware(
			quotaMiddleware.Middleware(
				http.HandlerFunc(aiHandler.Generate),
			),
		),
	)

	http.Handle("/v1/generate", aiRoute)

	// ----------------------

	fmt.Println("Server running on port", cfg.Port)

	http.ListenAndServe(":"+cfg.Port, nil)
}
