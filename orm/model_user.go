package orm

import (
	"fmt"
	"time"
)

type User struct {
	ModelCore
	Name      string
	IsAdmin   bool
	LastLogin *time.Time
}

func (u User) String() string {
	s := fmt.Sprintf("U%d: \"%s\"\n",
		u.ID, u.Name)
	return s
}
