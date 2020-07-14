package client

import (
	"github.com/PagerDuty/go-pagerduty"
)

type EventClient struct {
	client *pagerduty.V2Event
}

func NewEventClient(client *pagerduty.V2Event) *EventClient {
	return &EventClient{client: client}
}

type Manager interface {
	ManageEvent() (*pagerduty.V2EventResponse, error)
}

//ManageEvent create or resolve event using PagerDuty API
func (c *EventClient) ManageEvent() (*pagerduty.V2EventResponse, error) {
	response, error := pagerduty.ManageEvent(*c.client)
	return response, error
}
