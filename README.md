# NCSU_GEARS

This repository is for NCSU_GEARS summer program(Topic: Serverless).

## Task 1 & 2
### Distribution of Workload
Yang: Task 1

Zhe: Task 2

### Tasks Description:

Task 1: Parse the input function chain and convert it into the desired format:

We are expecting to receive a JSON formatted function chain description from the customer. A sample of such a JSON is attached with this email in the "Gears-Task1" file.

Unfortunately, using the input data as is, would lead to suboptimal performance since we first need to access all the entries in the top level JSON to identify the correct function. Then further nesting complicates the parsing more.

Thus, the task is to convert this JSON into a format/data structure, which allows us to have a faster access to the metadata whenever needed. An example approach could be to have it converted into a hashmap with function name as the key and function data/metadata hashmap, as the value. This is just for example. You can think of other approaches also. It's just that we want to have as fast as possible access to all the function metadata.

As a target, consider parsing the JSON into a format/data structure that allows O(1) access to version, timeout, isLast, next, and dependsOn fields for each function. Further, next and dependsOn should also be parsed in a way that allows O(1) access to metadata of each dependent function in the respective list.

Our aim is to run this on the function chain as soon as the customer registers it, so that all the executions of the function chain would run on this faster accessible data structure.



Task 2: Implement parallel execution of warm state handling:

In our current setup, we maintain some functions in a warm state while the other functions are being executed, to ensure low cold start times. Please read about cold start times in serverless computing, if you are not already aware.

However, in the current state of the code, when a warm function executes, we update the warm state array after the execution of the said function, which could lead to cold start problem not being solved at all.

Therefore, the task is to use multi-threading/goroutines to enable parallel updates to the warm state queue while the current function is being executed. The point to take care of is that the warm state update thread time should be the limiting one. Thus, even if the function execution completes before warm state update, the program shall wait for the warm state update to complete.

For your reference, I have attached the relevant code snippet with this email. Kindly go through this code snippet and implement this parallel warm state update. Note that warm state updates that need to be parallelized are happening in the "RunFunction" function. There are "Todo" comments for identification of the relevant part.

Note that some code has been removed to create this relevant code snippet. So, if any part of the code is unclear, let me know and we can discuss.

### How to Run:
```shell
cd ${APPROOT}
go run main/main.go
```
### Result & Explanation
#### Experiment Configurations:
```yaml
numParallel: 2
functionsIter: 0
```
#### Experiment Results:
```shell
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 18 WarmStateUpdateEventHandler received event, updating warm state now...
GoRoutineId: 18 Current fns: [f1 f2]
GoRoutineId: 29 Current fns: [f1 f2]
GoRoutineId: 5 function: f11 with parameters: data10 has been well processed
GoRoutineId: 29 Result for function: map[f11:processed]

GoRoutineId: 22 Current fns: [f1 f2]
GoRoutineId: 32 function: f4 with parameters: data3 has been well processed
GoRoutineId: 22 Result for function: map[f4:processed]

GoRoutineId: 28 Current fns: [f1 f2]
GoRoutineId: 7 function: f10 with parameters: data9 has been well processed
GoRoutineId: 28 Result for function: map[f10:processed]

GoRoutineId: 26 Current fns: [f1 f2]
GoRoutineId: 34 function: f8 with parameters: data7 has been well processed
GoRoutineId: 26 Result for function: map[f8:processed]

GoRoutineId: 27 Current fns: [f1 f2]
GoRoutineId: 9 function: f9 with parameters: data8 has been well processed
GoRoutineId: 11 function: f13 with parameters:  has been well processed
GoRoutineId: 27 Result for function: map[f13:processed f9:processed]

GoRoutineId: 31 Current fns: [f1 f2]
GoRoutineId: 31 Result for function: map[]

GoRoutineId: 19 Current fns: [f1 f2]
GoRoutineId: 37 function: f1 with parameters: data0 has been well processed
GoRoutineId: 39 function: f2 with parameters:  has been well processed
GoRoutineId: 40 function: f3 with parameters:  has been well processed
GoRoutineId: 14 function: f5 with parameters:  has been well processed
GoRoutineId: 15 function: f6 with parameters:  has been well processed
GoRoutineId: 16 function: f7 with parameters:  has been well processed
GoRoutineId: 66 function: f12 with parameters:  has been well processed
GoRoutineId: 19 Result for function: map[f1:processed f12:processed f2:processed f3:processed f5:processed f6:processed f7:processed]

```
#### Explanation
1. The json strings were mapped to "map[string]models.function" and kept O(1) access to the metadata.
2. [UPDATED IN V2]"dependsOn" and "next" fields has been modified to "map[string]models.function" to remain O(1) access instead of traversing all the json again.
3. For each call to "ScheduleFunctionOnNode()", the main GoRoutine(GoRoutineId: 1) will wait for the user function GoRoutine and warm state update GoRoutine.
4. Use mutex lock to remain thread-safe: Could be switched to read-write lock if reading is significantly larger than writing.
5. [UPDATED IN V2] Add a WarmStateUpdateEventHandler to maintain warm state once a request(event) comes.
6. [UPDATED IN V2] Update test code to simulate concurrency scenarios.

