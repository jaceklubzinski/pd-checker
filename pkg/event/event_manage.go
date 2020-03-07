package event

import (
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/avast/retry-go"
)

func (o *ManageEvent) manageIncident() {
	fmt.Println("create incident")
	err := retry.Do(
		func() error {
			resp, err := pagerduty.ManageEvent(*o.Options)
			if err != nil {
				return err
			}
			o.Response = resp
			fmt.Println(resp)
			return nil
		},
		retry.Delay(time.Second*5),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: %s\n", n, err)
		}),
	)
	if err != nil {
		panic(err)
	}
}
