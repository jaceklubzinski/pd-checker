package event

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/metrics"
)

//Event struct that operate on alert
type Event struct {
	Options  *pagerduty.V2Event
	Response *pagerduty.V2EventResponse
	Manager  client.Manager
	message  string
	metrics.RecordMetricsEvent
}

//NewEvent new event
func NewEvent(Options *pagerduty.V2Event) *Event {
	return &Event{Options: Options}
}
