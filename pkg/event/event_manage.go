package event

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/avast/retry-go"
	log "github.com/sirupsen/logrus"
)

//ManageIncident create or resolve incident
func (o *ManageEvent) ManageIncident() error {
	var delaySecond time.Duration = 5
	var retryAttempt uint = 5
	err := retry.Do(
		func() error {
			resp, err := pagerduty.ManageEvent(*o.Options)
			if err != nil {
				return err
			}
			o.Response = resp
			o.EventMetrics.RecordMetricsEvent(o.Response.Status)
			log.Info(o.message)
			return nil
		},
		retry.Attempts(retryAttempt),
		retry.Delay(time.Second*delaySecond),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: %s\n", n, err)
			o.EventMetrics.RecordMetricsEventError(err.Error())
		}),
	)
	return err
}
