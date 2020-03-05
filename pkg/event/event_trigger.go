package event

func (o *optionEvent) triggerEvent() {
	o.options.Action = "trigger"
}
