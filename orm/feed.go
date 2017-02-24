package orm

import (
	"fmt"
	"time"
)

type Feed struct {
	Core
	Name                string
	NumMessagesRecieved uint
	LastMessageAt       *time.Time
	IsEnabled           bool
	Owner               User
	OwnerID             uint
	Persistenceduration time.Duration
	FeedType            FeedType
	FeedTypeID          uint
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
