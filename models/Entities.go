package models

type Function struct {
	Name            string
	Version         string
	DependsOn       map[string]Function
	Next            map[string]Function
	Timeout         string
	IsLast          bool
	BreakConditions []Condition
	Data            string
	IsWarm          bool
}

type UnparsedFunction struct {
	Name            string
	Version         string
	DependsOn       []struct{ Name string }
	Next            []struct{ Name string }
	Timeout         string
	IsLast          bool
	BreakConditions []Condition
	Data            string
	IsWarm          bool
}

type Condition struct {
	Key      string
	Operator string
	Val      string
}

type Podresult struct {
	ResultsMap map[string]interface{}
}
