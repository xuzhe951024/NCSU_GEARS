package tests

import (
	"NCSU_Gears/common/constants"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestRegisterFunctionChain(t *testing.T) {
	viper.SetConfigName(constants.APPLICATION_CONFIG_NAME)
	viper.AddConfigPath("." + constants.RESOURCES_PATH)
	viper.SetConfigType(constants.CONFIG_FILE_TYPE)
	if err := viper.ReadInConfig(); err != nil {
		t.Error("Error reading config file", err)
	}

	port := viper.Get(constants.PORT_CONFIG_NAME)

	pathToJson := "static/jsons/Task1.json"

	jsonFile, err := os.Open(pathToJson)
	if err != nil {
		t.Error(err)
	}
	defer jsonFile.Close()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Post(fmt.Sprintf("https://localhost:%v%s", port, constants.URI_V1_REGISTER_FUNCTION_CHAIN),
		constants.CONTANT_TYPE_JSON,
		jsonFile)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	defer resp.Body.Close()

	logrus.Info("Response status:", resp.Status)
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Error reading response body:", err)
	}
	logrus.Info("Response body:", string(rspBody))

	assert.Equal(t, "200 OK", resp.Status)
}
