package event

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/avast/retry-go"

	log "github.com/sirupsen/logrus"
)

//Trigger create or resolve PagerDuty event and retry on failure
func (e *Event) Trigger() error {
	var delaySecond time.Duration = 5
	var retryAttempt uint = 5
	err := retry.Do(
		func() error {
			resp, err := e.manager.ManageEvent(e.Options)
			if err != nil {
				return err
			}
			e.Response = resp
			e.RecordMetricsEventOk(e.Response.Status)
			log.Info(e.message)
			return nil
		},
		retry.Attempts(retryAttempt),
		retry.Delay(time.Second*delaySecond),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: %s\n", n, err)
			e.RecordMetricsEventError(err.Error())
		}),
	)
	return err
}

//SetOptionsResolve set action as resolve for incident
func (e *Event) SetOptionsResolve() {
	e.Options.Action = "resolve"
	e.Options.DedupKey = e.Response.DedupKey
	e.message = "Alert Resolved"
}

//SetOptionsTrigger set action as trigger for incident
func (e *Event) SetOptionsTrigger() {
	e.Options.Action = "trigger"
	e.message = "Alert Triggered"
}

//SetPayLoad add additional information to new alert
func (e *Event) SetPayLoad(triggerEvery string) {
	e.Options.Payload = &pagerduty.V2Payload{
		Summary:  "PD CHECKER - OK",
		Severity: "info",
		Source:   "localhost",
		Details:  triggerEvery,
	}
}
