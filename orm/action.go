package orm

import "fmt"

type Action struct {
	Core
	ActionAttributes

	ActionType   ActionType
	ActionTypeID uint
}

type ActionAttributes struct {
	Name      string
	IsEnabled bool
}

func (a Action) String() string {
	s := fmt.Sprintf("a.%d: %s", a.ID, a.Name)
	return s
}
