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

//---------------------------------------------------------------------

type UserFieldsForCreate struct {
	Name      string
	IsAdmin   bool
	IsEnabled bool
}

type UserFieldsForRead struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	IsAdmin     bool
	IsEnabled   bool
	LastLoginAt time.Time
}

type UserFieldsForUpdate struct {
	Name      string
	IsAdmin   bool
	IsEnabled bool
}

func CreateUser(fields *UserFieldsForCreate) (*User, error) {
	user := &User{
		Name:      fields.Name,
		IsAdmin:   fields.IsAdmin,
		IsEnabled: fields.IsEnabled,
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

		IsAdmin:   user.IsAdmin,
		IsEnabled: user.IsEnabled,
	}

	return read, nil
}

func (user *User) Update(update *UserFieldsForUpdate) error {

	if update.Name != "" {
		user.Name = update.Name
	}

	user.IsAdmin = update.IsAdmin
	user.IsEnabled = update.IsEnabled

	return nil
}

func (u User) String() string {
	s := fmt.Sprintf("u.%d: %s", u.ID, u.Name)
	return s
}
