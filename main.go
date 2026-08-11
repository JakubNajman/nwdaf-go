package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var store = NewSubscriptionStore()

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

func V1SubscriptionCreate(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if !req.AnalyticsId.IsValid() {
		writeError(w, http.StatusBadRequest, "invalid analyticsId")
		return
	}
	if req.NotificationUri == "" {
		writeError(w, http.StatusBadRequest, "notificationUri required")
		return
	}

	resp := store.Add(req)
	log.Printf("[Sub] Created %s for %s → %s", resp.SubscriptionId, req.AnalyticsId, req.NotificationUri)

	writeJSON(w, http.StatusCreated, resp) // 201
}

func V1SubscriptionDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deletion := store.Delete(id)

	if !deletion {
		writeError(w, http.StatusNotFound, id+" not found")
		return
	}

	log.Printf("[Sub] Removed %s", id)
	w.WriteHeader(http.StatusNoContent) // 204
}

func V1SubscriptionUpdate(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest
	id := r.PathValue("id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if !req.AnalyticsId.IsValid() {
		writeError(w, http.StatusBadRequest, "invalid analyticsId")
		return
	}
	if req.NotificationUri == "" {
		writeError(w, http.StatusBadRequest, "notificationUri required")
		return
	}

	resp, succ := store.Update(id, req)
	if !succ {
		writeError(w, http.StatusBadRequest, "update failed")
		return
	}

	log.Printf("[Sub] Updated %s for %s → %s", resp.SubscriptionId, req.AnalyticsId, req.NotificationUri)
	writeJSON(w, http.StatusAccepted, resp) // 202
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /nnwdaf-analyticsinfo/v1/health", V1Health)
	mux.HandleFunc("POST /nnwdaf-analyticsinfo/v1/analytics", V1Analytics)
	mux.HandleFunc("POST /nnwdaf-eventssubscription/v1/subscriptions", V1SubscriptionCreate)
	// mux.HandleFunc("GET /nnwdaf-eventssubscription/v1/subscriptions", V1SubscriptionList)
	mux.HandleFunc("PUT /nnwdaf-eventssubscription/v1/subscriptions/{id}", V1SubscriptionUpdate)
	mux.HandleFunc("DELETE /nnwdaf-eventssubscription/v1/subscriptions/{id}", V1SubscriptionDelete)

	log.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
