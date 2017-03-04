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

func (model *Orm) AddUser(u *User) (uint, error) {
	user := *u
	err := model.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	id := user.ID
	return id, nil
}

func (model *Orm) UpdateUser(id uint, user *User) error {
	u, err := model.GetUser(id)
	if err != nil {
		return err
	}
	*u = *user
	return model.db.Save(u).Error
}

func (model *Orm) DeleteUser(id uint) error {
	u, err := model.GetUser(id)
	if err != nil {
		return err
	}
	if u == nil {
		return fmt.Errorf("record not found u.%d", id)
	}
	err = model.db.Delete(u).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetUser(id uint) (*User, error) {

	u := &User{}
	err := model.db.First(u, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

//---------------------------------------------------------------------

func (model *Orm) AddFeed(feed *Feed) (uint, error) {
	f := *feed
	err := model.db.Create(f).Error
	if err != nil {
		return 0, err
	}
	return f.ID, nil
}

func (model *Orm) UpdateFeed(id uint, feed *Feed) error {
	f, err := model.GetFeed(id)
	if err != nil {
		return err
	}
	*f = *feed
	return model.db.Save(f).Error
}

func (model *Orm) DeleteFeed(id uint) error {
	f, err := model.GetFeed(id)
	if err != nil {
		return err
	}
	if f == nil {
		return fmt.Errorf("record not found f.%d", id)
	}
	err = model.db.Delete(f).Error
	if err != nil {
		return err
	}

	return nil
}

func (model *Orm) GetFeed(id uint) (*Feed, error) {

	f := &Feed{}
	err := model.db.First(f, "id = ?", id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return f, nil
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
