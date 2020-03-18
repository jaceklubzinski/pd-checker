package cmd

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/event"
	"github.com/jaceklubzinski/pd-checker/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
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
		clientEvent.EventMetrics.NewRecordMetricsEvent()
		clientEvent.PayLoad(triggerEvery.String())
		go func() {
			ticker := time.NewTicker(triggerEvery)
			for ; true; <-ticker.C {
				clientEvent.TriggerEvent()
				err := clientEvent.ManageIncident()
				if err == nil {
					clientEvent.ResolveEvent()
					_ = clientEvent.ManageIncident()
				}
			}
		}()
		metrics.Server()
	},
}

func init() {
	eventCmd.AddCommand(serverCmd)
	defaultRepeat := 60 * time.Second
	serverCmd.Flags().DurationP("repeat", "r", defaultRepeat, "Trigger new alert every duration minutes")
}
