package incident

import (
	"time"

	"github.com/jaceklubzinski/pd-checker/pkg/base"
	log "github.com/sirupsen/logrus"
)

type incidenter interface {
	MarkToCheck()
	SetAlertState()
}

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

//MarkToCheck check if service incident should be considered again base on last creation time and duraiton
func (i *Incident) MarkToCheck() error {
	serviceTimer, err := time.ParseDuration(i.Timer)
	if err != nil {
		return err
	}
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
	return err
}

//SetAlertState check if incident for service was already created
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
