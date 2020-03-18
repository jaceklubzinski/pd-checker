package services

import (
	"fmt"

	"github.com/jaceklubzinski/pd-checker/pkg/client"
)

//Services PagerDuty service struct
type Services struct {
	Service client.ServiceClient
}

//GetAll Get all services
func (s *Services) GetAll() {
	getAll := s.Service.ListServices()
	for _, p := range getAll.Services {
		fmt.Printf("ID: %s Name: %s\n", p.APIObject.ID, p.Name)
	}
}
