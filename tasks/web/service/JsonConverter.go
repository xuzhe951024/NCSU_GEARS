package service

import (
	"NCSU_Gears/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
)

func parseJSON(pathToJson string) models.RegisterFunctionChainVO {
	jsonFile, err := os.Open(pathToJson)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result models.RegisterFunctionChainVO
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		logrus.Error("Unmarshal jason failed: ", err)
		return models.RegisterFunctionChainVO{}
	}

	return result
}

func detectCycle(functions []models.UnparsedFunction) bool {
	colors := make(map[string]string)
	for _, function := range functions {
		colors[function.Name] = "white"
	}

	for _, function := range functions {
		if colors[function.Name] == "white" {
			if dfs(function.Name, functions, colors) {
				return true
			}
		}
	}

	return false
}

func dfs(node string, functions []models.UnparsedFunction, colors map[string]string) bool {
	colors[node] = "gray"

	for _, function := range functions {
		if function.Name == node {
			for _, next := range function.Next {
				if colors[next.Name] == "gray" {
					return true
				}
				if colors[next.Name] == "white" && dfs(next.Name, functions, colors) {
					return true
				}
			}
		}
	}

	colors[node] = "black"
	return false
}

func convertToFunction(functions []models.UnparsedFunction) map[string]models.Function {
	fnMappings := make(map[string]models.Function)

	// First, create all Function structs and add them to fnMappings.
	for _, unparsedFunction := range functions {
		function := models.Function{
			Name:            unparsedFunction.Name,
			Version:         unparsedFunction.Version,
			DependsOn:       make(map[string]models.FunctionIndex),
			Next:            make(map[string]struct{ Name string }),
			Timeout:         unparsedFunction.Timeout,
			IsLast:          unparsedFunction.IsLast,
			BreakConditions: unparsedFunction.BreakConditions,
			Data:            unparsedFunction.Data,
			IsWarm:          unparsedFunction.IsWarm,
		}
		fnMappings[unparsedFunction.Name] = function
	}

	// Then, fill in the DependsOn and Next fields.
	for _, unparsedFunction := range functions {
		function := fnMappings[unparsedFunction.Name]

		for _, dep := range unparsedFunction.DependsOn {
			if _, ok := fnMappings[dep.Name]; ok {
				fnMappings[unparsedFunction.Name].DependsOn[dep.Name] = dep
			}
		}

		for _, next := range unparsedFunction.Next {
			if _, ok := fnMappings[next.Name]; ok {
				fnMappings[unparsedFunction.Name].Next[next.Name] = struct{ Name string }{Name: next.Name}
			}
		}

		fnMappings[unparsedFunction.Name] = function
	}

	return fnMappings
}

func ParseJsonToMaps(unparsedFunctions models.RegisterFunctionChainVO) (map[string]models.Function, []string) {
	funcs := make([]string, len(unparsedFunctions.Functions))
	unparsedFnMappings := make(map[string]models.UnparsedFunction)

	for i, function := range unparsedFunctions.Functions {
		funcs[i] = function.Name
		unparsedFnMappings[function.Name] = function
	}

	if detectCycle(unparsedFunctions.Functions) {
		log.Fatal("Detected a cycle in the function graph.")
	}

	fnMappings := convertToFunction(unparsedFunctions.Functions)

	return fnMappings, funcs
}

func JsonToMaps(pathToJson string) (map[string]models.Function, []string) {
	result := parseJSON(pathToJson)

	return ParseJsonToMaps(result)
}
