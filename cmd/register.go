package cmd

import (
	broker "github.com/longguikeji/arkid-broker/base_broker"
	"github.com/spf13/viper"
)

var config = viper.New()

var registerCmd = broker.NewRegisterCmd(config)

func init() {
	broker.BindFlag(config, registerCmd)
	rootCmd.AddCommand(registerCmd)
}
