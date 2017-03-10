package orm

import (
	"os"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Orm struct {
	db      *gorm.DB
	adminID uint
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
		&Feed{}, &FeedRuleAssociation{},
		&Rule{},
		&Action{}, &ActionRuleAssociation{},
		&User{},
	}

	for _, table := range tables {

		err = model.db.CreateTable(table).Error
		if err != nil {
			return nil, err
		}
	}

	// add the admin user -- requestor is 0, since there really isn't one
	model.adminID, err = model.CreateUser(0, &UserFieldsForCreate{
		Name:        "admin",
		IsEnabled:   true,
		Permissions: 1,
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
	if requestorID != model.adminID {
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
	if requestorID != model.adminID {
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
	if requestorID != model.adminID {
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
	if requestorID != model.adminID {
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

	user, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isAuthorized(user, CanCreateFeed) {
		return 0, fmt.Errorf("Permission denied.")
	}

	feed, err := CreateFeed(fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(feed).Error
	if err != nil {
		return 0, err
	}
	return feed.ID, nil
}

func (model *Orm) UpdateFeed(requestorID uint, id uint, fields *FeedFieldsForUpdate) error {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	if !isAuthorized(user, CanUpdateFeed) {
		return fmt.Errorf("Permission denied.")
	}

	feed, err := model.readFeedById(id)
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

func (model *Orm) DeleteFeed(requestorID uint, id uint) error {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	if !isAuthorized(user, CanDeleteFeed) {
		return fmt.Errorf("Permission denied.")
	}

	feed, err := model.readFeedById(id)
	if err != nil {
		return err
	}
	if feed == nil {
		return fmt.Errorf("record not found f.%d", id)
	}
	err = model.db.Delete(feed).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) ReadFeed(requestorID uint, id uint) (*FeedFieldsForRead, error) {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	if !isAuthorized(user, CanReadFeed) {
		return nil, fmt.Errorf("Permission denied.")
	}

	feed, err := model.readFeedById(id)
	if err != nil {
		return nil, err
	}

	if feed == nil {
		return nil, nil
	}

	fields, err := feed.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
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
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isAuthorized(user, CanCreateAction) {
		return 0, fmt.Errorf("Permission denied.")
	}

	rule, err := CreateAction(fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(rule).Error
	if err != nil {
		return 0, err
	}
	return rule.ID, nil
}

func (model *Orm) UpdateAction(requestorID uint, id uint, fields *ActionFieldsForUpdate) error {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	if !isAuthorized(user, CanUpdateAction) {
		return fmt.Errorf("Permission denied.")
	}

	rule, err := model.readActionById(id)
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

func (model *Orm) DeleteAction(requestorID uint, id uint) error {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	if !isAuthorized(user, CanDeleteAction) {
		return fmt.Errorf("Permission denied.")
	}

	rule, err := model.readActionById(id)
	if err != nil {
		return err
	}
	if rule == nil {
		return fmt.Errorf("record not found f.%d", id)
	}
	err = model.db.Delete(rule).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) ReadAction(requestorID uint, id uint) (*ActionFieldsForRead, error) {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	if !isAuthorized(user, CanReadAction) {
		return nil, fmt.Errorf("Permission denied.")
	}

	rule, err := model.readActionById(id)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, nil
	}

	fields, err := rule.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
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
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return 0, err
	}

	if !isAuthorized(user, CanCreateRule) {
		return 0, fmt.Errorf("Permission denied.")
	}

	rule, err := CreateRule(fields)
	if err != nil {
		return 0, err
	}

	err = model.db.Create(rule).Error
	if err != nil {
		return 0, err
	}
	return rule.ID, nil
}

func (model *Orm) UpdateRule(requestorID uint, id uint, fields *RuleFieldsForUpdate) error {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	if !isAuthorized(user, CanUpdateRule) {
		return fmt.Errorf("Permission denied.")
	}

	rule, err := model.readRuleById(id)
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

func (model *Orm) DeleteRule(requestorID uint, id uint) error {
	requestor, err := model.readUserByID(requestorID)
	if err != nil {
		return err
	}

	if !isAuthorized(requestor, CanDeleteRule) {
		return fmt.Errorf("Permission denied.")
	}

	rule, err := model.readRuleById(id)
	if err != nil {
		return err
	}
	if rule == nil {
		return fmt.Errorf("record not found f.%d", id)
	}
	err = model.db.Delete(rule).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) ReadRule(requestorID uint, id uint) (*RuleFieldsForRead, error) {
	user, err := model.readUserByID(requestorID)
	if err != nil {
		return nil, err
	}

	if !isAuthorized(user, CanReadRule) {
		return nil, fmt.Errorf("Permission denied.")
	}

	rule, err := model.readRuleById(id)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, nil
	}

	fields, err := rule.Read()
	if err != nil {
		return nil, err
	}

	return fields, nil
}
