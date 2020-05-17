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
	UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{}) error
	SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) error
	GetIncident() (inc []*Incident, err error)
	UpdateIncidentState(incident *Incident) error
}

//UpdateIncident update incident creation time and timer base on service id
func (d *Store) UpdateIncident(incident *pagerduty.Incident, incidentTimer interface{}) error {
	title := incident.Title
	id := incident.IncidentNumber
	serviceID := incident.Service.ID
	serviceName := incident.Service.Summary
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("Update incidents set id = ?,title = ?,createat = ?,timer = ? where serviceid = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, title, createAt, incidentTimer, serviceID)
	if err == nil {
		log.WithFields(log.Fields{
			"serviceName":    serviceName,
			"ServiceID":      serviceID,
			"IncidentNumber": id,
		}).Info("New incident for service registered")
	}
	return err

}

//SaveIncident insert incident to database
func (d *Store) SaveIncident(incident *pagerduty.Incident, incidentTimer interface{}) error {
	title := incident.Title
	id := incident.IncidentNumber
	serviceID := incident.Service.ID
	serviceName := incident.Service.Summary
	createAt := incident.CreatedAt
	stmt, err := d.db.Prepare("REPLACE INTO incidents(id,title,createat,timer,serviceid) values(?,?,?,?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, title, createAt, incidentTimer, serviceID)
	if err == nil {
		log.WithFields(log.Fields{
			"ServiceName":    serviceName,
			"ServiceID":      serviceID,
			"IncidentNumber": id,
			"CreatedAt":      createAt,
		}).Info("New incident for service registered")
	}
	return err
}

//UpdateIncidentState update incident state to database
func (d *Store) UpdateIncidentState(incident *Incident) error {
	alert := incident.Alert
	tocheck := incident.ToCheck
	trigger := incident.Trigger
	serviceID := incident.ServiceID
	stmt, err := d.db.Prepare("UPDATE incidents set alert = ? , tocheck = ? , trigger = ? where serviceid = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alert, tocheck, trigger, serviceID)
	return err
}

//GetIncident get all arleady triggered incidents
func (d *Store) GetIncident() (inc []*Incident, err error) {
	var ServiceID string
	r, err := d.db.Query("select * from incidents inner join services on (incidents.serviceid=services.id);")
	if err != nil {
		return nil, err
	}
	for r.Next() {
		var incTmp Incident
		err := r.Scan(&incTmp.ID, &incTmp.Title, &incTmp.ServiceID, &incTmp.CreateAt, &incTmp.Timer, &incTmp.Alert, &incTmp.ToCheck, &incTmp.Trigger, &ServiceID, &incTmp.ServiceName)
		if err != nil {
			return nil, err
		}
		inc = append(inc, &incTmp)
	}
	return inc, nil
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
