package orm

import (
	"fmt"
	"time"
)

type Feed struct {
	Core
	FeedAttributes

	NumMessagesRecieved uint
	LastMessageAt       *time.Time
	Owner               User
	OwnerID             uint
	FeedType            FeedType
	FeedTypeID          uint
}

type FeedAttributes struct {
	Name                string
	IsEnabled           bool
	PersistenceDuration time.Duration
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
