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
	"github.com/jaceklubzinski/pd-checker/pkg/ui"

	log "github.com/sirupsen/logrus"
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
		checkEvery, _ := cmd.Flags().GetDuration("check-repeat")
		base.CheckErr(err)
		config := config.NewConfig("", databasePath)
		db, err := database.ConnectDatabase(config)
		base.CheckErr(err)
		DbRepository := database.NewIncidentRepository(db)
		DbRepository.InitIncidentRepository()
		server := ui.NewServer(DbRepository)
		pdclient := pagerduty.NewClient(getFlagAuthToken())
		conn := client.NewApiClient(pdclient)
		incidents := incident.IncidentService{IncidentClient: conn, DbRepository: DbRepository}
		incidents.IncidentOptions()
		serviceClient := services.Services{Service: conn}
		ticker := time.NewTicker(checkEvery)
		go func() {
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
				log.WithFields(log.Fields{
					"checkEvery": checkEvery,
				}).Info("Waitig to next check")
			}
		}()
		server.Listen()
	},
}

func init() {
	serviceCmd.AddCommand(serviceServerCmd)
	defaultRepeatService := 60 * time.Second
	serviceServerCmd.Flags().DurationP("check-repeat", "t", defaultRepeatService, "Check for new new alert every duration minutes")
}
