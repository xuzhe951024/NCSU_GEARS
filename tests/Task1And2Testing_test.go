package tests

import (
	"NCSU_Gears/tasks/web/service"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestForTask1And2(t *testing.T) {
	pathToJson := "static/jsons/Task1.json"
	fnMappings, funcs := service.JsonToMaps(pathToJson)

	// Start the event handler goroutine
	go service.WarmStateUpdateEventHandler(funcs, fnMappings)

	var wg sync.WaitGroup
	for i, functionName := range funcs {
		wg.Add(1)
		go func(i int, functionName string) {
			defer wg.Done()

			// Random delay
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(delay)

			_, err := service.ScheduleFunctionOnNode(functionName, fmt.Sprintf("data%d", i), fnMappings, funcs)
			if err != nil {
				t.Error(fmt.Sprintf("Error scheduling function: %s", err))
				return
			}

		}(i, functionName)
	}

	wg.Wait()

}
