package orm

import (
	"fmt"
	"time"
)

type RuleToAction struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsPublic  bool

	OwnerID  uint
	ActionID uint
	RuleID   uint
}

func CreateRuleToAction(requestorID uint, ruleID uint, actionID uint, isPublic bool) (*RuleToAction, error) {
	ruleToAction := &RuleToAction{
		OwnerID:  requestorID,
		IsPublic: isPublic,
		RuleID:   ruleID,
		ActionID: actionID,
	}

	return ruleToAction, nil
}

func (ruleToAction *RuleToAction) GetOwnerID() uint {
	return ruleToAction.OwnerID
}

func (ruleToAction *RuleToAction) GetIsPublic() bool {
	return ruleToAction.IsPublic
}

func (raa RuleToAction) String() string {
	s := fmt.Sprintf("raa.%d", raa.ID)
	return s
}
