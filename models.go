package main

import (
	"time"
)

type AnalyticsId string

const (
	LoadLevelInformation AnalyticsId = "LOAD_LEVEL_INFORMATION"
	NetworkPerformance   AnalyticsId = "NETWORK_PERFORMANCE"
	AbnormalBehaviour    AnalyticsId = "ABNORMAL_BEHAVIOUR"
	ServiceExperience    AnalyticsId = "SERVICE_EXPERIENCE"
	UeMobility           AnalyticsId = "UE_MOBILITY"
	UeCommunication      AnalyticsId = "UE_COMMUNICATION"
	ExpectedUeBehaviour  AnalyticsId = "EXPECTED_UE_BEHAVIOURAL_PARAMETERS"
	DnPerformance        AnalyticsId = "DN_PERFORMANCE"
	Dispersion           AnalyticsId = "DISPERSION"
)

func (a AnalyticsId) IsValid() bool {
	switch a {
	case LoadLevelInformation, NetworkPerformance, AbnormalBehaviour,
		ServiceExperience, UeMobility, UeCommunication,
		ExpectedUeBehaviour, DnPerformance, Dispersion:
		return true
	}
	return false
}

type V1HealthResponse struct {
	Status        string `json:"status"`
	NRF           string `json:"nrf"`
	ADRF          string `json:"adrf"`
	Subscriptions int    `json:"subscriptions"`
}

type V1AnalyticsRequest struct {
	AnalyticsId AnalyticsId `json:"analyticsId"`
	ReqAnaType  string      `json:"reqAnaType"`
	ReqPeriod   string      `json:"reqPeriod"`
	StartTs     *time.Time  `json:"startTs"`
	EndTs       *time.Time  `json:"endTs"`
}

type V1AnalyticsResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

type ReqAnaType string

const (
	Hist        ReqAnaType = "HIST"
	Pred        ReqAnaType = "PRED"
	HistAndPred ReqAnaType = "HIST_AND_PRED"
)

type NotifMethod string

const (
	Periodic  NotifMethod = "PERIODIC"
	Threshold NotifMethod = "THRESHOLD"
	OneTime   NotifMethod = "ONE_TIME"
)

type Snssai struct {
	Sst int    `json:"sst"`
	Sd  string `json:"sd"`
}

type PlmnId struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

type SnssaiDnnFilter struct {
	Snssai Snssai `json:"snssai"`
	Dnn    string `json:"dnn"`
}

type TgtUe struct {
	AnyUe        *bool   `json:"anyUe,omitempty"`
	Supi         *string `json:"supi,omitempty"`
	ExterGroupId *string `json:"exterGroupId,omitempty"`
}

type AnalyticsRequest struct {
	AnalyticsId     AnalyticsId      `json:"analyticsId"`
	ReqAnaType      ReqAnaType       `json:"reqAnaType"`
	ReqPeriod       string           `json:"reqPeriod"`
	TgtUe           *TgtUe           `json:"tgtUe,omitempty"`
	SnssaiDnnFilter *SnssaiDnnFilter `json:"snssaiDnnFilter,omitempty"`
	StartTs         *time.Time       `json:"startTs,omitempty"`
	EndTs           *time.Time       `json:"endTs,omitempty"`
}

func NewAnalyticsRequest() AnalyticsRequest {
	return AnalyticsRequest{
		ReqAnaType: Hist,
		ReqPeriod:  "PT5M",
	}
}

type EventReportingRequirement struct {
	NotifMethod    NotifMethod `json:"notifMethod"`
	RepPeriod      int         `json:"repPeriod"`
	NotifThreshold *float64    `json:"notifThreshold,omitempty"`
}

type SubscriptionRequest struct {
	AnalyticsId               AnalyticsId               `json:"analyticsId"`
	NotificationUri           string                    `json:"notificationUri"`
	EventReportingRequirement EventReportingRequirement `json:"eventReportingRequirement"`
	SnssaiDnnFilter           *SnssaiDnnFilter          `json:"snssaiDnnFilter,omitempty"`
	TgtUe                     *TgtUe                    `json:"tgtUe,omitempty"`
	Expiry                    *string                   `json:"expiry,omitempty"`
}

type SubscriptionResponse struct {
	SubscriptionId string  `json:"subscriptionId"`
	AnalyticsId    string  `json:"analyticsId"`
	Expiry         *string `json:"expiry,omitempty"`
}

type EventNotification struct {
	SubscriptionId string         `json:"subscriptionId"`
	AnalyticsId    string         `json:"analyticsId"`
	TimeStamp      string         `json:"timeStamp"`
	Data           map[string]any `json:"data"`
}
