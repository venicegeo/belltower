package orm

import (
	"fmt"
	"time"
)

type Condition struct {
	ModelCore
	Name        string
	Owner       User
	OwnerID     uint
	DslType     string
	NumRecieved uint
	FeedType    FeedType
	FeedTypeID  uint
	LastMessage *time.Time
	IsEnabled   bool
	Persistence time.Duration
}

type ConditionAcl struct {
	ModelCore
	UserID      uint
	ConditionID uint
}

type ConditionFeed struct {
	ModelCore
	ConditionID uint
	FeedID      uint
}

func (acl ConditionAcl) String() string {
	s := fmt.Sprintf("A%d: C%d U%d\n",
		acl.ID, acl.ConditionID, acl.UserID)
	return s
}

func (acl ConditionFeed) String() string {
	s := fmt.Sprintf("CF%d: C%d F%d\n",
		acl.ID, acl.ConditionID, acl.FeedID)
	return s
}
