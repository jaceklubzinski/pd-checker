package incident

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarkToCheck(t *testing.T) {
	i := Incident{
		Timer:    "24h",
		CreateAt: "2020-05-17T20:21:27Z",
		ToCheck:  "N",
	}
	err := i.MarkToCheck()
	assert.Equal(t, i.ToCheck, "Y")
	assert.Nil(t, err)
	i.CreateAt = time.Now().Format("2006-01-02T15:04:05Z")
	i.ToCheck = "N"
	err = i.MarkToCheck()
	assert.Equal(t, i.ToCheck, "N")
	assert.Nil(t, err)
}

func TestSetAlertState(t *testing.T) {
	i := Incident{
		ToCheck: "Y",
		Trigger: "N",
		Alert:   "N",
	}
	i.SetAlertState()
	assert.Equal(t, i.Alert, "Y")
	i.Trigger = "Y"
	i.Alert = "N"
	i.SetAlertState()
	assert.Equal(t, i.Alert, "N")
}
