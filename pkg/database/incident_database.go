package database

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"

	_ "github.com/mattn/go-sqlite3"
)

//IncidentDb structure for incidents stored in database
type IncidentDb struct {
	Id          string
	Title       string
	Service     string
	ServiceName string
	CreateAt    string
	Timer       string
	Alert       string
	Tocheck     string
	Trigger     string
}

//IncidentRepository interface
type IncidentRepository interface {
	UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{})
	SaveIncident(incident *pagerduty.Incident, incidentTimer interface{})
	GetIncident() (inc []*IncidentDb)
	UpdateIncidentState(incident *IncidentDb)
}

//UpdateIncident insert incident to database
func (d *Store) UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	service := incident.Service.ID
	serviceName := incident.Service.Summary
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("Update incidents set id = ?,title = ?,createat = ?,timer = ?, servicename = ? where service = ?")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, serviceName, service)
	base.CheckErr(err)
}

//SaveIncident insert incident to database
func (d *Store) SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	service := incident.Service.ID
	serviceName := incident.Service.Summary
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("REPLACE INTO incidents(id,title,createat,timer,service,servicename) values(?,?,?,?,?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, service, serviceName)
	base.CheckErr(err)
}

//UpdateIncidentState insert incident to database
func (d *Store) UpdateIncidentState(incident *IncidentDb) {
	alert := incident.Alert
	tocheck := incident.Tocheck
	trigger := incident.Trigger
	service := incident.Service
	stmt, err := d.db.Prepare("UPDATE incidents set alert = ? , tocheck = ? , trigger = ? where service = ?")
	base.CheckErr(err)
	_, err = stmt.Exec(alert, tocheck, trigger, service)
	base.CheckErr(err)
}

func (d *Store) GetIncident() (inc []*IncidentDb) {
	var incTmp IncidentDb
	r, err := d.db.Query("select * from incidents")
	base.CheckErr(err)
	for r.Next() {
		err := r.Scan(&incTmp.Id, &incTmp.Title, &incTmp.Service, &incTmp.ServiceName, &incTmp.CreateAt, &incTmp.Timer, &incTmp.Alert, &incTmp.Tocheck, &incTmp.Trigger)
		base.CheckErr(err)
		inc = append(inc, &incTmp)
	}
	return
}

//InitIncidentRepository create schema
func (d *Store) InitIncidentRepository() {
	incidentTable := `
	CREATE TABLE IF NOT EXISTS incidents(
		id TEXT NOT NULL,
		title TEXT NOT NULL,
		service TEXT NOT NULL UNIQUE,
		servicename TEXT NOT NULL,
		createat TEXT NOT NULL,
		timer TEXT NOT NULL,
		alert TEXT DEFAULT "N",
		tocheck TEXT DEFAULT "N",
		trigger TEXT DEFAULT "N"

	);
	`

}

func (d *Store) createTable(sqlTable string) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS incidents(
		id TEXT NOT NULL,
		title TEXT NOT NULL,
		service TEXT NOT NULL UNIQUE,
		servicename TEXT NOT NULL,
		createat TEXT NOT NULL,
		timer TEXT NOT NULL,
		alert TEXT DEFAULT "N",
		tocheck TEXT DEFAULT "N",
		trigger TEXT DEFAULT "N"

	);
	`
	_, err := d.db.Exec(sql_table)
	base.CheckErr(err)
}
