package cmd

import (
	"github.com/spf13/viper"
)

func getFlagAuthToken() string {
	return viper.GetString("pagerduty_auth_token")
}

func getFlagMetricsPort() string {
	return viper.GetString("metrics-port")
}

func getFlagDatabasePath() string {
	return viper.GetString("database-path")
}
