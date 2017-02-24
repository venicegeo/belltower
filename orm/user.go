package orm

import (
	"fmt"
	"time"
)

type User struct {
	Core
	UserAttributes

	LastLoginAt time.Time
}

type UserAttributes struct {
	Name      string
	IsAdmin   bool
	IsEnabled bool
}

func (u User) String() string {
	s := fmt.Sprintf("u.%d: %s", u.ID, u.Name)
	return s
}
