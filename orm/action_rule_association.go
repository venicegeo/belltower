package orm

import (
	"fmt"
	"time"
)

type ActionRuleAssociation struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Rule     Rule
	RuleID   uint
	Action   Action
	ActionID uint
}

func (raa ActionRuleAssociation) String() string {
	s := fmt.Sprintf("raa.%d", raa.ID)
	return s
}
