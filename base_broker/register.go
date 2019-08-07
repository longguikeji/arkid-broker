package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	pb "github.com/longguikeji/arkid-broker/protos/apistore"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func readContent(filePath, serverURI, instanceID, instanceName string) ([]byte, error) {
	var content []byte
	rawContent, _ := ioutil.ReadFile(filePath)
	jsonContent, err := yaml.YAMLToJSON(rawContent)
	if err != nil {
		jsonContent = rawContent
	}
	api := map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonContent), &api)
	if err != nil {
		return content, err
	}

	serverInfo := map[string]string{
		"instanceName": instanceName,
		"instanceID":   instanceID,
	}
	description, _ := json.Marshal(serverInfo)
	api["servers"] = []map[string]string{{
		"url":         serverURI,
		"description": string(description),
	}}

	content, err = json.Marshal(api)
	return content, err
}

// RegisterAPI OpenAPI to API-SVR
func RegisterAPI(filePath, apiSVRURI, serverURI, instanceID, instanceName string) error {
	content, err := readContent(filePath, serverURI, instanceID, instanceName)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(map[string]string{"content": string(content)})

	req, err := http.NewRequest(
		"POST",
		apiSVRURI+"/v1/apis",
		bytes.NewReader(data),
	)
	if err != err {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	res := &pb.CreateAPIRes{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	if res.Code != "0" {
		return errors.New(res.Msg)
	}

	return nil
}

// NewRegisterCmd ...
func NewRegisterCmd(config *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register ...",
		Run: func(cmd *cobra.Command, args []string) {
			file := config.GetString("file")
			target := config.GetString("target")
			server := config.GetString("server")
			instanceID := config.GetString("instanceID")
			instanceName := config.GetString("instanceName")

			fmt.Println("load openapi from", file)
			fmt.Println("register to", target)
			fmt.Println("access serve by", server)
			fmt.Println("instanceID:", instanceID)
			fmt.Println("instanceName:", instanceName)

			err := RegisterAPI(file, target, server, instanceID, instanceName)
			if err != nil {
				fmt.Println("register failed")
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("register successfully")
		},
	}
	cmd.Flags().StringP("file", "f", "openapi", "the path of openapi")
	cmd.Flags().StringP("target", "t", "", "the uri of api-svr")
	cmd.Flags().StringP("server", "i", "", "the uri of broker self")
	cmd.Flags().StringP("instanceID", "u", "", "the ID of api instance")
	cmd.Flags().StringP("instanceName", "n", "", "the human-readable name of api instance")

	return cmd
}

// BindFlag ...
func BindFlag(config *viper.Viper, cmd *cobra.Command) {
	config.BindPFlag("file", cmd.Flags().Lookup("file"))
	config.BindPFlag("target", cmd.Flags().Lookup("target"))
	config.BindPFlag("server", cmd.Flags().Lookup("server"))
	config.BindPFlag("instanceID", cmd.Flags().Lookup("instanceID"))
	config.BindPFlag("instanceName", cmd.Flags().Lookup("instanceName"))
}
