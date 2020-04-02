package database

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
)

//ServiceDB structure for service stored in database
type ServiceDB struct {
	ID   string
	Name string
}

//SaveService insert service to database
func (d *Store) SaveService(service *pagerduty.Service) {
	id := service.APIObject.ID
	name := service.APIObject.Summary
	stmt, err := d.db.Prepare("REPLACE INTO services values(?,?)")
	base.CheckErr(err)
	_, err = stmt.Exec(id, name)
	base.CheckErr(err)
}

//GetService get all PagerDuty services without checker incidents
func (d *Store) GetService() (service []*ServiceDB) {
	r, err := d.db.Query("select services.id,services.name FROM services left join incidents on (services.id = incidents.serviceid) where incidents.serviceid is NULL;")
	base.CheckErr(err)
	for r.Next() {
		var serviceTmp ServiceDB
		err := r.Scan(&serviceTmp.ID, &serviceTmp.Name)
		base.CheckErr(err)
		service = append(service, &serviceTmp)
	}
	return
}
