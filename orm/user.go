package orm

import (
	"fmt"
	"time"
)

// Users are never publically visible: only the admin has access rights.
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Rights      uint // what rights this user has
	IsAdmin     bool
	IsEnabled   bool
	LastLoginAt time.Time
	OwnerID     uint
}

//---------------------------------------------------------------------

type UserFieldsForCreate struct {
	Name      string
	IsEnabled bool
	Rights    uint
	IsAdmin   bool
}

type UserFieldsForRead struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	IsEnabled   bool
	IsAdmin     bool
	Rights      uint
	LastLoginAt time.Time
}

type UserFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	Rights    uint
}

func CreateUser(fields *UserFieldsForCreate) (*User, error) {
	user := &User{
		Name:      fields.Name,
		IsEnabled: fields.IsEnabled,
		Rights:    fields.Rights,
		IsAdmin:   fields.IsAdmin,
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

		IsEnabled: user.IsEnabled,
		Rights:    user.Rights,
		IsAdmin:   user.IsAdmin,
	}

	return read, nil
}

func (user *User) Update(update *UserFieldsForUpdate) error {

	if update.Name != "" {
		user.Name = update.Name
	}

	user.IsEnabled = update.IsEnabled
	user.Rights = update.Rights

	return nil
}

func (u User) String() string {
	s := fmt.Sprintf("u.%d: %s", u.ID, u.Name)
	return s
}
