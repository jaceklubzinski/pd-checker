package cmd

import (
	"github.com/spf13/viper"
)

func getFlagAuthToken() string {
	return viper.GetString("pagerduty_auth_token")
}
