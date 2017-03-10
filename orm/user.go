package orm

import (
	"fmt"
	"time"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Rights      uint // what rights this user has
	IsPublic    bool // what restrictions are there for others trying to access this object
	IsAdmin     bool
	IsEnabled   bool
	LastLoginAt time.Time
	OwnerID     uint
}

//---------------------------------------------------------------------

type UserFieldsForCreate struct {
	Name        string
	IsEnabled   bool
	Permissions uint
	IsAdmin     bool
}

type UserFieldsForRead struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	IsEnabled   bool
	IsAdmin     bool
	Permissions uint
	LastLoginAt time.Time
}

type UserFieldsForUpdate struct {
	Name        string
	IsEnabled   bool
	Permissions uint
}

func CreateUser(fields *UserFieldsForCreate) (*User, error) {
	user := &User{
		Name:        fields.Name,
		IsEnabled:   fields.IsEnabled,
		Permissions: fields.Permissions,
		IsAdmin:     fields.IsAdmin,
	}

	return user, nil
}

//---------------------------------------------------------------------

func (user *User) Read() (*UserFieldsForRead, error) {

	read := &UserFieldsForRead{
		ID:   user.ID,
		Name: user.Name,

		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,

		IsEnabled:   user.IsEnabled,
		Permissions: user.Permissions,
		IsAdmin:     user.IsAdmin,
	}

	return read, nil
}

func (user *User) Update(update *UserFieldsForUpdate) error {

	if update.Name != "" {
		user.Name = update.Name
	}

	user.IsEnabled = update.IsEnabled
	user.Permissions = update.Permissions

	return nil
}

func (u User) String() string {
	s := fmt.Sprintf("u.%d: %s", u.ID, u.Name)
	return s
}
