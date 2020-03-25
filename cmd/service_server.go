package cmd

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-checker/pkg/base"
	"github.com/jaceklubzinski/pd-checker/pkg/client"
	"github.com/jaceklubzinski/pd-checker/pkg/config"
	"github.com/jaceklubzinski/pd-checker/pkg/database"
	"github.com/jaceklubzinski/pd-checker/pkg/incident"
	"github.com/jaceklubzinski/pd-checker/pkg/services"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serviceServerCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		databasePath, err := cmd.Flags().GetString("database-path")
		base.CheckErr(err)
		triggerEvery, _ := cmd.Flags().GetDuration("check-repeat")
		base.CheckErr(err)
		config := config.NewConfig("", databasePath)
		db, err := database.ConnectDatabase(config)
		base.CheckErr(err)
		repository := database.NewIncidentRepository(db)
		repository.InitIncidentRepository()
		pdclient := pagerduty.NewClient(getFlagAuthToken())
		conn := client.NewApiClient(pdclient)
		serviceClient := services.Services{Service: conn}
		service := serviceClient.Service.ListServices()
		incidents := incident.IncidentService{IncidentClient: conn, Repository: repository}
		incidents.IncidentOptions()
		for _, s := range service.Services {
			incidents.Options.ServiceIDs = []string{s.APIObject.ID}
			incidents.CountIncidentService()
			incidents.CounterInfo()
		}
		//server := incident.NewServer(&incidents, repository)
		ticker := time.NewTicker(triggerEvery)
		for ; true; <-ticker.C {
			incidents.Incidents = repository.GetIncident()
			for _, v := range incidents.Incidents {
				fmt.Println(v.Service)
			}
		}

	},
}

func init() {
	serviceCmd.AddCommand(serviceServerCmd)
	defaultRepeatService := 60 * time.Second
	serviceServerCmd.Flags().DurationP("check-repeat", "t", defaultRepeatService, "Check for new new alert every duration minutes")
}
