package orm

import (
	"fmt"
	"time"
)

type Rule struct {
	Core
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	Owner        User
	OwnerID      uint
}

func (r Rule) String() string {
	s := fmt.Sprintf("r.%d: %s", r.ID, r.Name)
	return s
}
