package orm

import (
	"fmt"
	"time"
)

type RuleToAction struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerID  uint
	ActionID uint
	RuleID   uint
}

func CreateRuleToAction(requestorID uint, ruleID uint, actionID uint) (*RuleToAction, error) {
	ruleToAction := &RuleToAction{
		OwnerID:  requestorID,
		RuleID:   ruleID,
		ActionID: actionID,
	}

	return ruleToAction, nil
}

func (raa RuleToAction) String() string {
	s := fmt.Sprintf("raa.%d", raa.ID)
	return s
}
