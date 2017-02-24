package orm

import (
	"os"

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
		&FeedType{}, &Feed{}, &FeedRuleAssociation{}, &FeedAccessList{},
		&Rule{}, &RuleAccessList{},
		&ActionType{}, &Action{}, &ActionRuleAssociation{}, &ActionAccessList{},
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

func (model *Orm) AddUser(name string, attrs *UserAttributes) error {
	user := &User{}
	user.UserAttributes = *attrs
	return model.db.Create(user).Error
}

func (model *Orm) UpdateUser(name string, attrs *UserAttributes) error {
	user, err := model.GetUserByName(name)
	if err != nil {
		return err
	}
	user.UserAttributes = *attrs
	return model.db.Save(user).Error
}

func (model *Orm) GetUserByName(name string) (*User, error) {

	user := &User{}
	err := model.db.First(user, "name = ?", name).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
