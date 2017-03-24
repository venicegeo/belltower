package btorm

import (
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type BtOrm struct {
	Orm *esorm.Orm
}

//---------------------------------------------------------------------

func NewBtOrm() (*BtOrm, error) {
	orm, err := esorm.NewOrm()
	if err != nil {
		return nil, err
	}

	btOrm := &BtOrm{
		Orm: orm,
	}

	types := []esorm.Elasticable{
		&Action{},
		&Feed{},
		&Rule{},
		&User{},
	}

	for _, t := range types {
		exists, err := orm.IndexExists(t)
		if err != nil {
			return nil, err
		}

		if exists {
			err = orm.DeleteIndex(t)
			if err != nil {
				return nil, err
			}
		}

		err = orm.CreateIndex(t)
		if err != nil {
			return nil, err
		}
	}

	return btOrm, nil
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateAction(requestorID common.Ident, fields *ActionFieldsForCreate) (common.Ident, error) {
	action := &Action{}
	return orm.Orm.CreateThing(requestorID, action, fields)
}

func (orm *BtOrm) ReadAction(id common.Ident) (*ActionFieldsForRead, error) {
	action := &Action{}
	action.Id = id
	fields, err := orm.Orm.ReadThing(action)
	if err != nil {
		return nil, err
	}
	return fields.(*ActionFieldsForRead), nil
}

func (orm *BtOrm) UpdateAction(id common.Ident, fields *ActionFieldsForUpdate) error {
	action := &Action{}
	action.Id = id
	return orm.Orm.UpdateThing(action, fields)
}

func (orm *BtOrm) DeleteAction(id common.Ident) error {
	action := &Action{}
	action.Id = id
	return orm.Orm.DeleteThing(action)
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateFeed(requestorID common.Ident, fields *FeedFieldsForCreate) (common.Ident, error) {
	feed := &Feed{}
	return orm.Orm.CreateThing(requestorID, feed, fields)
}

func (orm *BtOrm) ReadFeed(id common.Ident) (*FeedFieldsForRead, error) {
	feed := &Feed{}
	feed.Id = id
	fields, err := orm.Orm.ReadThing(feed)
	if err != nil {
		return nil, err
	}
	return fields.(*FeedFieldsForRead), nil
}

func (orm *BtOrm) UpdateFeed(id common.Ident, fields *FeedFieldsForUpdate) error {
	feed := &Feed{}
	feed.Id = id
	return orm.Orm.UpdateThing(feed, fields)
}

func (orm *BtOrm) DeleteFeed(id common.Ident) error {
	feed := &Feed{}
	feed.Id = id
	return orm.Orm.DeleteThing(feed)
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateRule(requestorID common.Ident, fields *RuleFieldsForCreate) (common.Ident, error) {
	rule := &Rule{}
	return orm.Orm.CreateThing(requestorID, rule, fields)
}

func (orm *BtOrm) ReadRule(id common.Ident) (*RuleFieldsForRead, error) {
	rule := &Rule{}
	rule.Id = id
	fields, err := orm.Orm.ReadThing(rule)
	if err != nil {
		return nil, err
	}
	return fields.(*RuleFieldsForRead), nil
}

func (orm *BtOrm) UpdateRule(id common.Ident, fields *RuleFieldsForUpdate) error {
	rule := &Rule{}
	rule.Id = id
	return orm.Orm.UpdateThing(rule, fields)
}

func (orm *BtOrm) DeleteRule(id common.Ident) error {
	rule := &Rule{}
	rule.Id = id
	return orm.Orm.DeleteThing(rule)
}

//---------------------------------------------------------------------

func (orm *BtOrm) CreateUser(requestorID common.Ident, fields *UserFieldsForCreate) (common.Ident, error) {
	user := &User{}
	return orm.Orm.CreateThing(requestorID, user, fields)
}

func (orm *BtOrm) ReadUser(id common.Ident) (*UserFieldsForRead, error) {
	user := &User{}
	user.Id = id
	fields, err := orm.Orm.ReadThing(user)
	if err != nil {
		return nil, err
	}
	return fields.(*UserFieldsForRead), nil
}

func (orm *BtOrm) UpdateUser(id common.Ident, fields *UserFieldsForUpdate) error {
	user := &User{}
	user.Id = id
	return orm.Orm.UpdateThing(user, fields)
}

func (orm *BtOrm) DeleteUser(id common.Ident) error {
	user := &User{}
	user.Id = id
	return orm.Orm.DeleteThing(user)
}
