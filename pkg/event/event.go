package event

import "github.com/PagerDuty/go-pagerduty"

type optionEvent struct {
	event *pagerduty.V2Event
}

type responseEvent struct {
	response *pagerduty.V2EventResponse
}

func NewEvent(event *pagerduty.V2Event) *optionEvent {
	return &optionEvent{event}
}

type manageEvent struct {
	options  *optionEvent
	response *pagerduty.V2EventResponse
}
