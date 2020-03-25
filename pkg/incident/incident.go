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
	IncidentClient client.IncidentClient
	Options        pagerduty.ListIncidentsOptions
	Repository     database.IncidentRepository
	count          int
	repeatTimer    interface{}
	Incidents      []*database.IncidentDb
}

//IncidentOptions set options to get incidents
func (i *IncidentService) IncidentOptions() {
	i.Options.Statuses = []string{"resolved"}
	i.Options.Until = time.Now().String()
	i.Options.Since = time.Now().AddDate(0, 0, -1).String()
	i.Options.SortBy = "created_at:asc"
}

//CountIncidentService count incidents for service
func (i *IncidentService) CountIncidentService() {
	i.count = 0
	registry := i.IncidentClient.ListIncidents(i.Options)
	for _, p := range registry.Incidents {
		if p.Title == "PD CHECKER - OK" {
			i.repeatTimer = i.IncidentClient.AlertDetail(p.Id)
			i.Repository.SaveIncident(&p, i.repeatTimer)
			i.count++
		}
	}
}

//CounterInfo check incichdent checker count for specific service
func (i *IncidentService) CounterInfo() {
	log.Printf("Service name: %s Incident created: %d", i.Options.ServiceIDs, i.count)
}
