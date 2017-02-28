package orm

import "fmt"

type Action struct {
	Core
	ActionAttributes
}

type ActionAttributes struct {
	Name       string
	IsEnabled  bool
	ConfigInfo string
}

func (a Action) String() string {
	s := fmt.Sprintf("a.%d: %s", a.ID, a.Name)
	return s
}
