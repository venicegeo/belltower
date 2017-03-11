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
	Role        Role
	IsEnabled   bool
	LastLoginAt time.Time
	OwnerID     uint
}

//---------------------------------------------------------------------

type UserFieldsForCreate struct {
	Name      string
	IsEnabled bool
	Role      Role
}

type UserFieldsForRead struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	IsEnabled   bool
	Role        Role
	LastLoginAt time.Time
}

type UserFieldsForUpdate struct {
	Name      string
	IsEnabled bool
	Role      Role
}

func CreateUser(fields *UserFieldsForCreate) (*User, error) {
	user := &User{
		Name:      fields.Name,
		IsEnabled: fields.IsEnabled,
		Role:      fields.Role,
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
		Role:      user.Role,
	}

	return read, nil
}

func (user *User) Update(update *UserFieldsForUpdate) error {

	if update.Name != "" {
		user.Name = update.Name
	}

	user.IsEnabled = update.IsEnabled
	user.Role = update.Role

	return nil
}

func (u User) String() string {
	s := fmt.Sprintf("u.%d: %s", u.ID, u.Name)
	return s
}
