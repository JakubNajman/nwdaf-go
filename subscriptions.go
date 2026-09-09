package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Subscription struct {
	SubscriptionId string
	Request        SubscriptionRequest
	cancel         context.CancelFunc
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

	s.NextId++
	subID := fmt.Sprintf("sub-%d", s.NextId)

	ctx, cancel := context.WithCancel(context.Background())
	s.Subscriptions[subID] = &Subscription{
		SubscriptionId: subID,
		Request:        req,
		cancel:         cancel,
	}

	s.mu.Unlock()

	go s.notificationLoop(ctx, subID, req)

	return SubscriptionResponse{
		SubscriptionId: subID,
		AnalyticsId:    string(req.AnalyticsId),
		Expiry:         req.Expiry,
	}
}

func (s *SubscriptionStore) Get(id string) (SubscriptionRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.Subscriptions[id]

	if !ok {
		return SubscriptionRequest{}, ErrSubscriptionNotFound
	}
	return sub.Request, nil
}

func (s *SubscriptionStore) Update(id string, req SubscriptionRequest) (SubscriptionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.Subscriptions[id]

	if !ok {
		return SubscriptionResponse{}, ErrSubscriptionNotFound
	}

	sub.cancel()

	ctx, cancel := context.WithCancel(context.Background())
	sub.Request = req
	sub.cancel = cancel

	s.mu.Unlock()

	go s.notificationLoop(ctx, id, req)

	return SubscriptionResponse{
		SubscriptionId: id,
		AnalyticsId:    string(req.AnalyticsId),
		Expiry:         req.Expiry,
	}, nil
}

func (s *SubscriptionStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, ok := s.Subscriptions[id]

	if !ok {
		return ErrSubscriptionNotFound
	}

	sub.cancel()
	delete(s.Subscriptions, id)
	return nil
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

func (s *SubscriptionStore) notificationLoop(ctx context.Context, id string, req SubscriptionRequest) {
	period := time.Duration(req.EventReportingRequirement.RepPeriod) * time.Second
	threshold := req.EventReportingRequirement.NotifThreshold

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	log.Printf("[Sub] Loop started for %s (every %v)", id, period)

	for {
		select {
		case <-ctx.Done():
			log.Print("[Sub] Loop stopped for %s", id)
			return

		case <-ticker.C:
			result := s.compute(req)

			if threshold != nil {
				load, _ := result["loadLevelInformaion"].(float64)
				if load < *threshold {
					continue
				}
			}

			log.Printf("[Sub] Would notify %s → %v", id, result)
		}
	}
}

func (s *SubscriptionStore) compute(req SubscriptionRequest) map[string]any {
	return map[string]any{
		"loadLevelInformation": 42.0,
		"analyticsId":          string(req.AnalyticsId),
	}
}
