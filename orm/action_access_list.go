package orm

import (
	"fmt"
	"time"
)

type ActionAccessList struct {
	ID   uint `gorm:"primary_key"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Action   Action
	ActionID uint
	User     User
	UserID   uint
	CanRead  bool
	CanWrite bool
}

func (aal ActionAccessList) String() string {
	s := fmt.Sprintf("aal.%d: \n",
		aal.ID)
	return s
}
