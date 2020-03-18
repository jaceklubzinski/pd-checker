package database

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"

	_ "github.com/mattn/go-sqlite3"
)

//Incident structure for incidents stored in database
type Incident struct {
	id       string
	title    string
	service  string
	createAt string
	timer    string
}

//IncidentRegister register for incidents
type IncidentRegister struct {
	Incidents []Incident
}

//IncidentRepository
type IncidentRepository interface {
	SaveIncident(incident *pagerduty.Incident, incidentTimer interface{})
}

//SaveIncident insert incident to database
func (d *store) SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	service := incident.Service.ID
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("REPLACE INTO incidents values(?,?,?,?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, service, createAt, incidentTimer)
	base.CheckErr(err)
}

func (d *store) GetIncident() (incidents IncidentRegister) {
	var inc Incident
	r, err := d.db.Query("select * from incidents")
	base.CheckErr(err)
	for r.Next() {
		err := r.Scan(&inc.id, &inc.title, &inc.service, &inc.createAt, &inc.timer)
		base.CheckErr(err)
		incidents.Incidents = append(incidents.Incidents, inc)
	}
	return
}

//InitIncidentRepository create schema
func (d *store) InitIncidentRepository() {
	_, err := d.db.Exec("create table if not exists incidents (id text NOT NULL,title text NOT NULL, service text NOT NULL UNIQUE, createat text NOT NULL, timer text)")
	base.CheckErr(err)
}
