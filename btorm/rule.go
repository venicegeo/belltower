package btorm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

type Rule struct {
	Id           common.Ident  `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Name         string        `json:"name"`
	PollDuration time.Duration `json:"poll_duration"`
	IsEnabled    bool          `json:"is_enabled"`
	Expression   string        `json:"expression"`
	OwnerId      common.Ident  `json:"owner_id"`
	IsPublic     bool          `json:"is_public"`
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
	Id           common.Ident
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	OwnerId      common.Ident
	IsPublic     bool
}

type RuleFieldsForUpdate struct {
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	IsPublic     bool
}

//---------------------------------------------------------------------

func (rule *Rule) GetIndexName() string {
	return "rule_index"
}

func (rule *Rule) GetTypeName() string {
	return "rule_type"
}

func (rule *Rule) GetMapping() string {
	mapping := `{
	"settings":{
	},
	"mappings":{
		"rule_type":{
			"dynamic":"strict",
			"properties":{
				"id":{
					"type":"string"
				},
				"name":{
					"type":"string"
				},
				"created_at":{
					"type":"date"
				},
				"updated_at":{
					"type":"date"
				},
				"is_enabled":{
					"type":"boolean"
				},
				"is_public":{
					"type":"boolean"
				},
				"owner_id":{
					"type":"string"
				},
				"poll_duration":{
					"type":"integer"
				},
				"expression":{
					"type":"string"
				}
			}
		}
	}
}`

	return mapping
}

func (rule *Rule) GetId() common.Ident {
	return rule.Id
}

func (rule *Rule) SetId() common.Ident {
	rule.Id = common.NewId()
	return rule.Id
}

//---------------------------------------------------------------------

func (rule *Rule) SetFieldsForCreate(ownerId common.Ident, ifields interface{}) error {

	fields := ifields.(*RuleFieldsForCreate)

	rule.Name = fields.Name
	rule.IsEnabled = fields.IsEnabled
	rule.Expression = fields.Expression
	rule.PollDuration = fields.PollDuration
	rule.OwnerId = ownerId
	rule.IsPublic = fields.IsPublic

	return nil
}

func (rule *Rule) GetFieldsForRead() (interface{}, error) {

	read := &RuleFieldsForRead{
		Id:           rule.Id,
		Name:         rule.Name,
		CreatedAt:    rule.CreatedAt,
		UpdatedAt:    rule.UpdatedAt,
		PollDuration: rule.PollDuration,
		IsEnabled:    rule.IsEnabled,
		Expression:   rule.Expression,
		OwnerId:      rule.OwnerId,
		IsPublic:     rule.IsPublic,
	}

	return read, nil
}

func (rule *Rule) SetFieldsForUpdate(ifields interface{}) error {

	fields := ifields.(*RuleFieldsForUpdate)

	if fields.Name != "" {
		rule.Name = fields.Name
	}

	rule.IsEnabled = fields.IsEnabled
	rule.IsPublic = fields.IsPublic

	if fields.PollDuration != 0 {
		rule.PollDuration = fields.PollDuration
	}
	if fields.Expression != "" {
		rule.Expression = fields.Expression
	}
	if fields.PollDuration != 0 {
		rule.PollDuration = fields.PollDuration
	}

	return nil
}

//---------------------------------------------------------------------

func (rule *Rule) GetOwnerId() common.Ident {
	return rule.OwnerId
}

func (rule *Rule) GetIsPublic() bool {
	return rule.IsPublic
}

func (rule Rule) String() string {
	s := fmt.Sprintf("a.%s: %s", rule.Id, rule.Name)
	return s
}
