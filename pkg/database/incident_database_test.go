package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/incident"

	"github.com/stretchr/testify/assert"
)

func TestUpdateIncident(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbRepository := NewIncidentRepository(db)
	mock.ExpectPrepare("Update incidents set").ExpectExec().
		WithArgs(12345, "PD CHECKER - TEST", "2020-05-12T19:25:11Z", "1s", "PK6TEST1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	i := &pagerduty.Incident{
		IncidentNumber: 12345,
		Title:          "PD CHECKER - TEST",
		CreatedAt:      "2020-05-12T19:25:11Z",
		Service: pagerduty.APIObject{
			ID:      "PK6TEST1",
			Summary: "test-service1",
		},
	}
	repeatTimer := "1s"
	err = dbRepository.UpdateIncident(i, repeatTimer)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when saving data to database", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveIncident(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbRepository := NewIncidentRepository(db)
	mock.ExpectPrepare("REPLACE INTO incidents").ExpectExec().
		WithArgs(12345, "PD CHECKER - TEST", "2020-05-12T19:25:11Z", "1s", "PK6TEST1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	i := &pagerduty.Incident{
		IncidentNumber: 12345,
		Title:          "PD CHECKER - TEST",
		CreatedAt:      "2020-05-12T19:25:11Z",
		Service: pagerduty.APIObject{
			ID:      "PK6TEST1",
			Summary: "test-service1",
		},
	}
	repeatTimer := "1s"
	err = dbRepository.SaveIncident(i, repeatTimer)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when saving data to database", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateIncidentState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbRepository := NewIncidentRepository(db)
	mock.ExpectPrepare("UPDATE incidents set").ExpectExec().
		WithArgs("Y", "Y", "Y", "PK6TEST1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	i := &incident.Incident{
		Alert:     "Y",
		ToCheck:   "Y",
		Trigger:   "Y",
		ServiceID: "PK6TEST1",
	}
	err = dbRepository.UpdateIncidentState(i)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when update incident state to database", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetIncident(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbRepository := NewIncidentRepository(db)

	rows := sqlmock.NewRows([]string{"ID", "Title", "ServiceID", "CreateAt", "Timer", "Alert", "ToCheck", "Trigger", "ServiceIDtmp", "ServiceName"}).
		AddRow("12345", "PD CHECKER - TEST", "PK6TEST1", "2020-05-12T19:25:11Z", "1s", "Y", "Y", "Y", "PK6TEST1", "test-service1")
	mock.ExpectQuery("select (.+) from incidents inner join services").
		WillReturnRows(rows)

	incidents, err := dbRepository.GetIncident()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when select incidents", err)
	}
	for _, i := range incidents {
		assert.Equal(t, i.ID, "12345")
		assert.Equal(t, i.ServiceName, "test-service1")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
