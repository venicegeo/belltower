package orm

import "fmt"

type ActionType struct {
	Core
	ConfigInfo string
}

func (at ActionType) String() string {
	s := fmt.Sprintf("at.%d:", at.ID)
	return s
}
