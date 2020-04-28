package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
)

type IncidentClient interface {
	ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse
	IncidentAlerts(id string) *pagerduty.ListAlertsResponse
}

func (c *ApiClient) ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse {
	eps, err := c.client.ListIncidents(Options)
	base.CheckErr(err)
	return eps
}

func (c *ApiClient) IncidentAlerts(id string) *pagerduty.ListAlertsResponse {
	eps, err := c.client.ListIncidentAlerts(id)
	base.CheckErr(err)
	return eps
}
