package orm

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Orm struct {
	MyDB        *common.MyDB
	prefix      string
	objectTypes []*Table{}
}

func NewOrm() *Orm {
	orm := &Orm{MyDB: &MyDB{}}
	orm.objectTypes = []interface{}{
		&Action{}, &Rule{}, &User{}, &Feed{},
	}
}

//---------------------------------------------------------------------

func (orm *Orm) Create() error {
	var err error

	err = orm.MyDB.Open(true)
	if err != nil {
		return err
	}

	ok, err := orm.isInitialized(orm)
	if err != nil {
		_ = orm.Close()
		return err
	}
	if !ok {
		_ = orm.Close()
		return fmt.Errorf("DB indexes not all valid")
	}

	return nil
}

func (orm *Orm) Open() error {
	var err error

	err = orm.MyDB.Open(false)
	if err != nil {
		return err
	}

	ok, err := orm.isInitialized(orm)
	if err != nil {
		_ = orm.Close()
		return err
	}
	if !ok {
		_ = orm.Close()
		return fmt.Errorf("DB indexes not all valid")
	}

	return nil
}

// remove all the indxes
func deleteTables(orm *BtOrm) error {
	for _, t := range orm.objectTypes {
		exists, err := orm.Orm.IndexExists(t)
		if err != nil {
			return err
		}
		if exists {
			err = orm.Orm.DeleteIndex(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// return true if all the DB tables exist
func (orm *Orm) HasTables() (bool, error) {
	for _, t := range []string{"a", "b"} {
		ok, err := orm.MyDB.HasTable(t)
		if err != nil {
			return false, nil
		}
		if !ok {
			return false, fmt.Errorf("Table %s does not exist", t)
		}
	}
	return true, nil
}

// create all the tables
func createTables(orm *BtOrm) error {
	for _, t := range orm.objectTypes {
		err := orm.Orm.CreateIndex(t, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func DatabaseDelete() error {
	var err error

	orm := &BtOrm{}
	err = orm.Open()
	if err != nil {
		return err
	}
	defer orm.Close() // ignore errors

	// try three times
	for i := 0; i < 3; i++ {
		err = deleteIndexes(orm)
		if err != nil {
			return err
		}

		exists, err := anyIndexExists(orm)
		if err != nil {
			return err
		}
		if !exists {
			// success!
			return nil
		}
	}

	// at least one isn't getting deleted
	return fmt.Errorf("Unable to delete all indexes from database")
}

func DatabaseInit() error {
	var err error

	orm := &BtOrm{}
	err = orm.Open()
	if err != nil {
		return err
	}
	defer orm.Orm.Close() // ignore error

	err = deleteIndexes(orm)
	if err != nil {
		return err
	}

	err = createIndexes(orm)
	if err != nil {
		return err
	}

	return nil
}

func (orm *BtOrm) Close() error {
	return orm.Orm.Close()
}

func (orm *BtOrm) GetAdminID() common.Ident {
	panic("internal error")
}

//---------------------------------------------------------------------

// THE BELOW FUNCTIONS ARE LARGELY TO PROMOTE TYPE-SAFETY

//---------------------------------------------------------------------

func (orm *BtOrm) CreateAction(requestorID common.Ident, fields *Action) (common.Ident, error) {
	fields.Id = ""
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadAction(id common.Ident) (*Action, error) {
	action := &Action{}
	action.Id = id
	action2, err := orm.Orm.ReadDocument(action)
	if err != nil {
		return nil, err
	}
	return action2.(*Action), nil
}

func (orm *BtOrm) ReadActions(from int, size int) ([]*Action, int64, error) {
	ary2, count, err := orm.Orm.ReadDocuments(&Action{}, from, size)
	if err != nil {
		return nil, 0, err
	}

	ret := make([]*Action, count)
	for i, v := range ary2 {
		ret[i] = v.(*Action)
	}
	return ret, count, nil
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
	fields.Id = ""
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadFeed(id common.Ident) (*Feed, error) {
	feed := &Feed{}
	feed.Id = id
	feed2, err := orm.Orm.ReadDocument(feed)
	if err != nil {
		return nil, err
	}
	return feed2.(*Feed), nil
}

func (orm *BtOrm) ReadFeeds(from int, size int) ([]*Feed, int64, error) {
	ary2, count, err := orm.Orm.ReadDocuments(&Feed{}, from, size)
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
	fields.Id = ""
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadRule(id common.Ident) (*Rule, error) {
	rule := &Rule{}
	rule.Id = id
	rule2, err := orm.Orm.ReadDocument(rule)
	if err != nil {
		return nil, err
	}
	return rule2.(*Rule), nil
}

func (orm *BtOrm) ReadRules(from int, size int) ([]*Rule, int64, error) {
	ary2, count, err := orm.Orm.ReadDocuments(&Rule{}, from, size)
	if err != nil {
		return nil, 0, err
	}

	ret := make([]*Rule, count)
	for i, v := range ary2 {
		ret[i] = v.(*Rule)
	}
	return ret, count, nil
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
	fields.Id = ""
	fields.OwnerId = requestorID
	return orm.Orm.CreateDocument(fields)
}

func (orm *BtOrm) ReadUser(id common.Ident) (*User, error) {
	user := &User{}
	user.Id = id
	user2, err := orm.Orm.ReadDocument(user)
	if err != nil {
		return nil, err
	}
	return user2.(*User), nil
}

func (orm *BtOrm) ReadUsers(from int, size int) ([]*User, int64, error) {
	ary2, count, err := orm.Orm.ReadDocuments(&User{}, from, size)
	if err != nil {
		return nil, 0, err
	}

	ret := make([]*User, count)
	for i, v := range ary2 {
		ret[i] = v.(*User)
	}
	return ret, count, nil
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
