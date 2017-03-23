package orm2

import (
	"fmt"
	"time"
)

type Rule struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	OwnerId      string
	IsPublic     bool
}

//---------------------------------------------------------------------

type RuleFieldsForCreate struct {
	Id           string
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	IsPublic     bool
}

type RuleFieldsForRead struct {
	Id   string
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	OwnerId      string
	IsPublic     bool
}

type RuleFieldsForUpdate struct {
	Id           string
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	IsPublic     bool
}

func CreateRule(requestorId string, fields *RuleFieldsForCreate) (*Rule, error) {
	rule := &Rule{
		Id:           fields.Id,
		Name:         fields.Name,
		IsEnabled:    fields.IsEnabled,
		Expression:   fields.Expression,
		PollDuration: fields.PollDuration,
		OwnerId:      requestorId,
		IsPublic:     fields.IsPublic,
	}

	return rule, nil
}

func (rule *Rule) GetOwnerId() string {
	return rule.OwnerId
}

func (rule *Rule) GetIsPublic() bool {
	return rule.IsPublic
}

//---------------------------------------------------------------------

func (rule *Rule) Read() (*RuleFieldsForRead, error) {

	read := &RuleFieldsForRead{
		Id:   rule.Id,
		Name: rule.Name,

		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,

		PollDuration: rule.PollDuration,
		IsEnabled:    rule.IsEnabled,
		Expression:   rule.Expression,
		OwnerId:      rule.OwnerId,
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
	s := fmt.Sprintf("f.%s: %s", f.Id, f.Name)
	return s
}
