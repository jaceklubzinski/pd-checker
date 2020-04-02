package incident

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/database"
	log "github.com/sirupsen/logrus"
)

//IncidentService struct to manage incidents
type IncidentService struct {
	IncidentClient client.IncidentClient
	Options        pagerduty.ListIncidentsOptions
	DbRepository   database.IncidentRepository
	repeatTimer    interface{}
	Incidents      []*database.IncidentDb
	Incident       *database.IncidentDb
	Services       []*database.ServiceDB
}

//IncidentOptions set options to get incidents
func (i *IncidentService) IncidentOptions() {
	i.Options.Statuses = []string{"resolved"}
	i.Options.Until = time.Now().String()
	i.Options.Since = time.Now().AddDate(0, 0, -1).String()
	i.Options.SortBy = "created_at:asc"
}

//CheckServiceIncident count incidents for service
func (i *IncidentService) CheckServiceIncident() {
	for _, service := range i.Services {
		i.Options.ServiceIDs = []string{service.ID}
		registry := i.IncidentClient.ListIncidents(i.Options)
		for _, incident := range registry.Incidents {
			if incident.Title == "PD CHECKER - OK" {
				i.repeatTimer = i.IncidentClient.AlertDetail(incident.Id)
				i.DbRepository.SaveIncident(&incident, i.repeatTimer)
				log.WithFields(log.Fields{
					"ServiceName":    incident.ServiceName,
					"ServiceID":      incident.ServiceID,
					"IncidentNumber": p.IncidentNumber,
				}).Info("New incident for service registered")
			}
		}
	}
}

//MarkToCheck check if service incident should be considered again base on last creation time and duraiton
func (i *IncidentService) MarkToCheck() {
	for _, incident := range i.Incidents {
		serviceTimer, err := time.ParseDuration(incident.Timer)
		base.CheckErr(err)
		lastTillNow := base.LastTillNowDuration(incident.CreateAt)
		if lastTillNow > serviceTimer {
			incident.ToCheck = "Y"
			log.WithFields(log.Fields{
				"ServiceName":  incident.ServiceName,
				"ServiceID":    incident.ServiceID,
				"LastCreateAt": incident.CreateAt,
			}).Info("Checking for new alert")
		} else {
			log.WithFields(log.Fields{
				"ServiceName": incident.ServiceName,
				"ServiceID":   incident.ServiceID,
			}).Info("Service has working PagerDuty integration")
		}
	}
}

//CheckToAlert check if service incident was created
func (i *IncidentService) CheckToAlert() {
	for _, incident := range i.Incidents {
		if incident.ToCheck == "Y" {
			i.Options.ServiceIDs = []string{incident.ServiceID}
			i.Options.Since = base.AddDurationToDate(incident.CreateAt, incident.Timer)
			i.Options.Until = base.DateNow()
			registry := i.IncidentClient.ListIncidents(i.Options)
			for _, p := range registry.Incidents {
				if p.Title == "PD CHECKER - OK" {
					i.repeatTimer = i.IncidentClient.AlertDetail(p.Id)
					incident.ToCheck = "N"
					incident.Trigger = "N"
					incident.Alert = "N"
					i.DbRepository.UpdateIncident(&p, i.repeatTimer)
					log.WithFields(log.Fields{
						"serviceName":    incident.ServiceName,
						"ServiceID":      incident.ServiceID,
						"IncidentNumber": p.IncidentNumber,
					}).Info("New incident for service registered")
				}
			}
			if incident.ToCheck == "Y" && incident.Trigger == "N" {
				incident.Alert = "Y"
				log.WithFields(log.Fields{
					"ServiceName": incident.ServiceName,
					"ServiceID":   incident.ServiceID,
				}).Info("New alert will be created")
			} else if incident.ToCheck == "Y" && incident.Trigger == "Y" {
				log.WithFields(log.Fields{
					"ServiceName": incident.ServiceName,
					"ServiceID":   incident.ServiceID,
				}).Info("Alert for service already created")
			}
		}
	}
}

//Alert Trigger new alert for service
func (i *IncidentService) Alert() {
	for _, incident := range i.Incidents {
		if incident.Alert == "Y" && incident.Trigger == "N" {
			log.WithFields(log.Fields{
				"ServiceName": incident.ServiceName,
				"ServiceID":   incident.ServiceID,
			}).Info("Trigger alert for service")
			incident.Trigger = "Y"
			incident.Alert = "N"
		}
		i.DbRepository.UpdateIncidentState(incident)
	}
}
