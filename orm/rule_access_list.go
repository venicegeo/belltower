package orm

import (
	"fmt"
	"time"
)

type RuleAccessList struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	RuleID   uint
	User     User
	UserID   uint
	CanRead  bool
	CanWrite bool
}

func (ral RuleAccessList) String() string {
	s := fmt.Sprintf("ral.%d", ral.ID)
	return s
}
