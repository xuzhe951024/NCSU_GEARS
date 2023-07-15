package tests

import (
	"NCSU_Gears/tasks"
	"NCSU_Gears/utils"
	"fmt"
)

func RunTestsForTask1And2() {
	pathToJson := "static/jsons/Task1.json"
	fnMappings, funcs := tasks.JsonToMaps(pathToJson)
	for i, functionName := range funcs {
		res, err := tasks.ScheduleFunctionOnNode(functionName, fmt.Sprintf("data%d", i), fnMappings, funcs)
		if err != nil {
			return
		}
		fmt.Println(fmt.Sprintf("GoRoutineId: %s Result for function: %v\n", utils.GetGoroutineID(), res.ResultsMap))
	}
}
