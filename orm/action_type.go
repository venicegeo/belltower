package orm

import "fmt"

type ActionType struct {
	Core
	ActionTypeAttributes
}

type ActionTypeAttributes struct {
	Name       string
	IsEnabled  bool
	ConfigInfo string
}

func (at ActionType) String() string {
	s := fmt.Sprintf("at.%d:", at.ID)
	return s
}
