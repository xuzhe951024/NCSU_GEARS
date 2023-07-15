package tasks

import (
	"NCSU_Gears/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func JsonToMaps(pathToJson string) (map[string]models.Function, []string) {
	jsonFile, err := os.Open(pathToJson)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string][]models.Function
	json.Unmarshal(byteValue, &result)

	funcs := make([]string, len(result["functions"]))
	fnMappings := make(map[string]models.Function)

	for i, function := range result["functions"] {
		funcs[i] = function.Name
		fnMappings[function.Name] = function
	}

	return fnMappings, funcs
}
