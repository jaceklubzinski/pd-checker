package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
)

type ServiceClient interface {
	Get() *pagerduty.ListServiceResponse
}

var serviceOpts pagerduty.ListServiceOptions

//Get all available serices from PagerDuty
func (c *ApiClient) Get() *pagerduty.ListServiceResponse {
	eps, err := c.client.ListServices(serviceOpts)
	base.CheckErr(err)
	return eps
}
