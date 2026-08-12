package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, detail string) {
	writeJSON(w, status, map[string]any{
		"status": status,
		"detail": detail,
	})
}

func validateSubscriptionRequest(req SubscriptionRequest) error {
	if !req.AnalyticsId.IsValid() {
		return ErrInvalidAnalyticsId
	}
	if req.NotificationUri == "" {
		return ErrMissingNotifUri
	}
	return nil
}
