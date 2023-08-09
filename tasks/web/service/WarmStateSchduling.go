package service

import (
	"NCSU_Gears/common/utils"
	"NCSU_Gears/models"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

// Maintain status of each identified function
var funcRunStatus = map[string]string{}

var resultsMap = make(map[string]interface{})

// Maximum number of warm state functions
var numParallel = 5

// Warm state functions
var warmFunctions = []string{}

var functionsIter = 0

var functionExecutionError error

var mutex = &sync.Mutex{}

// Event channel for triggering warm state updates
var WarmStateUpdateChan = make(chan bool)

func PrepareWarmState(funcs []string, fnMappings map[string]models.Function, nextFuncs []string) {
	// Prepare warm state. First use the passed function list, which will be child
	// functions of last executed parent functions. If space left, use backup BFS
	// list to fill the warm state slots. This also includes finding the best nodes
	// for each warm state function.
	temp := []string{}
	if len(nextFuncs) > 0 {
		localIter := 0
		for len(temp) < numParallel-len(warmFunctions) && localIter < len(nextFuncs) {
			function := nextFuncs[localIter]
			_, ok := funcRunStatus[function]
			if !ok && fnMappings[function].IsWarm == false {
				temp = append(temp, function)
			}
			localIter += 1
		}
	}
	localIter := functionsIter
	for len(temp) < numParallel-len(warmFunctions) && localIter < len(funcs) {
		function := funcs[localIter]
		_, ok := funcRunStatus[function]
		if !ok && fnMappings[function].IsWarm == false {
			temp = append(temp, function)
		}
		localIter += 1
	}
	functionsIter = localIter
	warmFunctions = append(warmFunctions, temp...)
	//logrus.Info(fmt.Sprintf("GoRoutineId: %s Current fns: %v", utils.GetGoroutineID(), warmFunctions))
	// get best nodes for the functions in temp array
}

// Warm state update event handler
func WarmStateUpdateEventHandler(funcs []string, fnMappings map[string]models.Function) {
	for range WarmStateUpdateChan {

		logrus.Info(fmt.Sprintf("GoRoutineId: %s WarmStateUpdateEventHandler received event, updating warm state now...", utils.GetGoroutineID()))

		mutex.Lock()
		PrepareWarmState(funcs, fnMappings, []string{})
		mutex.Unlock()
	}

}

func RunFunction(fn string, fnMappings map[string]models.Function, funcs []string, wg *sync.WaitGroup) {
	// Runs the input function
	mutex.Lock()
	val, ok := funcRunStatus[fn]
	mutex.Unlock()
	flag := false
	if ok {
		if val != "Completed" && val != "Submitted" {
			flag = true
		}
	} else {
		flag = true
	}
	if flag {
		mutex.Lock()
		funcRunStatus[fn] = "Submitted"
		mutex.Unlock()

		// execute the function and get the result. If error, then break out, else update the function status and continue
		// Generate event
		WarmStateUpdateChan <- true
		dataString, _ := base64.StdEncoding.DecodeString(fnMappings[fn].Data)
		logrus.Info(fmt.Sprintf("GoRoutineId: %s function: %s with parameters: %s has been well processed", utils.GetGoroutineID(), fn, string(dataString)))

		mutex.Lock()
		resultsMap[fn] = "processed"
		funcRunStatus[fn] = "Completed"
		mutex.Unlock()
		logrus.Info(fmt.Sprintf("GoRoutineId: %s Result for function: %v\n", utils.GetGoroutineID(), resultsMap))
		// If executed function was in warm state, clear up its slot.
		// Finished: Need to run this in parallel with the execution.

		var wgWarmState sync.WaitGroup

		wgWarmState.Add(1)
		go func() {
			defer wgWarmState.Done()
			mutex.Lock()
			idx := -1
			for i, funcn := range warmFunctions {
				if funcn == fn {
					idx = i
					break
				}
			}
			if idx != -1 {
				warmFunctions = append(warmFunctions[:idx], warmFunctions[idx+1:]...)
			}
			mutex.Unlock()
		}()
		// Find next potential runnable functions from child functions of the currently executed function. Store the potential functions in the potentialRunnables array
		// Use potential runnable functions to fill the empty slots in warm state

		// Find next potential runnable functions from child functions of the currently executed function.
		nextFuncs := fnMappings[fn].Next
		potentialRunnables := []string{}
		for _, nextFunc := range nextFuncs {
			// Check if all dependencies of the next function have been completed.
			dependsOn := fnMappings[nextFunc.Name].DependsOn
			allDependenciesCompleted := true
			for _, dependency := range dependsOn {
				mutex.Lock()
				if funcRunStatus[dependency.Name] != "Completed" {
					allDependenciesCompleted = false
				}
				mutex.Unlock()
				if !allDependenciesCompleted {
					break
				}
			}
			// If all dependencies of the next function have been completed, then the next function is a potential runnable function.
			if allDependenciesCompleted {
				potentialRunnables = append(potentialRunnables, nextFunc.Name)
			}
		}

		// If any new warm state function is runnable, add it to the list of potential runnables
		// For non-warm state functions, find the best nodes
		// Using recursive multi-threading, parallely execute all runnable functions
		var wgRecursive sync.WaitGroup
		for _, fn := range potentialRunnables {
			wgRecursive.Add(1)
			go func(fn string, fnMappings map[string]models.Function, funcs []string, resultsMap map[string]interface{}) {
				defer wgRecursive.Done()
				RunFunction(fn, fnMappings, funcs, &wgRecursive)
				if functionExecutionError != nil {
					return
				}
			}(fn, fnMappings, funcs, resultsMap)
		}
		wgRecursive.Wait()
	}
}

func ScheduleFunctionOnNode(functionName string, data string, fnMappings map[string]models.Function, funcs []string) (*models.Podresult, error) {
	// fnMappings would contain parsed function chain in the following format:
	//	{
	//		"f1": {
	//			"dependsOn": [],
	//			"nextFuncs": [],
	//			"isLast": false,
	//			"breakConditions": [],
	//			"isWarm": false
	//		}
	//	}
	// funcs is an array of all functions in the function chain
	// resultsMap is a placeholder for results

	res := models.Podresult{}
	res.ResultsMap = resultsMap
	// Prepare initial warm state
	PrepareWarmState(funcs, fnMappings, []string{})
	data = base64.StdEncoding.EncodeToString([]byte(data))
	functionInfo := fnMappings[functionName]
	functionInfo.Data = data
	fnMappings[functionName] = functionInfo
	// Using recursive multithreading, execute the function chain
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		functionExecutionError = nil
		RunFunction(functionName, fnMappings, funcs, &wg)
	}()
	wg.Wait()
	return &res, nil
}
