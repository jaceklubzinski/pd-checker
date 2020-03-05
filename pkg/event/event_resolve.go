package event

func (o *manageEvent) resolveEvent() {
	o.options.event.Action = "resolve"
	o.options.event.DedupKey = o.response.DedupKey

}
