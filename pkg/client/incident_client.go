package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
)

//IncidentClient interface for managing PagerDuty client library
type IncidentClient interface {
	ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse
	IncidentAlerts(id string) *pagerduty.ListAlertsResponse
	CreateIncident(Options *pagerduty.CreateIncidentOptions) error
}

//ListIncidents list of the PagerDuty incidents for specific service
func (c *ApiClient) ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse {
	eps, err := c.client.ListIncidents(Options)
	base.CheckErr(err)
	return eps
}

//IncidentAlerts List of the PagerDuty alerts for specific incient
func (c *ApiClient) IncidentAlerts(id string) *pagerduty.ListAlertsResponse {
	eps, err := c.client.ListIncidentAlerts(id)
	base.CheckErr(err)
	return eps
}

//CreateIncident create new PagerDuty incident
func (c *ApiClient) CreateIncident(Options *pagerduty.CreateIncidentOptions) error {
	from := "pd-checker@no-replay"
	_, err := c.client.CreateIncident(from, Options)
	return err
}
