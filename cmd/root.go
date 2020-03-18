package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pd-checker",
	Short: "PagerDuty alert checker",
	Long:  `Program trigger and insantly resolved incident in PagerDuty`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("metrics-port", "p", "2112", "Port for metrics endpoint")
	rootCmd.PersistentFlags().StringP("database-path", "d", "./incidentRepository.db", "sqlite database path")
}
