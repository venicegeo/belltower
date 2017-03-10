package orm

import (
	"fmt"
	"time"
)

type Rule struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name           string
	PollDuration   time.Duration
	IsEnabled      bool
	Expression     string
	OwnerID        uint
	PublicCanRead  bool
	PublicCanWrite bool
}

//---------------------------------------------------------------------

type RuleFieldsForCreate struct {
	Name           string
	PollDuration   time.Duration
	IsEnabled      bool
	Expression     string
	OwnerID        uint
	PublicCanRead  bool
	PublicCanWrite bool
}

type RuleFieldsForRead struct {
	ID   uint
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	PollDuration   time.Duration
	IsEnabled      bool
	Expression     string
	OwnerID        uint
	PublicCanRead  bool
	PublicCanWrite bool
}

type RuleFieldsForUpdate struct {
	Name           string
	PollDuration   time.Duration
	IsEnabled      bool
	Expression     string
	PublicCanRead  bool
	PublicCanWrite bool
}

func CreateRule(fields *RuleFieldsForCreate) (*Rule, error) {
	rule := &Rule{
		Name:           fields.Name,
		IsEnabled:      fields.IsEnabled,
		Expression:     fields.Expression,
		PollDuration:   fields.PollDuration,
		OwnerID:        fields.OwnerID,
		PublicCanRead:  fields.PublicCanRead,
		PublicCanWrite: fields.PublicCanWrite,
	}

	return rule, nil
}

//---------------------------------------------------------------------

func (rule *Rule) Read() (*RuleFieldsForRead, error) {

	read := &RuleFieldsForRead{
		ID:   rule.ID,
		Name: rule.Name,

		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,

		PollDuration:   rule.PollDuration,
		IsEnabled:      rule.IsEnabled,
		Expression:     rule.Expression,
		OwnerID:        rule.OwnerID,
		PublicCanRead:  rule.PublicCanRead,
		PublicCanWrite: rule.PublicCanWrite,
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
	rule.PublicCanRead = update.PublicCanRead
	rule.PublicCanWrite = update.PublicCanWrite
	return nil
}

func (f Rule) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
