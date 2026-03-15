package system

import (
	"encoding/json"
	"net/http"
)

type Metrics struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	m := Metrics{
		Status:  "ok",
		Service: "ai-api-saas",
	}

	json.NewEncoder(w).Encode(m)
}
