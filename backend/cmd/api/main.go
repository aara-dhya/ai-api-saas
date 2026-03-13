package main

import (
	"fmt"
	"net/http"

	"ai-api-saas/pkg/config"
	"ai-api-saas/pkg/database"
)

func main() {
	cfg := config.Load()

	db := database.NewPostgres(cfg.DatabaseURL)

	_ = db // temporary until we use it

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	http.ListenAndServe(":"+cfg.Port, nil)
}
