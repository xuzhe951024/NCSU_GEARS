package tests

import (
	"NCSU_Gears/tasks"
	"NCSU_Gears/utils"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func RunTestsForTask1And2() {
	pathToJson := "static/jsons/Task1.json"
	fnMappings, funcs := tasks.JsonToMaps(pathToJson)

	// Start the event handler goroutine
	go tasks.WarmStateUpdateEventHandler(funcs, fnMappings)

	var wg sync.WaitGroup
	for i, functionName := range funcs {
		wg.Add(1)
		go func(i int, functionName string) {
			defer wg.Done()

			// Generate event
			tasks.WarmStateUpdateChan <- true

			// Random delay
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(delay)

			res, err := tasks.ScheduleFunctionOnNode(functionName, fmt.Sprintf("data%d", i), fnMappings, funcs)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error scheduling function: %s", err))
				return
			}
			fmt.Println(fmt.Sprintf("GoRoutineId: %s Result for function: %v\n", utils.GetGoroutineID(), res.ResultsMap))
		}(i, functionName)
	}

	wg.Wait()

}
