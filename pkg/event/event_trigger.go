package event

//TriggerEvent set action as trigger for incident
func (o *ManageEvent) TriggerEvent() {
	o.Options.Action = "trigger"
	o.message = "Alert Triggered"
}
