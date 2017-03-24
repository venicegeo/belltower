package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
)

type Rule struct {
	Common
	PollDuration time.Duration `json:"poll_duration"`
	Expression   string        `json:"expression"`
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
	Common
	PollDuration time.Duration
	Expression   string
}

type RuleFieldsForUpdate struct {
	Name         string
	PollDuration time.Duration
	IsEnabled    bool
	Expression   string
	IsPublic     bool
}

//---------------------------------------------------------------------

func (rule *Rule) GetLoweredName() string { return "rule" }

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
		PollDuration: rule.PollDuration,
		Expression:   rule.Expression,
	}
	read.Common = rule.Common

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
