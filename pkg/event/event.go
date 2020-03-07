package event

import "github.com/PagerDuty/go-pagerduty"

type ManageEvent struct {
	Options  *pagerduty.V2Event
	Response *pagerduty.V2EventResponse
}
