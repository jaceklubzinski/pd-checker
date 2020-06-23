package event

import (
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
)

type mockEvent struct {
}

type mockEventManager interface {
	ManageEvent(options *pagerduty.V2Event) (*pagerduty.V2EventResponse, error)
}

func (m *mockEvent) ManageEvent(options *pagerduty.V2Event) (*pagerduty.V2EventResponse, error) {
	response := pagerduty.V2EventResponse{
		Status: "success",
	}
	return &response, nil
}
func TestManageIncident(t *testing.T) {
	client := &mockEvent{}
	var opts pagerduty.V2Event
	clientEvent := Event{
		Options: &opts,
		manager: client,
	}
	clientEvent.NewRecordMetricsEvent()
	clientEvent.SetPayLoad("24h")
	clientEvent.SetOptionsTrigger()
	assert.Equal(t, clientEvent.Options.Action, "trigger")
	err := clientEvent.Trigger()
	if err == nil {
		clientEvent.SetOptionsResolve()
		_ = clientEvent.Trigger()
	}
	assert.Equal(t, clientEvent.Response.Status, "success")
	assert.Equal(t, clientEvent.Options.Payload.Details, "24h")
	assert.Equal(t, clientEvent.Options.Action, "resolve")

}
