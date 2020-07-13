package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/services"
	"github.com/stretchr/testify/assert"
)

type mockApiClient struct{}

type mockIncidentClient interface {
	Get() *pagerduty.ListServiceResponse
}

func (c *mockApiClient) Get() *pagerduty.ListServiceResponse {
	return &pagerduty.ListServiceResponse{
		Services: []pagerduty.Service{
			pagerduty.Service{
				APIObject: pagerduty.APIObject{
					ID:      "PK6TEST1",
					Summary: "test-service1",
				},
			},
		},
	}
}

func TestSaveService(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbRepository := NewIncidentRepository(db)

	mock.ExpectPrepare("REPLACE INTO services").ExpectExec().
		WithArgs("PK6TEST1", "test-service1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	client := &mockApiClient{}
	serviceClient := services.Services{Service: client}
	service := serviceClient.Service.Get()
	for _, s := range service.Services {
		assert.Equal(t, s.APIObject.ID, "PK6TEST1")
		assert.Equal(t, s.APIObject.Summary, "test-service1")
		err := dbRepository.SaveService(&s)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when saving data to database", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetService(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbRepository := NewIncidentRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("PK6TEST1", "test-service1")
	mock.ExpectQuery("select (.+) FROM services*").
		WillReturnRows(rows)

	services, err := dbRepository.GetService()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when select services", err)
	}
	for _, s := range services {
		assert.Equal(t, s.ID, "PK6TEST1")
		assert.Equal(t, s.Name, "test-service1")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
