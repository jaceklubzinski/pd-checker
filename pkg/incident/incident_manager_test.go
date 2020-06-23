package incident

import (
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/magiconair/properties/assert"
)

type MockApiClient struct{}

type MockIncidentClient interface {
	ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse
	IncidentAlerts(id string) *pagerduty.ListAlertsResponse
	CreateIncident(Options *pagerduty.CreateIncidentOptions) error
}

func (c *MockApiClient) ListIncidents(Options pagerduty.ListIncidentsOptions) *pagerduty.ListIncidentsResponse {
	return &pagerduty.ListIncidentsResponse{
		Incidents: []pagerduty.Incident{
			pagerduty.Incident{
				Title: "PD CHECKER - OK",
			},
		},
	}
}

func (c *MockApiClient) CreateIncident(Options *pagerduty.CreateIncidentOptions) error {
	return nil
}

func (c *MockApiClient) IncidentAlerts(id string) *pagerduty.ListAlertsResponse {
	return &pagerduty.ListAlertsResponse{
		Alerts: []pagerduty.IncidentAlert{
			pagerduty.IncidentAlert{
				Body: map[string]interface{}{
					"details": "1s",
				},
			},
		},
	}
}

func TestCheckForNew(t *testing.T) {
	client := &MockApiClient{}
	incidents := Manager{IncidentClient: client}
	if i := incidents.CheckForNew(); i != nil {
		assert.Equal(t, i.Title, "PD CHECKER - OK")
	}
}

func TestAlertDetails(t *testing.T) {
	client := &MockApiClient{}
	incidents := Manager{IncidentClient: client}
	repeatTimer := incidents.AlertDetails("id")
	assert.Equal(t, repeatTimer, "1s")
}

func TestTriggerAlert(t *testing.T) {
	client := &MockApiClient{}
	i := Manager{
		IncidentClient: client,
		Incident: Incident{
			Alert:   "Y",
			Trigger: "N",
		},
	}
	i.TriggerAlert()
	assert.Equal(t, i.Alert, "N")
	assert.Equal(t, i.Trigger, "Y")
	i.Trigger = "N"
	i.TriggerAlert()
	assert.Equal(t, i.Alert, "N")
	assert.Equal(t, i.Trigger, "N")
}
