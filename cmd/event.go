package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(eventCmd)
	viper.BindEnv("pagerduty_integration_key")
	eventCmd.Flags().StringP("integration-key", "i", viper.GetString("PAGERDUTY_INTEGRATION_KEY"), "Integration key for PagerDuty event api (env PAGERDUTY_INTEGRATION_KEY)")
}
