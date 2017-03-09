package orm

import (
	"os"

	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Orm struct {
	db *gorm.DB
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
		&Feed{}, &FeedRuleAssociation{}, &FeedAccessList{},
		&Rule{}, &RuleAccessList{},
		&Action{}, &ActionRuleAssociation{}, &ActionAccessList{},
		&User{},
	}

	for _, table := range tables {

		err = model.db.CreateTable(table).Error
		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

func (model *Orm) Close() error {
	return model.db.Close()
}

//---------------------------------------------------------------------
// User

func (model *Orm) readUserById(id uint) (*User, error) {
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

func (model *Orm) CreateUser(fields *UserFieldsForCreate) (uint, error) {
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

func (model *Orm) UpdateUser(id uint, fields *UserFieldsForUpdate) error {
	user, err := model.readUserById(id)
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

func (model *Orm) DeleteUser(id uint) error {
	user, err := model.readUserById(id)
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

func (model *Orm) ReadUser(id uint) (*UserFieldsForRead, error) {

	user, err := model.readUserById(id)
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

func (model *Orm) CreateFeed(fields *FeedFieldsForCreate) (uint, error) {
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

func (model *Orm) UpdateFeed(id uint, fields *FeedFieldsForUpdate) error {
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

func (model *Orm) DeleteFeed(id uint) error {
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

func (model *Orm) ReadFeed(id uint) (*FeedFieldsForRead, error) {

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

func (model *Orm) CreateAction(fields *ActionFieldsForCreate) (uint, error) {
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

func (model *Orm) UpdateAction(id uint, fields *ActionFieldsForUpdate) error {
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

func (model *Orm) DeleteAction(id uint) error {
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

func (model *Orm) ReadAction(id uint) (*ActionFieldsForRead, error) {

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

func (model *Orm) CreateRule(fields *RuleFieldsForCreate) (uint, error) {
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

func (model *Orm) UpdateRule(id uint, fields *RuleFieldsForUpdate) error {
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

func (model *Orm) DeleteRule(id uint) error {
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

func (model *Orm) ReadRule(id uint) (*RuleFieldsForRead, error) {

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
