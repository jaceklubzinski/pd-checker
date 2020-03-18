package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
)

type IncidentClient interface {
	ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse
	AlertDetail(id string) (repeatTimer interface{})
}

func (c *ApiClient) ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse {
	eps, err := c.client.ListIncidents(Options)
	base.CheckErr(err)
	return eps
}

func (c *ApiClient) AlertDetail(id string) (repeatTimer interface{}) {
	eps, err := c.client.ListIncidentAlerts(id)
	base.CheckErr(err)
	for _, i := range eps.Alerts {
		repeatTimer = i.Body["details"]
	}
	return repeatTimer
}
