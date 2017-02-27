package orm

import "fmt"

type ActionAccessList struct {
	Core
	Action   Action
	ActionID uint
	User     User
	UserID   uint
	CanRead  bool
	CanWrite bool
}

func (aal ActionAccessList) String() string {
	s := fmt.Sprintf("aal.%d: \n",
		aal.ID)
	return s
}
