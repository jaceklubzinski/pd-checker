package cmd

import (
	"fmt"
	"os"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/event"
	"github.com/spf13/cobra"
)

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("trigger called")
		var opts pagerduty.V2Event
		integrationKey := os.Getenv("PAGERDUTY_INTEGRATION_KEY")
		opts.RoutingKey = integrationKey
		client := &event.ManageEvent{Options: &opts}
		client.TriggerEvent()
		client.ResolveEvent()
	},
}

func init() {
	eventCmd.AddCommand(triggerCmd)
}
