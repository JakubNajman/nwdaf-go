package main

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidAnalyticsId   = errors.New("invalid analytics id")
	ErrMissingNotifUri      = errors.New("missing notif uri")
	ErrInvalidRepPerdio     = errors.New("repPeriod must be > 0")
)

func validateSubscriptionRequest(req SubscriptionRequest) error {
	if !req.AnalyticsId.IsValid() {
		return ErrInvalidAnalyticsId
	}
	if req.NotificationUri == "" {
		return ErrMissingNotifUri
	}
	if req.EventReportingRequirement.RepPeriod <= 0 {
		return ErrInvalidRepPerdio
	}
	return nil
}
