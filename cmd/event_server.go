package cmd

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/event"
	"github.com/jaceklubzinski/pd-checker/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var eventServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Trigger and resolve alerts in deamon mode",
	Long: `Server mode trigger and resolve alerts in repetable mode.
Metrics are available on url 127.0.0.1:2112/metrics`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		var opts pagerduty.V2Event
		integrationKey := viper.GetString("pagerduty_integration_key")
		triggerEvery, _ := cmd.Flags().GetDuration("repeat")
		opts.RoutingKey = integrationKey
		clientEvent := event.NewEvent(&opts)
		clientEvent.NewRecordMetricsEvent()
		clientEvent.SetPayLoad(triggerEvery.String())
		log.WithFields(log.Fields{
			"triggerEvery": triggerEvery,
		}).Info("Waitig to trigger next event")
		go func() {
			ticker := time.NewTicker(triggerEvery)
			for ; true; <-ticker.C {
				clientEvent.SetOptionsTrigger()
				err := clientEvent.Trigger()
				if err == nil {
					clientEvent.SetOptionsResolve()
					_ = clientEvent.Trigger()
				}
			}
		}()
		metrics.Server()
	},
}

func init() {
	eventCmd.AddCommand(eventServerCmd)
	defaultRepeat := 60 * time.Second
	eventServerCmd.Flags().DurationP("repeat", "r", defaultRepeat, "Trigger new alert every duration minutes")
}
