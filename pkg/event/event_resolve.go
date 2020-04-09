package event

//ResolveEvent set action as resolve for incident
func (o *ManageEvent) ResolveEvent() {
	o.Options.Action = "resolve"
	o.Options.DedupKey = o.Response.DedupKey
	o.message = "Alert Resolved"
}