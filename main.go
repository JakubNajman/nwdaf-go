package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (s *Server) newV1AnalyticsRequest() V1AnalyticsRequest {
	return V1AnalyticsRequest{
		ReqAnaType: "HIST",
		ReqPeriod:  "PT5M",
	}
}

func (s *Server) V1Health(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) V1Analytics(w http.ResponseWriter, r *http.Request) {
	req := s.newV1AnalyticsRequest()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":400,"detail":"bad json"}`, http.StatusBadRequest)
		return
	}

	if req.AnalyticsId == "" {
		http.Error(w, `{"status":400,"detail":"analyticsId required"}`, http.StatusBadRequest)
		return
	}
	if !req.AnalyticsId.IsValid() {
		http.Error(w, `{"status":400,"detail":"unknown analyticsId"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)

}

func (s *Server) V1SubscriptionCreate(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := validateSubscriptionRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := s.store.Add(req)
	log.Printf("[Sub] Created %s for %s → %s", resp.SubscriptionId, req.AnalyticsId, req.NotificationUri)

	writeJSON(w, http.StatusCreated, resp) // 201
}

func (s *Server) V1SubscriptionDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deletion := s.store.Delete(id)

	if deletion != nil {
		if errors.Is(deletion, ErrSubscriptionNotFound) {
			writeError(w, http.StatusNotFound, id+" not found")
			return
		} else {
			writeError(w, http.StatusInternalServerError, deletion.Error())
			return
		}
	}

	log.Printf("[Sub] Removed %s", id)
	w.WriteHeader(http.StatusNoContent) // 204
}

func (s *Server) V1SubscriptionUpdate(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest
	id := r.PathValue("id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := validateSubscriptionRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, succ := s.store.Update(id, req)
	if succ != nil {
		writeError(w, http.StatusBadRequest, "update failed")
		return
	}

	log.Printf("[Sub] Updated %s for %s → %s", resp.SubscriptionId, req.AnalyticsId, req.NotificationUri)
	writeJSON(w, http.StatusAccepted, resp) // 202
}

func (s *Server) V1SubscriptionList(w http.ResponseWriter, r *http.Request) {
	subs := s.store.List()
	writeJSON(w, http.StatusAccepted, map[string]any{"subscriptions": subs}) // 202
}

func main() {
	srv := NewServer()

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", srv.Routes()); err != nil {
		log.Fatal(err)
	}
}
