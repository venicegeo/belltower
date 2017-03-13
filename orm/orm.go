package orm

import (
	"os"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Orm struct {
	db      *gorm.DB
	AdminID uint
}

func NewOrm() (*Orm, error) {
	var err error

	model := &Orm{}

	// ignore errors
	os.Remove("test.db")

	model.db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	// create all the tables
	tables := []interface{}{
		&Feed{}, &FeedToRule{},
		&Rule{}, &RuleToAction{},
		&Action{},
		&User{},
	}

	for _, table := range tables {

		err = model.db.CreateTable(table).Error
		if err != nil {
			return nil, err
		}
	}

	// add the admin user -- requestor is 0, since there really isn't one
	model.AdminID, err = model.CreateUser(0, &UserFieldsForCreate{
		Name:      "admin",
		IsEnabled: true,
		Role:      AdminRole,
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (model *Orm) Close() error {
	return model.db.Close()
}

//---------------------------------------------------------------------
// User

func (model *Orm) readUserByID(id uint) (*User, error) {
	user := &User{}
	err := model.db.First(user, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (model *Orm) CreateUser(requestorID uint, fields *UserFieldsForCreate) (uint, error) {

	if requestorID != model.AdminID {
		return 0, fmt.Errorf("Permission denied.")
	}

	user, err := CreateUser(fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (model *Orm) UpdateUser(requestorID uint, id uint, fields *UserFieldsForUpdate) error {
	if requestorID != model.AdminID {
		return fmt.Errorf("Permission denied.")
	}

	user, err := model.readUserByID(id)
	if err != nil {
		return err
	}
	err = user.Update(fields)
	if err != nil {
		return err
	}
	err = model.db.Save(user).Error
	return err
}

func (model *Orm) DeleteUser(requestorID uint, id uint) error {
	if requestorID != model.AdminID {
		return fmt.Errorf("Permission denied.")
	}

	user, err := model.readUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("record not found f.%d", id)
	}
	err = model.db.Delete(user).Error
	if err != nil {
		return err
	}

	return nil

}

func (model *Orm) ReadUser(requestorID uint, id uint) (*UserFieldsForRead, error) {
	if requestorID != model.AdminID {
		return nil, fmt.Errorf("Permission denied.")
	}

	user, err := model.readUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	fields, err := user.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

//---------------------------------------------------------------------
// Feed

func (model *Orm) readFeedById(id uint) (*Feed, error) {
	feed := &Feed{}
	err := model.db.First(feed, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return feed, nil
}

func (model *Orm) CreateFeed(requestorID uint, fields *FeedFieldsForCreate) (uint, error) {

	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isCreator(requestor.Role) {
		return 0, fmt.Errorf("Permission denied.")
	}

	feed, err := CreateFeed(requestorID, fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(feed).Error
	if err != nil {
		return 0, err
	}
	return feed.ID, nil
}

func (model *Orm) UpdateFeed(requestorID uint, feedID uint, fields *FeedFieldsForUpdate) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	feed, err := model.getFeedAndCheckAuthn(requestor, feedID, UpdateOperation)
	if err != nil {
		return err
	}

	err = feed.Update(fields)
	if err != nil {
		return err
	}
	err = model.db.Save(feed).Error
	return err
}

func (model *Orm) DeleteFeed(requestorID uint, feedID uint) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	feed, err := model.getFeedAndCheckAuthn(requestor, feedID, DeleteOperation)
	if err != nil {
		return err
	}

	err = model.db.Delete(feed).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) ReadFeed(requestorID uint, feedID uint) (*FeedFieldsForRead, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	feed, err := model.getFeedAndCheckAuthn(requestor, feedID, ReadOperation)
	if err != nil {
		return nil, err
	}

	fields, err := feed.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// returns "field,nil" or "nil,err" -- exactly one of the two will be set
func (model *Orm) getFeedAndCheckAuthn(requestor *User, feedID uint, operation Operation) (*Feed, error) {
	feed, err := model.readFeedById(feedID)
	if err != nil {
		return nil, err
	}

	if feed == nil {
		if isCreator(requestor.Role) {
			return nil, fmt.Errorf("Feed %d not found", feedID)
		} else {
			return nil, fmt.Errorf("Permission denied")
		}
	}

	if !isAuthorized(requestor, feed, operation) {
		return nil, fmt.Errorf("Permission denied.")
	}

	return feed, nil
}

//---------------------------------------------------------------------
// Action

func (model *Orm) readActionById(id uint) (*Action, error) {
	rule := &Action{}
	err := model.db.First(rule, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return rule, nil
}

func (model *Orm) CreateAction(requestorID uint, fields *ActionFieldsForCreate) (uint, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isCreator(requestor.Role) {
		return 0, fmt.Errorf("Permission denied.")
	}

	rule, err := CreateAction(requestorID, fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(rule).Error
	if err != nil {
		return 0, err
	}
	return rule.ID, nil
}

func (model *Orm) UpdateAction(requestorID uint, actionID uint, fields *ActionFieldsForUpdate) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	action, err := model.getActionAndCheckAuthn(requestor, actionID, UpdateOperation)
	if err != nil {
		return err
	}

	err = action.Update(fields)
	if err != nil {
		return err
	}
	err = model.db.Save(action).Error
	return err
}

func (model *Orm) DeleteAction(requestorID uint, actionID uint) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	action, err := model.getActionAndCheckAuthn(requestor, actionID, DeleteOperation)
	if err != nil {
		return err
	}

	err = model.db.Delete(action).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) ReadAction(requestorID uint, actionID uint) (*ActionFieldsForRead, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	action, err := model.getActionAndCheckAuthn(requestor, actionID, ReadOperation)
	if err != nil {
		return nil, err
	}

	fields, err := action.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// returns "field,nil" or "nil,err" -- exactly one of the two will be set
func (model *Orm) getActionAndCheckAuthn(requestor *User, actionID uint, operation Operation) (*Action, error) {
	action, err := model.readActionById(actionID)
	if err != nil {
		return nil, err
	}

	if action == nil {
		if isCreator(requestor.Role) {
			return nil, fmt.Errorf("Action %d not found", actionID)
		} else {
			return nil, fmt.Errorf("Permission denied")
		}
	}

	if !isAuthorized(requestor, action, operation) {
		return nil, fmt.Errorf("Permission denied.")
	}

	return action, nil
}

//---------------------------------------------------------------------
// Rule

func (model *Orm) readRuleById(id uint) (*Rule, error) {
	rule := &Rule{}
	err := model.db.First(rule, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return rule, nil
}

func (model *Orm) CreateRule(requestorID uint, fields *RuleFieldsForCreate) (uint, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isCreator(requestor.Role) {
		return 0, fmt.Errorf("Permission denied.")
	}

	rule, err := CreateRule(requestorID, fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(rule).Error
	if err != nil {
		return 0, err
	}

	return rule.ID, nil
}

func (model *Orm) UpdateRule(requestorID uint, ruleID uint, fields *RuleFieldsForUpdate) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	rule, err := model.getRuleAndCheckAuthn(requestor, ruleID, UpdateOperation)
	if err != nil {
		return err
	}

	err = rule.Update(fields)
	if err != nil {
		return err
	}
	err = model.db.Save(rule).Error
	return err
}

func (model *Orm) DeleteRule(requestorID uint, ruleID uint) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	rule, err := model.getRuleAndCheckAuthn(requestor, ruleID, DeleteOperation)
	if err != nil {
		return err
	}

	err = model.db.Delete(rule).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) ReadRule(requestorID uint, ruleID uint) (*RuleFieldsForRead, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	rule, err := model.getRuleAndCheckAuthn(requestor, ruleID, ReadOperation)
	if err != nil {
		return nil, err
	}

	fields, err := rule.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// returns "field,nil" or "nil,err" -- exactly one of the two will be set
func (model *Orm) getRuleAndCheckAuthn(requestor *User, ruleID uint, operation Operation) (*Rule, error) {
	rule, err := model.readRuleById(ruleID)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		if isCreator(requestor.Role) {
			return nil, fmt.Errorf("Rule %d not found", ruleID)
		} else {
			return nil, fmt.Errorf("Permission denied")
		}
	}

	if !isAuthorized(requestor, rule, operation) {
		return nil, fmt.Errorf("Permission denied.")
	}

	return rule, nil
}

//---------------------------------------------------------------------
// Action -> Rule Association

func (model *Orm) CreateFeedToRule(requestorID uint, feedID uint, ruleID uint, isPublic bool) (uint, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isCreator(requestor.Role) {
		return 0, fmt.Errorf("Permission denied.")
	}

	feed, err := model.readFeedById(feedID)
	if err != nil {
		return 0, err
	}

	if feed == nil {
		return 0, fmt.Errorf("Permission denied.")
	}

	rule, err := model.readRuleById(ruleID)
	if err != nil {
		return 0, err
	}
	if rule == nil {
		return 0, fmt.Errorf("Permission denied.")
	}

	if !isAuthorized(requestor, feed, ReadOperation) {
		return 0, fmt.Errorf("Permission denied.")
	}

	if !isAuthorized(requestor, rule, ReadOperation) {
		return 0, fmt.Errorf("Permission denied.")
	}

	assoc, err := CreateFeedToRule(requestorID, feedID, ruleID, isPublic)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(assoc).Error
	if err != nil {
		return 0, err
	}

	return assoc.ID, nil
}

func (model *Orm) ReadFeedToRule(requestorID uint, feedToRuleID uint) (*FeedToRule, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	feedToRule, err := model.getFeedToRuleAndCheckAuthn(requestor, feedToRuleID, ReadOperation)
	if err != nil {
		return nil, err
	}

	return feedToRule, nil
}

func (model *Orm) DeleteFeedToRule(requestorID uint, feedToRuleID uint) (uint, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	feedToRule, err := model.getFeedToRuleAndCheckAuthn(requestor, feedToRuleID, DeleteOperation)
	if err != nil {
		return 0, err
	}

	err = model.db.Delete(feedToRule).Error
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (model *Orm) getFeedToRuleAndCheckAuthn(requestor *User, feedToRuleID uint, operation Operation) (*FeedToRule, error) {
	feedToRule, err := model.readFeedToRuleById(feedToRuleID)
	if err != nil {
		return nil, err
	}

	if feedToRule == nil {
		if isCreator(requestor.Role) {
			return nil, fmt.Errorf("FeedToRule %d not found", feedToRuleID)
		} else {
			return nil, fmt.Errorf("Permission denied")
		}
	}

	if !isAuthorized(requestor, feedToRule, operation) {
		return nil, fmt.Errorf("Permission denied.")
	}

	return feedToRule, nil
}

func (model *Orm) readFeedToRuleById(id uint) (*FeedToRule, error) {
	feedToRule := &FeedToRule{}
	err := model.db.First(feedToRule, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return feedToRule, nil
}

//---------------------------------------------------------------------

func (model *Orm) CreateRuleToAction(requestorID uint, ruleID uint, actionID uint, isPublic bool) (uint, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isCreator(requestor.Role) {
		return 0, fmt.Errorf("Permission denied.")
	}

	rule, err := model.readRuleById(ruleID)
	if err != nil {
		return 0, err
	}

	if rule == nil {
		return 0, fmt.Errorf("Permission denied.")
	}

	action, err := model.readActionById(actionID)
	if err != nil {
		return 0, err
	}
	if action == nil {
		return 0, fmt.Errorf("Permission denied.")
	}

	if !isAuthorized(requestor, rule, ReadOperation) {
		return 0, fmt.Errorf("Permission denied.")
	}

	if !isAuthorized(requestor, action, ReadOperation) {
		return 0, fmt.Errorf("Permission denied.")
	}

	assoc, err := CreateRuleToAction(requestorID, ruleID, actionID, isPublic)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(assoc).Error
	if err != nil {
		return 0, err
	}

	return assoc.ID, nil
}

func (model *Orm) ReadRuleToAction(requestorID uint, ruleToActionID uint) (*RuleToAction, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	ruleToAction, err := model.getRuleToActionAndCheckAuthn(requestor, ruleToActionID, ReadOperation)
	if err != nil {
		return nil, err
	}

	return ruleToAction, nil
}

func (model *Orm) DeleteRuleToAction(requestorID uint, ruleToActionID uint) (uint, error) {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	ruleToAction, err := model.getRuleToActionAndCheckAuthn(requestor, ruleToActionID, DeleteOperation)
	if err != nil {
		return 0, err
	}

	err = model.db.Delete(ruleToAction).Error
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (model *Orm) getRuleToActionAndCheckAuthn(requestor *User, ruleToActionID uint, operation Operation) (*RuleToAction, error) {
	ruleToAction, err := model.readRuleToActionById(ruleToActionID)
	if err != nil {
		return nil, err
	}

	if ruleToAction == nil {
		if isCreator(requestor.Role) {
			return nil, fmt.Errorf("RuleToAction %d not found", ruleToActionID)
		} else {
			return nil, fmt.Errorf("Permission denied")
		}
	}

	if !isAuthorized(requestor, ruleToAction, operation) {
		return nil, fmt.Errorf("Permission denied.")
	}

	return ruleToAction, nil
}

func (model *Orm) readRuleToActionById(id uint) (*RuleToAction, error) {
	ruleToAction := &RuleToAction{}
	err := model.db.First(ruleToAction, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return ruleToAction, nil
}
