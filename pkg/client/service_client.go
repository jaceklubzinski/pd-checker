package client

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
)

type ServiceClient interface {
	ListServices() *pagerduty.ListServiceResponse
}

var serviceOpts pagerduty.ListServiceOptions

func (c *ApiClient) ListServices() *pagerduty.ListServiceResponse {
	eps, err := c.client.ListServices(serviceOpts)
	base.CheckErr(err)
	return eps
}
