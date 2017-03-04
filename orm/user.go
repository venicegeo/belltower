package orm

import (
	"fmt"
	"time"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Name                string
	IsAdmin             bool
	IsEnabled           bool
	LastLoginAt         time.Time
	LastLoginAtInternal string
}

func (u *User) BeforeSave() {
	u.LastLoginAtInternal = u.LastLoginAt.Format(time.RFC3339)
}

func (u User) String() string {
	s := fmt.Sprintf("u.%d: %s", u.ID, u.Name)
	return s
}
