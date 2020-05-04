package database

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

//Incident structure for incidents stored in database
type Incident struct {
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
	GetIncident() (inc []*Incident)
	UpdateIncidentState(incident *Incident)
}

//UpdateIncident update incident creation time and timer base on service id
func (d *Store) UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	serviceID := incident.Service.ID
	serviceName := incident.Service.Summary
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("Update incidents set id = ?,title = ?,createat = ?,timer = ? where serviceid = ?")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, serviceID)
	base.CheckErr(err)
	log.WithFields(log.Fields{
		"serviceName":    serviceName,
		"ServiceID":      serviceID,
		"IncidentNumber": id,
	}).Info("New incident for service registered")
}

//SaveIncident insert incident to database
func (d *Store) SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) {
	title := incident.Title
	id := incident.IncidentNumber
	serviceID := incident.Service.ID
	serviceName := incident.Service.Summary
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("REPLACE INTO incidents(id,title,createat,timer,serviceid) values(?,?,?,?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, serviceID)
	base.CheckErr(err)
	log.WithFields(log.Fields{
		"ServiceName":    serviceName,
		"ServiceID":      serviceID,
		"IncidentNumber": id,
		"CreatedAt":      createAt,
	}).Info("New incident for service registered")
}

//UpdateIncidentState update incident state to database
func (d *Store) UpdateIncidentState(incident *Incident) {
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
func (d *Store) GetIncident() (inc []*Incident) {
	var tmp string
	r, err := d.db.Query("select * from incidents inner join services on (incidents.serviceid=services.id);")
	base.CheckErr(err)
	for r.Next() {
		var incTmp Incident
		err := r.Scan(&incTmp.ID, &incTmp.Title, &incTmp.ServiceID, &incTmp.CreateAt, &incTmp.Timer, &incTmp.Alert, &incTmp.ToCheck, &incTmp.Trigger, &tmp, &incTmp.ServiceName)
		base.CheckErr(err)
		inc = append(inc, &incTmp)
	}
	return
}

//MarkToCheck check if service incident should be considered again base on last creation time and duraiton
func (i *Incident) MarkToCheck() {
	serviceTimer, err := time.ParseDuration(i.Timer)
	base.CheckErr(err)
	lastTillNow := base.LastTillNowDuration(i.CreateAt)
	if lastTillNow > serviceTimer {
		i.ToCheck = "Y"
		log.WithFields(log.Fields{
			"ServiceName":  i.ServiceName,
			"ServiceID":    i.ServiceID,
			"LastCreateAt": i.CreateAt,
		}).Info("Checking for new alert")
	} else {
		log.WithFields(log.Fields{
			"ServiceName": i.ServiceName,
			"ServiceID":   i.ServiceID,
		}).Info("Service has working PagerDuty integration")
	}
}

//SetAlertState check if service incident was created
func (i *Incident) SetAlertState() {
	if i.ToCheck == "Y" && i.Trigger == "N" {
		i.Alert = "Y"
		log.WithFields(log.Fields{
			"ServiceName": i.ServiceName,
			"ServiceID":   i.ServiceID,
		}).Warning("New alert for service will be created")
	} else if i.ToCheck == "Y" && i.Trigger == "Y" {
		log.WithFields(log.Fields{
			"ServiceName": i.ServiceName,
			"ServiceID":   i.ServiceID,
		}).Warning("Alert for service already created")
	}
}
