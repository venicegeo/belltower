package orm

import (
	"fmt"
	"time"
)

type FeedAccessList struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

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
