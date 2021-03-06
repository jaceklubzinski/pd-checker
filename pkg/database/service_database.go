package database

import (
	"github.com/PagerDuty/go-pagerduty"
)

//Service structure for service stored in database
type Service struct {
	ID   string
	Name string
}

//SaveService insert service to database
func (d *Store) SaveService(service *pagerduty.Service) error {
	id := service.APIObject.ID
	name := service.APIObject.Summary
	stmt, err := d.db.Prepare("REPLACE INTO services values(?,?)")
	if err != nil {
		return err
	}
	sqlResult, err := stmt.Exec(id, name)
	if err != nil {
		return err
	}
	_, err = sqlResult.RowsAffected()
	return err
}

//GetService get all PagerDuty services without checker incidents
func (d *Store) GetService() (service []*Service, err error) {
	r, err := d.db.Query("select services.id,services.name FROM services left join incidents on (services.id = incidents.serviceid) where incidents.serviceid is NULL;")
	if err != nil {
		return nil, err
	}
	for r.Next() {
		var serviceTmp Service
		err := r.Scan(&serviceTmp.ID, &serviceTmp.Name)
		if err != nil {
			return nil, err
		}
		service = append(service, &serviceTmp)
	}
	return service, nil
}
