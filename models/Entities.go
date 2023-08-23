package models

import "sync"

type Function struct {
	Name            string
	Version         string
	DependsOn       map[string]FunctionIndex
	Next            map[string]struct{ Name string }
	Timeout         string
	IsLast          bool
	BreakConditions []Condition
	Data            string
	IsWarm          bool
}

type UnparsedFunction struct {
	Name      string
	Version   string
	DependsOn []FunctionIndex
	Next      []struct {
		Name string `json:"name"`
	}
	Timeout         string
	IsLast          bool
	BreakConditions []Condition
	Data            string
	IsWarm          bool
}

type FunctionIndex struct {
	Name      string      `json:"name"`
	Required  bool        `json:"required"`
	Condition []Condition `json:"condition"`
}

type Condition struct {
	Key      string
	Operator string
	Val      string
}

type Podresult struct {
	ResultsMap map[string]interface{}
}

type RegisterFunctionChainVO struct {
	Identifier string
	Functions  []UnparsedFunction
}

type ThreadSafeMap struct {
	mu    sync.RWMutex
	items map[string]map[string]Function
}

// NewThreadSafeMap creates a new instance of a ThreadSafeMap.
func NewThreadSafeMap() *ThreadSafeMap {
	return &ThreadSafeMap{
		items: make(map[string]map[string]Function),
	}
}

// Set sets the key to value.
func (tsm *ThreadSafeMap) SetThreadSafeMap(key string, value map[string]Function) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	tsm.items[key] = value
}

// Get gets the value for the given key.
func (tsm *ThreadSafeMap) GetThreadSafeMap(key string) (map[string]Function, bool) {
	tsm.mu.RLock()
	defer tsm.mu.RUnlock()
	val, ok := tsm.items[key]
	return val, ok
}

type ThreadSafeCount struct {
	mu    sync.RWMutex
	count int64
}

// NewThreadSafeCount creates a new instance of a ThreadSafeCount.
func NewThreadSafeCount() *ThreadSafeCount {
	return &ThreadSafeCount{
		count: 0,
	}
}

// Set sets count.
func (tsc *ThreadSafeCount) SetThreadSafeCount(num int64) {
	tsc.mu.Lock()
	defer tsc.mu.Unlock()
	tsc.count += num
}

// Get gets count.
func (tsc *ThreadSafeCount) GetThreadSafeCount() int64 {
	tsc.mu.RLock()
	defer tsc.mu.RUnlock()
	return tsc.count
}
