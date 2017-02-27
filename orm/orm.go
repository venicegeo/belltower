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

func (model *Orm) AddUser(attrs *UserAttributes) (uint, error) {
	user := &User{}
	user.UserAttributes = *attrs
	err := model.db.Create(user).Error
	id := user.ID
	return id, err
}

func (model *Orm) UpdateUser(id uint, attrs *UserAttributes) error {
	r, err := model.GetUser(id)
	if err != nil {
		return err
	}
	r.UserAttributes = *attrs
	return model.db.Save(r).Error
}

func (model *Orm) DeleteUser(id uint) error {
	r, err := model.GetUser(id)
	if err != nil {
		return err
	}
	if r == nil {
		return fmt.Errorf("record not found u.%d", id)
	}
	err = model.db.Delete(r).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetUser(id uint) (*User, error) {

	r := &User{}
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

func (model *Orm) AddFeed(attrs *FeedAttributes) (uint, error) {
	r := &Feed{}
	r.FeedAttributes = *attrs
	err := model.db.Create(r).Error
	id := r.ID
	return id, err
}

func (model *Orm) UpdateFeed(id uint, attrs *FeedAttributes) error {
	r, err := model.GetFeed(id)
	if err != nil {
		return err
	}
	r.FeedAttributes = *attrs
	return model.db.Save(r).Error
}

func (model *Orm) DeleteFeed(id uint) error {
	r, err := model.GetFeed(id)
	if err != nil {
		return err
	}
	if r == nil {
		return fmt.Errorf("record not found f.%d", id)
	}
	err = model.db.Delete(r).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetFeed(id uint) (*Feed, error) {

	r := &Feed{}
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

func (model *Orm) AddRule(attrs *RuleAttributes) (uint, error) {
	r := &Rule{}
	r.RuleAttributes = *attrs
	err := model.db.Create(r).Error
	id := r.ID
	return id, err
}

func (model *Orm) UpdateRule(id uint, attrs *RuleAttributes) error {
	r, err := model.GetRule(id)
	if err != nil {
		return err
	}
	r.RuleAttributes = *attrs
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

func (model *Orm) AddAction(attrs *ActionAttributes) (uint, error) {
	r := &Action{}
	r.ActionAttributes = *attrs
	err := model.db.Create(r).Error
	id := r.ID
	return id, err
}

func (model *Orm) UpdateAction(id uint, attrs *ActionAttributes) error {
	r, err := model.GetAction(id)
	if err != nil {
		return err
	}
	r.ActionAttributes = *attrs
	return model.db.Save(r).Error
}

func (model *Orm) DeleteAction(id uint) error {
	r, err := model.GetAction(id)
	if err != nil {
		return err
	}
	if r == nil {
		return fmt.Errorf("record not found a.%d", id)
	}
	err = model.db.Delete(r).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetAction(id uint) (*Action, error) {

	r := &Action{}
	err := model.db.First(r, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return r, nil
}
