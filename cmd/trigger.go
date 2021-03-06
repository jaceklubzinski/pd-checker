package cmd

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/event"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger and instantly resolve single alert",
	Long:  `Trigger and instantly resolve single alert`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("trigger called")
		var opts pagerduty.V2Event
		integrationKey := viper.GetString("pagerduty_integration_key")
		opts.RoutingKey = integrationKey
		client := event.NewEvent(&opts)
		client.NewRecordMetricsEvent()
		client.SetPayLoad("24h")
		client.SetOptionsTrigger()
		err := client.Trigger()
		if err == nil {
			client.SetOptionsResolve()
			_ = client.Trigger()
		}
	},
}

func init() {
	eventCmd.AddCommand(triggerCmd)
}
