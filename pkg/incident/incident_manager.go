package incident

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
	"github.com/jaceklubzinski/pd-checker/pkg/client"

	log "github.com/sirupsen/logrus"
)

//Manager struct to manage incidents
type Manager struct {
	IncidentClient client.IncidentClient
	Options        pagerduty.ListIncidentsOptions
	Incident
}

//SetOptions set options to get incidents
func (m *Manager) SetOptions() {
	m.Options.Statuses = []string{"resolved"}
	m.Options.Until = time.Now().String()
	m.Options.Since = time.Now().AddDate(0, 0, -1).String()
	m.Options.SortBy = "created_at:desc"
}

//SetOptionsFromIncident set options from existed incident to check again triggered incidents
func (m *Manager) SetOptionsFromIncident() {
	m.Options.ServiceIDs = []string{m.Incident.ServiceID}
	m.Options.Since = base.AddDurationToDate(m.Incident.CreateAt, m.Incident.Timer)
	m.Options.Until = base.DateNow()
}

//CheckForNew check if service incident was created
func (m *Manager) CheckForNew() *pagerduty.Incident {
	registry := m.IncidentClient.ListIncidents(m.Options)
	for _, p := range registry.Incidents {
		if p.Title == "PD CHECKER - OK" {
			return &p
		}
	}
	return nil
}

//TriggerAlert Trigger new alert for service
func (m *Manager) TriggerAlert() {
	if m.Incident.Alert == "Y" && m.Incident.Trigger == "N" {
		log.WithFields(log.Fields{
			"ServiceName": m.Incident.ServiceName,
			"ServiceID":   m.Incident.ServiceID,
		}).Info("Trigger alert for service")
		m.Incident.Trigger = "Y"
		m.Incident.Alert = "N"
	}
}

//AlertDetails get detailed information from PagerDuty Incident
func (m *Manager) AlertDetails(id string) (repeatTimer interface{}) {
	a := m.IncidentClient.IncidentAlerts(id)
	for _, i := range a.Alerts {
		repeatTimer = i.Body["details"]
	}
	return repeatTimer
}
