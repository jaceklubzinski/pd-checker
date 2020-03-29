package incident

import (
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/database"
)

//IncidentService struct to manage incidents
type IncidentService struct {
	IncidentClient client.IncidentClient
	Options        pagerduty.ListIncidentsOptions
	DbRepository   database.IncidentRepository
	repeatTimer    interface{}
	Incidents      []*database.IncidentDb
	Incident       *database.IncidentDb
}

//IncidentOptions set options to get incidents
func (i *IncidentService) IncidentOptions() {
	i.Options.Statuses = []string{"resolved"}
	i.Options.Until = time.Now().String()
	i.Options.Since = time.Now().AddDate(0, 0, -1).String()
	i.Options.SortBy = "created_at:asc"
}

//WriteToDBIncidentService count incidents for service
func (i *IncidentService) WriteToDBIncidentService() {
	registry := i.IncidentClient.ListIncidents(i.Options)
	for _, p := range registry.Incidents {
		if p.Title == "PD CHECKER - OK" {
			i.repeatTimer = i.IncidentClient.AlertDetail(p.Id)
			i.DbRepository.SaveIncident(&p, i.repeatTimer)
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
			incident.Tocheck = "Y"
			log.Printf("For service %s (%s) last alert was created at %s  - checking for new alert", incident.ServiceName, incident.Service, incident.CreateAt)
		} else {
			log.Printf("Service %s (%s) has working PagerDuty integration and created checker incident", incident.ServiceName, incident.Service)
		}
	}
}

//CheckToAlert check if service incident was created
func (i *IncidentService) CheckToAlert() {
	for _, incident := range i.Incidents {
		if incident.Tocheck == "Y" {
			i.Options.ServiceIDs = []string{incident.Service}
			i.Options.Since = base.AddDurationToDate(incident.CreateAt, incident.Timer)
			i.Options.Until = base.DateNow()
			registry := i.IncidentClient.ListIncidents(i.Options)
			for _, p := range registry.Incidents {
				if p.Title == "PD CHECKER - OK" {
					i.repeatTimer = i.IncidentClient.AlertDetail(p.Id)
					incident.Tocheck = "N"
					incident.Trigger = "N"
					incident.Alert = "N"
					i.DbRepository.UpdateIncident(&p, i.repeatTimer)
					log.Printf("New incident for service %s (%s) registered %d", incident.ServiceName, incident.Service, p.IncidentNumber)
				}
			}
			if incident.Tocheck == "Y" && incident.Trigger == "N" {
				incident.Alert = "Y"
				log.Printf("New alert will be created for service %s (%s)", incident.ServiceName, incident.Service)
			} else if incident.Tocheck == "Y" && incident.Trigger == "Y" {
				log.Printf("Alert for service %s (%s) already created", incident.ServiceName, incident.Service)
			}
		}
	}
}

func (i *IncidentService) Alert() {
	for _, incident := range i.Incidents {
		if incident.Alert == "Y" && incident.Trigger == "N" {
			log.Printf("Trigger alert for service %s (%s)", incident.ServiceName, incident.Service)
			incident.Trigger = "Y"
			incident.Alert = "N"
		}
		i.DbRepository.UpdateIncidentState(incident)
	}
}
