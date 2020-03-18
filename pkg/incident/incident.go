package incident

import (
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/database"
)

//IncidentService struct to manage incidents
type IncidentService struct {
	Incident    client.IncidentClient
	Options     pagerduty.ListIncidentsOptions
	Repository  database.IncidentRepository
	count       int
	repeatTimer interface{}
	Incidents   database.IncidentRegister
}

//IncidentOptions set options to get incidents
func (i *IncidentService) IncidentOptions() {
	i.Options.Statuses = []string{"resolved"}
	i.Options.Until = time.Now().String()
	i.Options.Since = time.Now().AddDate(0, 0, -1).String()
	i.Options.SortBy = "created_at:asc"
}

//GetCheckerCount count incidents for service
func (i *IncidentService) GetCheckerCount() {
	i.count = 0
	registry := i.Incident.ListIncidents(i.Options)
	for _, p := range registry.Incidents {
		if p.Title == "PD CHECKER - OK" {
			i.repeatTimer = i.Incident.AlertDetail(p.Id)
			i.Repository.SaveIncident(&p, i.repeatTimer)
			i.count++
		}
	}
}

//CheckerCount check incichdent checker count for specific service
func (i *IncidentService) CheckerCount() {
	log.Printf("Service name: %s Incident created: %d", i.Options.ServiceIDs, i.count)
}
