package event

import "github.com/PagerDuty/go-pagerduty"

func (o *manageEvent) createIncident() {
	o.response, err := pagerduty.ManageEvent(*o.options.event)
	if err != nil {
		panic(err)
	}
}
