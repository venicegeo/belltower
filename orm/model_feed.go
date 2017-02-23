package orm

import (
	"fmt"
	"time"
)

type Feed struct {
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

type FeedAcl struct {
	ModelCore
	UserID uint
	FeedID uint
}

func (acl FeedAcl) String() string {
	s := fmt.Sprintf("A%d: F%d U%d\n",
		acl.ID, acl.FeedID, acl.UserID)
	return s
}
