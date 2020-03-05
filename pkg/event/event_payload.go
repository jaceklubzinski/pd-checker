package event

import "github.com/PagerDuty/go-pagerduty"

func (o *optionEvent) payLoad() {
	o.event.Payload = &pagerduty.V2Payload{
		Summary:  "PD CHECKER - OK",
		Severity: "info",
		Source:   "localhost",
	}
}
