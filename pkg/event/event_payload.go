package event

import "github.com/PagerDuty/go-pagerduty"

func (o *ManageEvent) payLoad() {
	o.Options.Payload = &pagerduty.V2Payload{
		Summary:  "PD CHECKER - OK",
		Severity: "info",
		Source:   "localhost",
	}
}
