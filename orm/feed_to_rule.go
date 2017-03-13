package orm

import (
	"fmt"
	"time"
)

type FeedToRule struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	IsPublic bool
	OwnerID  uint
	RuleID   uint
	FeedID   uint
}

func CreateFeedToRule(requestorID uint, feedID uint, ruleID uint, isPublic bool) (*FeedToRule, error) {
	feedToRule := &FeedToRule{
		OwnerID:  requestorID,
		IsPublic: isPublic,
		FeedID:   feedID,
		RuleID:   ruleID,
	}

	return feedToRule, nil
}

func (feedToRule *FeedToRule) GetOwnerID() uint {
	return feedToRule.OwnerID
}

func (feedToRule *FeedToRule) GetIsPublic() bool {
	return feedToRule.IsPublic
}

func (fra FeedToRule) String() string {
	s := fmt.Sprintf("A%d: F%d\n",
		fra.ID, fra.FeedID)
	return s
}
