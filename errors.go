package main

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidAnalyticsId   = errors.New("invalid analytics id")
	ErrMissingNotifUri      = errors.New("missing notif uri")
)

func validateSubscriptionRequest(req SubscriptionRequest) error {
	if !req.AnalyticsId.IsValid() {
		return ErrInvalidAnalyticsId
	}
	if req.NotificationUri == "" {
		return ErrMissingNotifUri
	}
	return nil
}
