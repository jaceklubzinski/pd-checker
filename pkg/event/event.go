package event

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/metrics"
)

type ManageEvent struct {
	Options      *pagerduty.V2Event
	Response     *pagerduty.V2EventResponse
	EventMetrics metrics.RecordMetricsEvent
	message      string
}
