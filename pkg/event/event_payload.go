package event

import (
	"github.com/PagerDuty/go-pagerduty"
)

//PayLoad add additional information to new alert
func (o *ManageEvent) PayLoad(triggerEvery string) {
	o.Options.Payload = &pagerduty.V2Payload{
		Summary:  "PD CHECKER - OK",
		Severity: "info",
		Source:   "localhost",
		Details:  triggerEvery,
	}
}
