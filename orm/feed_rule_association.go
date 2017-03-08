package orm

import (
	"fmt"
	"time"
)

type FeedRuleAssociation struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Rule   Rule
	RuleID uint
	Feed   Feed
	FeedID uint
}

func (fra FeedRuleAssociation) String() string {
	s := fmt.Sprintf("A%d: F%d\n",
		fra.ID, fra.FeedID)
	return s
}
