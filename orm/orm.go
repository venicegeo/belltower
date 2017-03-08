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

func (model *Orm) AddRule(rule *Rule) (uint, error) {
	r := *rule
	err := model.db.Create(r).Error
	if err != nil {
		return 0, err
	}
	id := r.ID
	return id, err
}

func (model *Orm) UpdateRule(id uint, rule *Rule) error {
	r, err := model.GetRule(id)
	if err != nil {
		return err
	}
	*r = *rule
	return model.db.Save(r).Error
}

func (model *Orm) DeleteRule(id uint) error {
	r, err := model.GetRule(id)
	if err != nil {
		return err
	}
	if r == nil {
		return fmt.Errorf("record not found r.%d", id)
	}
	err = model.db.Delete(r).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetRule(id uint) (*Rule, error) {

	r := &Rule{}
	err := model.db.First(r, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return r, nil
}

//---------------------------------------------------------------------

func (model *Orm) AddAction(action *Action) (uint, error) {
	r := *action
	err := model.db.Create(r).Error
	if err != nil {
		return 0, err
	}
	id := r.ID
	return id, nil
}

func (model *Orm) UpdateAction(id uint, action *Action) error {
	a, err := model.GetAction(id)
	if err != nil {
		return err
	}
	*a = *action
	return model.db.Save(a).Error
}

func (model *Orm) DeleteAction(id uint) error {
	a, err := model.GetAction(id)
	if err != nil {
		return err
	}
	if a == nil {
		return fmt.Errorf("record not found a.%d", id)
	}
	err = model.db.Delete(a).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetAction(id uint) (*Action, error) {

	a := &Action{}
	err := model.db.First(a, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return a, nil
}
