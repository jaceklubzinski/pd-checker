package services

import (
	"github.com/jaceklubzinski/pd-checker/pkg/client"
)

//Services PagerDuty service struct
type Services struct {
	Service client.ServiceClient
}
