package orm

import "fmt"

type RuleAccessList struct {
	Core
	RuleID   uint
	User     User
	UserID   uint
	CanRead  bool
	CanWrite bool
}

func (ral RuleAccessList) String() string {
	s := fmt.Sprintf("ral.%d", ral.ID)
	return s
}
