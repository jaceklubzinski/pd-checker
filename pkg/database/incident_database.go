package database

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"

	_ "github.com/mattn/go-sqlite3"
)

//IncidentDb structure for incidents stored in database
type IncidentDb struct {
	ID          string
	Title       string
	ServiceID   string
	ServiceName string
	CreateAt    string
	Timer       string
	Alert       string
	ToCheck     string
	Trigger     string
}

//IncidentRepository interface
type IncidentRepository interface {
	UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{})
	SaveIncident(incident *pagerduty.Incident, incidentTimer interface{})
	GetIncident() (inc []*IncidentDb)
	UpdateIncidentState(incident *IncidentDb)
}

//UpdateIncident update incident creation time and timer base on service id
func (d *Store) UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	service := incident.Service.ID
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("Update incidents set id = ?,title = ?,createat = ?,timer = ? where serviceid = ?")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, service)
	base.CheckErr(err)
}

//SaveIncident insert incident to database
func (d *Store) SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	serviceID := incident.Service.ID
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("REPLACE INTO incidents(id,title,createat,timer,serviceid) values(?,?,?,?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, serviceID)
	base.CheckErr(err)
}

//UpdateIncidentState update incident state to database
func (d *Store) UpdateIncidentState(incident *IncidentDb) {
	alert := incident.Alert
	tocheck := incident.ToCheck
	trigger := incident.Trigger
	serviceID := incident.ServiceID
	stmt, err := d.db.Prepare("UPDATE incidents set alert = ? , tocheck = ? , trigger = ? where serviceid = ?")
	base.CheckErr(err)
	_, err = stmt.Exec(alert, tocheck, trigger, serviceID)
	base.CheckErr(err)
}

//GetIncident get all arleady triggered incidents
func (d *Store) GetIncident() (inc []*IncidentDb) {
	var tmp string
	r, err := d.db.Query("select * from incidents inner join services on (incidents.serviceid=services.id);")
	base.CheckErr(err)
	for r.Next() {
		var incTmp IncidentDb
		err := r.Scan(&incTmp.ID, &incTmp.Title, &incTmp.ServiceID, &incTmp.CreateAt, &incTmp.Timer, &incTmp.Alert, &incTmp.ToCheck, &incTmp.Trigger, &tmp, &incTmp.ServiceName)
		base.CheckErr(err)
		inc = append(inc, &incTmp)
	}
	return
}
