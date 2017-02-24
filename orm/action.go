package orm

import "fmt"

type Action struct {
	Core
	Name         string
	ActionType   ActionType
	ActionTypeID uint
}

func (a Action) String() string {
	s := fmt.Sprintf("a.%d: %s", a.ID, a.Name)
	return s
}
