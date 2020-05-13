package incident

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/database"

	log "github.com/sirupsen/logrus"
)

//Manager struct to manage incidents
type Manager struct {
	IncidentClient client.IncidentClient
	Options        pagerduty.ListIncidentsOptions
}

//SetOptions set options to get incidents
func (m *Manager) SetOptions() {
	m.Options.Statuses = []string{"resolved"}
	m.Options.Until = time.Now().String()
	m.Options.Since = time.Now().AddDate(0, 0, -1).String()
	m.Options.SortBy = "created_at:asc"
}

//SetOptionsFromIncident set options from existed incident to check again triggered incidents
func (m *Manager) SetOptionsFromIncident(inc *database.Incident) {
	m.Options.ServiceIDs = []string{inc.ServiceID}
	m.Options.Since = base.AddDurationToDate(inc.CreateAt, inc.Timer)
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

//Alert Trigger new alert for service
func (m *Manager) Alert(inc *database.Incident) {
	if inc.Alert == "Y" && inc.Trigger == "N" {
		log.WithFields(log.Fields{
			"ServiceName": inc.ServiceName,
			"ServiceID":   inc.ServiceID,
		}).Info("Trigger alert for service")
		inc.Trigger = "Y"
		inc.Alert = "N"
	}
}

func (m *Manager) AlertDetails(id string) (repeatTimer interface{}) {
	a := m.IncidentClient.IncidentAlerts(id)
	for _, i := range a.Alerts {
		repeatTimer = i.Body["details"]
	}
	return repeatTimer
}
