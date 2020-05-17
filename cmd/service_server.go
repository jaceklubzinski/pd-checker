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
	Short: "Check for new Pagerduty event in a server mode",
	Long:  "Check for new Pagerduty event in a server mode and generate status page",
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
		incidents := incident.Manager{IncidentClient: conn}
		serviceClient := services.Services{Service: conn}
		ticker := time.NewTicker(checkEvery)
		go func() {
			for ; true; <-ticker.C {
				service := serviceClient.Service.Get()
				for _, s := range service.Services {
					err := DbRepository.SaveService(&s)
					base.CheckErr(err)
				}
				dbservice, err := DbRepository.GetService()
				base.CheckErr(err)
				for _, s := range dbservice {
					incidents.SetOptions()
					incidents.Options.ServiceIDs = []string{s.ID}
					if i := incidents.CheckForNew(); i != nil {
						repeatTimer := incidents.AlertDetails(i.Id)
						DbRepository.SaveIncident(i, repeatTimer)
					}
				}
				dbincident, err := DbRepository.GetIncident()
				for _, inc := range dbincident {
					inc.MarkToCheck()
					if inc.ToCheck == "Y" {
						incidents.SetOptionsFromIncident(inc)
						if i := incidents.CheckForNew(); i != nil {
							repeatTimer := incidents.AlertDetails(i.Id)
							inc.ToCheck = "N"
							inc.Trigger = "N"
							inc.Alert = "N"
							DbRepository.UpdateIncident(i, repeatTimer)
						}
						inc.SetAlertState()
						incidents.Alert(inc)
						DbRepository.UpdateIncidentState(inc)
					}
				}
				log.WithFields(log.Fields{
					"checkEvery": checkEvery,
				}).Info("Waiting to next check")
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
