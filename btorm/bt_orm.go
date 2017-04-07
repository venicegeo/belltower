package btorm

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type BtOrm struct {
	Orm         *esorm.Orm
	prefix      string
	objectTypes []esorm.Elasticable
}

type OrmOption int

const (
	OrmOptionCreate       OrmOption = iota // removes old db if present, then creates
	OrmOptionOpen                          // fails if not already exists
	OrmOptionOpenOrCreate                  // if exists then open, else create
)

//---------------------------------------------------------------------

func NewBtOrm(prefix string, ormOption OrmOption) (*BtOrm, error) {
	orm, err := esorm.NewOrm()
	if err != nil {
		return nil, err
	}

	types := []esorm.Elasticable{
		&Action{},
		&Feed{},
		&Rule{},
		&User{},
	}

	btOrm := &BtOrm{
		Orm:         orm,
		prefix:      prefix,
		objectTypes: types,
	}

	switch ormOption {
	case OrmOptionCreate:
		err = btOrm.createDatabase()
		if err != nil {
			return nil, err
		}
	case OrmOptionOpen:
		err = btOrm.openDatabase()
		if err != nil {
			return nil, err
		}
	case OrmOptionOpenOrCreate:
		exists, err := btOrm.databaseExists()
		if exists {
			err = btOrm.createDatabase()
			if err != nil {
				return nil, err
			}
		} else {
			err = btOrm.createDatabase()
			if err != nil {
				return nil, err
			}
		}
	}

	return btOrm, nil
}

func (btOrm *BtOrm) databaseExists() (bool, error) {

	for _, t := range btOrm.objectTypes {
		exists, err := btOrm.Orm.IndexExists(t)
		if err != nil {
			return false, err
		}

		if !exists {
			return false, nil
		}
	}
	return true, nil
}

func (btOrm *BtOrm) createDatabase() error {

	err := btOrm.deleteDatabase()
	if err != nil {
		return err
	}

	for _, t := range btOrm.objectTypes {
		err = btOrm.Orm.CreateIndex(t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (btOrm *BtOrm) openDatabase() error {

	exists, err := btOrm.databaseExists()
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("database does not exist")
	}

	// nothing to actually do here

	return nil
}

func (btOrm *BtOrm) deleteDatabase() error {

	// try to delete all the indexes

	for _, t := range btOrm.objectTypes {
		exists, err := btOrm.Orm.IndexExists(t)
		if err != nil {
			return err
		}

		if exists {
			err = btOrm.Orm.DeleteIndex(t)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (orm *BtOrm) Close() error {
	return orm.Orm.Close()
}

func (orm *BtOrm) GetAdminID() common.Ident {
	panic("getadminid")
}

//---------------------------------------------------------------------

// THE BELOW FUNCTIONS ARE LARGELY TO PROMOTE TYPE-SAFETY

//---------------------------------------------------------------------

func (orm *BtOrm) CreateAction(requestorID common.Ident, fields *Action) (common.Ident, error) {
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadAction(id common.Ident) (*Action, error) {
	action := &Action{}
	action.Id = id
	fields, err := orm.Orm.ReadDocument(action)
	if err != nil {
		return nil, err
	}
	return fields.(*Action), nil
}

func (orm *BtOrm) UpdateAction(id common.Ident, fields *Action) error {
	fields.Id = id
	return orm.Orm.UpdateDocument(fields)
}

func (orm *BtOrm) DeleteAction(id common.Ident) error {
	action := &Action{}
	action.Id = id
	return orm.Orm.DeleteDocument(action)
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateFeed(requestorID common.Ident, fields *Feed) (common.Ident, error) {
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadFeed(id common.Ident) (*Feed, error) {
	feed := &Feed{}
	feed.Id = id
	fields, err := orm.Orm.ReadDocument(feed)
	if err != nil {
		return nil, err
	}
	return fields.(*Feed), nil
}

func (orm *BtOrm) ReadFeeds(from int, size int) ([]*Feed, int64, error) {
	ary := make([]esorm.Elasticable, size)
	for i := range ary {
		ary[i] = &Feed{}
	}
	ary2, count, err := orm.Orm.ReadDocuments(ary, from, size)
	if err != nil {
		return nil, 0, err
	}

	ret := make([]*Feed, count)
	for i, v := range ary2 {
		ret[i] = v.(*Feed)
	}
	return ret, count, nil
}

func (orm *BtOrm) UpdateFeed(id common.Ident, fields *Feed) error {
	fields.Id = id
	return orm.Orm.UpdateDocument(fields)
}

func (orm *BtOrm) DeleteFeed(id common.Ident) error {
	feed := &Feed{}
	feed.Id = id
	return orm.Orm.DeleteDocument(feed)
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateRule(requestorID common.Ident, fields *Rule) (common.Ident, error) {
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadRule(id common.Ident) (*Rule, error) {
	rule := &Rule{}
	rule.Id = id
	fields, err := orm.Orm.ReadDocument(rule)
	if err != nil {
		return nil, err
	}
	return fields.(*Rule), nil
}

func (orm *BtOrm) UpdateRule(id common.Ident, fields *Rule) error {
	fields.Id = id
	return orm.Orm.UpdateDocument(fields)
}

func (orm *BtOrm) DeleteRule(id common.Ident) error {
	rule := &Rule{}
	rule.Id = id
	return orm.Orm.DeleteDocument(rule)
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateUser(requestorID common.Ident, fields *User) (common.Ident, error) {
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadUser(id common.Ident) (*User, error) {
	user := &User{}
	user.Id = id
	fields, err := orm.Orm.ReadDocument(user)
	if err != nil {
		return nil, err
	}
	return fields.(*User), nil
}

func (orm *BtOrm) UpdateUser(id common.Ident, fields *User) error {
	fields.Id = id
	return orm.Orm.UpdateDocument(fields)
}

func (orm *BtOrm) DeleteUser(id common.Ident) error {
	user := &User{}
	user.Id = id
	return orm.Orm.DeleteDocument(user)
}
