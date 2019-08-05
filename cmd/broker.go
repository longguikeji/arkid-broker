package cmd

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/longguikeji/arkid-broker/router"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile   string
	brokerName   string
	brokerRouter *mux.Router
)

func init() {
	cobra.OnInitialize(initConfig)

	brokerCmd.Flags().StringVarP(&configFile, "config", "f", "", "the config for broker")
	brokerCmd.Flags().IntP("port", "p", 0, "the port listen to")
	brokerCmd.Flags().StringP("target", "t", "", "the target addr broker to")

	viper.BindPFlag("port", brokerCmd.Flags().Lookup("port"))
	viper.BindPFlag("target", brokerCmd.Flags().Lookup("target"))

	rootCmd.AddCommand(brokerCmd)
}

func initConfig() {
	if configFile == "" {
		configFile = fmt.Sprintf("config/config.yaml")
	}
	brokerRouter = router.GetRouter()
	viper.SetConfigFile(configFile)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "broker ...",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load config from", configFile)
		fmt.Println("listen on", viper.GetInt("port"))
		fmt.Println("forward request to", viper.GetString("target"))

		http.Handle("/", brokerRouter)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", viper.GetInt("port")), nil))
	},
}
