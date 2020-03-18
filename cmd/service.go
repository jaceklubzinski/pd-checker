package cmd

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/config"
	"github.com/jaceklubzinski/pd-checker/pkg/database"
	"github.com/jaceklubzinski/pd-checker/pkg/incident"
	"github.com/jaceklubzinski/pd-checker/pkg/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service called")
		databasePath, err := cmd.Flags().GetString("database-path")
		base.CheckErr(err)
		config := config.NewConfig("", databasePath)
		db, err := database.ConnectDatabase(config)
		base.CheckErr(err)
		repository := database.NewIncidentRepository(db)
		repository.InitIncidentRepository()
		incidentRegiser := repository.GetIncident()
		fmt.Println(incidentRegiser)
		pdclient := pagerduty.NewClient(getFlagAuthToken())
		conn := client.NewApiClient(pdclient)
		serviceClient := services.Services{Service: conn}
		service := serviceClient.Service.ListServices()
		incident := incident.IncidentService{Incident: conn, Repository: repository}
		incident.IncidentOptions()
		for _, s := range service.Services {
			incident.Options.ServiceIDs = []string{s.APIObject.ID}
			incident.GetCheckerCount()
			incident.CheckerCount()
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	viper.AutomaticEnv()
	serverCmd.PersistentFlags().String("pagerduty_auth_token", "", "Set your PagerDuty auth Token")
}
