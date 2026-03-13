package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"ai-api-saas/internal/apikey"
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

	// health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	// register API key endpoint
	http.HandleFunc("/api/keys", apiKeyHandler.CreateAPIKey)

	http.ListenAndServe(":"+cfg.Port, nil)
}
