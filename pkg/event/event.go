package event

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/metrics"
)

//ManageEvent struct that operate on alert
type ManageEvent struct {
	Options  *pagerduty.V2Event
	Response *pagerduty.V2EventResponse
	message  string
	metrics.RecordMetricsEvent
}

func NewEvent(Options *pagerduty.V2Event) *ManageEvent {
	return &ManageEvent{Options: Options}
}
