package models

type Function struct {
	Name            string
	Version         string
	DependsOn       []string
	Next            []string
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
