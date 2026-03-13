package main

import (
	"fmt"
	"net/http"

	"ai-api-saas/pkg/config"
)

func main() {
	cfg := config.Load()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API running on port", cfg.Port)
	})

	http.ListenAndServe(":"+cfg.Port, nil)
}
