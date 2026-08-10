package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func newV1AnalyticsRequest() V1AnalyticsRequest {
	return V1AnalyticsRequest{
		ReqAnaType: "HIST",
		ReqPeriod:  "PT5M",
	}
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

func V1Analytics(w http.ResponseWriter, r *http.Request) {
	req := newV1AnalyticsRequest()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":400,"detail":"bad json"}`, http.StatusBadRequest)
		return
	}

	if req.AnalyticsId == "" {
		http.Error(w, `{"status":400,"detail":"analyticsId required"}`, http.StatusBadRequest)
		return
	} else if req.AnalyticsId != "ABNORMAL_BEHAVIOUR" && req.AnalyticsId != "NF_LOAD" {
		http.Error(w, `{"status":400,"detail":"unknown analyticsId"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /nnwdaf-analyticsinfo/v1/health", V1Health)
	mux.HandleFunc("POST /nnwdaf-analyticsinfo/v1/analytics", V1Analytics)

	log.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
