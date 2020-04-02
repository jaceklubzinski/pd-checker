package cmd

import (
	"fmt"
	"log"
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
		DbRepository := database.NewIncidentRepository(db)
		DbRepository.InitIncidentRepository()
		pdclient := pagerduty.NewClient(getFlagAuthToken())
		conn := client.NewApiClient(pdclient)
		incidents := incident.IncidentService{IncidentClient: conn, DbRepository: DbRepository}
		incidents.IncidentOptions()
		serviceClient := services.Services{Service: conn}
		ticker := time.NewTicker(triggerEvery)
		for ; true; <-ticker.C {
			service := serviceClient.Service.ListServices()
			for _, s := range service.Services {
				DbRepository.SaveService(&s)
			}
			incidents.Services = DbRepository.GetService()
			incidents.Incidents = DbRepository.GetIncident()
			incidents.CheckServiceIncident()
			incidents.MarkToCheck()
			incidents.CheckToAlert()
			incidents.Alert()
			log.Printf("Waitig for %s to next check", triggerEvery)
		}

	},
}

func init() {
	serviceCmd.AddCommand(serviceServerCmd)
	defaultRepeatService := 60 * time.Second
	serviceServerCmd.Flags().DurationP("check-repeat", "t", defaultRepeatService, "Check for new new alert every duration minutes")
}
