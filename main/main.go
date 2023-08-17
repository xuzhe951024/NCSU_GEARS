package main

import (
	"NCSU_Gears/common/constants"
	"NCSU_Gears/tasks/web/rest/v1"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	//tests.RunTestsForTask1And2()
	logrus.SetLevel(logrus.InfoLevel)
	viper.SetConfigName(constants.APPLICATION_CONFIG_NAME)
	viper.AddConfigPath(constants.RESOURCES_PATH)
	viper.SetConfigType(constants.CONFIG_FILE_TYPE)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error("Error reading config file, %s", err)
	}

	port := viper.Get(constants.PORT_CONFIG_NAME)
	webFramework := viper.Get(constants.WEB_FRAMEWORK)

	if webFramework == constants.GORILLA {
		v1.RegisterGorillaFunctionChainController(fmt.Sprintf(":%v", port))
	} else {
		v1.RegisterEchoFunctionChainController(fmt.Sprintf(":%v", port))
	}

}
