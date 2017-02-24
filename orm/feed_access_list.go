package orm

import "fmt"

type FeedAccessList struct {
	Core
	Feed     Feed
	FeedID   uint
	User     User
	UserID   uint
	CanRead  bool
	CanWrite bool
}

func (fal FeedAccessList) String() string {
	s := fmt.Sprintf("fal.%d", fal.ID)
	return s
}
