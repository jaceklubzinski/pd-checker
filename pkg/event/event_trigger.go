package event

func (o *ManageEvent) TriggerEvent() {
	o.Options.Action = "trigger"
	o.payLoad()
	o.manageIncident()
}
