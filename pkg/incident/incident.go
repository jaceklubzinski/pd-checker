package incident

import (
	"fmt"
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
	Repository     database.IncidentRepository
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
			i.Repository.SaveIncident(&p, i.repeatTimer)
		}
	}
}

//MarkToCheck check if service incident should be considered again base on last creation time and duraiton
func (i *IncidentService) MarkToCheck() {
	fmt.Println("Mark")
	for _, incident := range i.Incidents {
		serviceTimer, err := time.ParseDuration(incident.Timer)
		base.CheckErr(err)
		lastTillNow := base.LastTillNowDuration(incident.CreateAt)
		if lastTillNow > serviceTimer {
			incident.Tocheck = "Y"
			fmt.Println("to check")
		}
	}
	fmt.Println("End Mark")
}

//CheckToAlert check if service incident was created
func (i *IncidentService) CheckToAlert() {
	count := 0
	fmt.Println("Check")
	for _, incident := range i.Incidents {
		fmt.Println("tocheck: ", incident.Tocheck)
		if incident.Tocheck == "Y" {
			i.Options.ServiceIDs = []string{incident.Service}
			i.Options.Since = base.NewStartDate(incident.CreateAt, incident.Timer)
			registry := i.IncidentClient.ListIncidents(i.Options)
			for _, p := range registry.Incidents {
				if p.Title == "PD CHECKER - OK" {
					count++
					i.repeatTimer = i.IncidentClient.AlertDetail(p.Id)
					i.Repository.SaveIncident(&p, i.repeatTimer)
				}
			}
		}
		if count == 0 && incident.Trigger == "N" && incident.Tocheck == "Y" {
			incident.Alert = "Y"
			incident.Tocheck = "N"
			fmt.Println("To alert")
		}
	}
	fmt.Println("End Check")
}

func (i *IncidentService) Alert() {
	for _, incident := range i.Incidents {
		if incident.Alert == "Y" && incident.Trigger == "N" {
			fmt.Println("Trigger alert for service ", incident.Service)
			incident.Trigger = "Y"
			incident.Alert = "N"
		}
		i.Repository.UpdateIncident(incident)
	}
}
