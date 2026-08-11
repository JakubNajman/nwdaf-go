package main

import (
	"net/http"
)

type Server struct {
	store *SubscriptionStore
}

func NewServer() *Server {
	server := Server{store: NewSubscriptionStore()}
	return &server
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /nnwdaf-analyticsinfo/v1/health", s.V1Health)
	mux.HandleFunc("POST /nnwdaf-analyticsinfo/v1/analytics", s.V1Analytics)
	mux.HandleFunc("POST /nnwdaf-eventssubscription/v1/subscriptions", s.V1SubscriptionCreate)
	mux.HandleFunc("GET /nnwdaf-eventssubscription/v1/subscriptions", s.V1SubscriptionList)
	mux.HandleFunc("PUT /nnwdaf-eventssubscription/v1/subscriptions/{id}", s.V1SubscriptionUpdate)
	mux.HandleFunc("DELETE /nnwdaf-eventssubscription/v1/subscriptions/{id}", s.V1SubscriptionDelete)
	return mux
}
