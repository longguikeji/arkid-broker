package cmd

import (
	"github.com/rockl2e/ark-apisvr/broker"
	"github.com/spf13/viper"
)

var config = viper.New()

var registerCmd = broker.NewRegisterCmd(config)

func init() {
	broker.BindFlag(config, registerCmd)
	rootCmd.AddCommand(registerCmd)
}
