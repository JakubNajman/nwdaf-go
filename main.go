package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type V1HealthResponse struct {
	Status        string `json:"status"`
	NRF           string `json:"nrf"`
	ADRF          string `json:"adrf"`
	Subscriptions int    `json:"subscriptions"`
}

func V1Health(w http.ResponseWriter, r *http.Request) {
	resp := V1HealthResponse{
		Status:        "ok",
		NRF:           "mockNRF",
		ADRF:          "mockADRF",
		Subscriptions: 1,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /nnwdaf-analyticsinfo/v1/health", V1Health)

	log.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
