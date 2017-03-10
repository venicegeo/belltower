package orm

import (
	"fmt"
	"time"
)

type Rule struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	OwnerID      uint
	IsPublic     bool
}

//---------------------------------------------------------------------

type RuleFieldsForCreate struct {
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	IsPublic     bool
}

type RuleFieldsForRead struct {
	ID   uint
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	OwnerID      uint
	IsPublic     bool
}

type RuleFieldsForUpdate struct {
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	IsPublic     bool
}

func CreateRule(requestorID uint, fields *RuleFieldsForCreate) (*Rule, error) {
	rule := &Rule{
		Name:         fields.Name,
		IsEnabled:    fields.IsEnabled,
		Expression:   fields.Expression,
		PollDuration: fields.PollDuration,
		OwnerID:      requestorID,
		IsPublic:     fields.IsPublic,
	}

	return rule, nil
}

func (rule *Rule) GetOwnerID() uint {
	return rule.OwnerID
}

func (rule *Rule) GetIsPublic() bool {
	return rule.IsPublic
}

//---------------------------------------------------------------------

func (rule *Rule) Read() (*RuleFieldsForRead, error) {

	read := &RuleFieldsForRead{
		ID:   rule.ID,
		Name: rule.Name,

		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,

		PollDuration: rule.PollDuration,
		IsEnabled:    rule.IsEnabled,
		Expression:   rule.Expression,
		OwnerID:      rule.OwnerID,
		IsPublic:     rule.IsPublic,
	}

	return read, nil
}

func (rule *Rule) Update(update *RuleFieldsForUpdate) error {

	if update.Name != "" {
		rule.Name = update.Name
	}

	rule.IsEnabled = update.IsEnabled
	if update.PollDuration != 0 {
		rule.PollDuration = update.PollDuration
	}
	if update.Expression != "" {
		rule.Expression = update.Expression
	}
	if update.PollDuration != 0 {
		rule.PollDuration = update.PollDuration
	}
	rule.IsPublic = update.IsPublic
	return nil
}

func (f Rule) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
