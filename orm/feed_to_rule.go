package orm

import (
	"fmt"
	"time"
)

type FeedToRule struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerID uint
	RuleID  uint
	FeedID  uint
}

func CreateFeedToRule(requestorID uint, feedID uint, ruleID uint) (*FeedToRule, error) {
	feedToRule := &FeedToRule{
		OwnerID: requestorID,
		FeedID:  feedID,
		RuleID:  ruleID,
	}

	return feedToRule, nil
}

func (fra FeedToRule) String() string {
	s := fmt.Sprintf("A%d: F%d\n",
		fra.ID, fra.FeedID)
	return s
}
