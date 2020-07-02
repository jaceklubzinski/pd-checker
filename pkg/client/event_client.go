package client

import "github.com/PagerDuty/go-pagerduty"

type Manager interface {
	ManageEvent(options *pagerduty.V2Event) (*pagerduty.V2EventResponse, error)
}

//ManageEvent create or resolve event using PagerDuty API
func ManageEvent(options *pagerduty.V2Event) (*pagerduty.V2EventResponse, error) {
	response, error := pagerduty.ManageEvent(*options)
	return response, error
}
