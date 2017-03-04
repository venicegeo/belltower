package orm

import (
	"fmt"
	"time"
)

type Feed struct {
	Core
	FeedAttributes

	FeedType            string
	NumMessagesRecieved uint
	LastMessageAt       *time.Time
	Owner               User
	OwnerID             uint
}

// FeedAttributes are the things that can be set (or modified)
// by the user.
type FeedAttributes struct {
	Name       string
	IsEnabled  bool
	ConfigInfo string // encodes info specifc to a feed type
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
