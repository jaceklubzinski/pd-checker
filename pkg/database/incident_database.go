package database

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"

	_ "github.com/mattn/go-sqlite3"
)

//IncidentDb structure for incidents stored in database
type IncidentDb struct {
	Id       string
	Title    string
	Service  string
	CreateAt string
	Timer    string
	Alert    string
	Tocheck  string
}

//IncidentRepository interface
type IncidentRepository interface {
	SaveIncident(incident *pagerduty.Incident, incidentTimer interface{})
	GetIncident() (inc []*IncidentDb)
}

//SaveIncident insert incident to database
func (d *Store) SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	service := incident.Service.ID
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("REPLACE INTO incidents values(?,?,?,?,?,?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, service, createAt, incidentTimer, "N", "N")
	base.CheckErr(err)
}

func (d *Store) GetIncident() (inc []*IncidentDb) {
	var incTmp IncidentDb
	r, err := d.db.Query("select * from incidents")
	base.CheckErr(err)
	for r.Next() {
		err := r.Scan(&incTmp.Id, &incTmp.Title, &incTmp.Service, &incTmp.CreateAt, &incTmp.Timer, &incTmp.Alert, &incTmp.Tocheck)
		base.CheckErr(err)
		inc = append(inc, &incTmp)
	}
	return
}

//InitIncidentRepository create schema
func (d *Store) InitIncidentRepository() {
	sql_table := `
	CREATE TABLE IF NOT EXISTS incidents(
		id TEXT NOT NULL,
		title TEXT NOT NULL,
		service TEXT NOT NULL UNIQUE,
		createat TEXT NOT NULL,
		timer TEXT NOT NULL,
		alert TEXT DEFAULT "N",
		tocheck TEXT DEFAULT "N"
	);
	`
	_, err := d.db.Exec(sql_table)
	base.CheckErr(err)
}
