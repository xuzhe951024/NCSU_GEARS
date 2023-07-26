package models

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
