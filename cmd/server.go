package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/event"
	"github.com/jaceklubzinski/pd-checker/pkg/metrics"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		var opts pagerduty.V2Event
		integrationKey := os.Getenv("PAGERDUTY_INTEGRATION_KEY")
		opts.RoutingKey = integrationKey
		client := &event.ManageEvent{Options: &opts}
		client.EventMetrics.NewRecordMetricsEvent()
		//go base.RepeatFunction(client.TriggerEvent())
		client.TriggerEvent()
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			for ; true; <-ticker.C {
				_ = client.ManageIncident()
			}
		}()
		metrics.MetricsServer()
	},
}

func init() {
	eventCmd.AddCommand(serverCmd)
}
