package main

import (
	"fmt"
	"sync"
)

type Subscription struct {
	SubscriptionId string
	Request        SubscriptionRequest
}

type SubscriptionInfo struct {
	SubscriptionId  string `json:"subscriptionId"`
	AnalyticsId     string `json:"analyticsId"`
	NotificationUri string `json:"notificationUri"`
	RepPeriod       int    `json:"repPeriod"`
}

func (s *Subscription) Info() SubscriptionInfo {
	out := SubscriptionInfo{
		SubscriptionId:  s.SubscriptionId,
		AnalyticsId:     string(s.Request.AnalyticsId),
		NotificationUri: s.Request.NotificationUri,
		RepPeriod:       s.Request.EventReportingRequirement.RepPeriod,
	}

	return out
}

type SubscriptionStore struct {
	mu            sync.Mutex
	Subscriptions map[string]*Subscription
	NextId        int
}

func NewSubscriptionStore() *SubscriptionStore {
	s := SubscriptionStore{}
	s.Subscriptions = make(map[string]*Subscription)
	return &s
}

func (s *SubscriptionStore) Add(req SubscriptionRequest) SubscriptionResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.NextId++
	subID := fmt.Sprintf("sub-%d", s.NextId)

	s.Subscriptions[subID] = &Subscription{SubscriptionId: subID, Request: req}

	return SubscriptionResponse{
		SubscriptionId: subID,
		AnalyticsId:    string(req.AnalyticsId),
		Expiry:         req.Expiry,
	}
}

func (s *SubscriptionStore) Get(id string) (SubscriptionRequest, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.Subscriptions[id]

	if !ok {
		return SubscriptionRequest{}, false
	}
	return sub.Request, true
}

func (s *SubscriptionStore) Update(id string, req SubscriptionRequest) (SubscriptionResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.Subscriptions[id]

	if !ok {
		return SubscriptionResponse{}, false
	}

	sub.Request = req

	return SubscriptionResponse{
		SubscriptionId: id,
		AnalyticsId:    string(req.AnalyticsId),
		Expiry:         req.Expiry,
	}, true
}

func (s *SubscriptionStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.Subscriptions[id]

	if !ok {
		return false
	}

	delete(s.Subscriptions, id)
	return true
}

func (s *SubscriptionStore) List() []SubscriptionInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]SubscriptionInfo, 0, len(s.Subscriptions))
	for _, sub := range s.Subscriptions {
		out = append(out, sub.Info())
	}
	return out
}
