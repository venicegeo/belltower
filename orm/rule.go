package orm

import (
	"fmt"
	"time"
)

type Rule struct {
	Core
	RuleAttributes

	Owner   User
	OwnerID uint
}

type RuleAttributes struct {
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
}

func (r Rule) String() string {
	s := fmt.Sprintf("r.%d: %s", r.ID, r.Name)
	return s
}
