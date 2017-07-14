package btorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func Reset(t *testing.T) {
	assert := assert.New(t)
	err := DatabaseInit()
	assert.NoError(err)
}

func TestDatabaseOpen(t *testing.T) {
	assert := assert.New(t)

	var err error
	var id common.Ident
	var action *Action
	var orm *BtOrm

	Reset(t)

	// add one, make sure it is there
	{
		orm := &BtOrm{}
		err = orm.Open()
		assert.NoError(err)

		action := &Action{}
		action.Name = "one"
		id, err = orm.CreateAction("u", action)
		assert.NoError(err)
		assert.NotEmpty(id)

		action, err = orm.ReadAction(id)
		assert.NoError(err)
		assert.NotNil(action)
		assert.EqualValues(id, action.Id)
		assert.EqualValues("one", action.Name)

		err = orm.Close()
		assert.NoError(err)
	}

	// verify still there
	{
		orm = &BtOrm{}
		err = orm.Open()
		assert.NoError(err)

		action, err = orm.ReadAction(id)
		assert.NoError(err)
		assert.NotNil(action)
		assert.EqualValues(id, action.Id)
		assert.EqualValues("one", action.Name)

		err = orm.Close()
		assert.NoError(err)
	}

	// wipe it all
	Reset(t)

	// verify not there anymore
	{
		orm = &BtOrm{}
		err = orm.Open()
		assert.NoError(err)

		action, err = orm.ReadAction(id)
		assert.Error(err)

		err = orm.Close()
		assert.NoError(err)
	}
}

func TestActionCRUD(t *testing.T) {
	assert := assert.New(t)

	Reset(t)

	orm := &BtOrm{}
	err := orm.Open()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &Action{}
	c.Name = "one"
	id, err := orm.CreateAction("u", c)
	assert.NoError(err)
	assert.NotEmpty(id)

	// does read work?
	r, err := orm.ReadAction(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("one", r.Name)

	// update it
	u := &Action{}
	u.Name = "two"
	err = orm.UpdateAction(id, u)
	assert.NoError(err)

	// read again, to check
	r, err = orm.ReadAction(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("two", r.Name)

	{
		// does create work again?
		c2 := &Action{}
		c2.Name = "oneone"
		id2, err := orm.CreateAction("uu", c2)
		assert.NoError(err)
		assert.NotEmpty(id2)

		// read again, to check
		r, err = orm.ReadAction(id2)
		assert.NoError(err)
		assert.NotNil(r)
		assert.EqualValues(id2, r.Id)
		assert.EqualValues("oneone", r.Name)

		// check bulk read
		ary, cnt, err := orm.ReadActions(0, 2)
		assert.NoError(err)
		assert.EqualValues(2, cnt)
		assert.Len(ary, 2)
		assert.EqualValues("two", ary[0].Name)
	}

	// try delete
	err = orm.DeleteAction(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadAction(id)
	assert.Error(err)
}

func TestFeedCRUD(t *testing.T) {
	assert := assert.New(t)

	Reset(t)

	orm := &BtOrm{}
	err := orm.Open()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &Feed{}
	c.Name = "one"
	id, err := orm.CreateFeed("u", c)
	assert.NoError(err)
	assert.NotEmpty(id)

	// does read work?
	r, err := orm.ReadFeed(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("one", r.Name)

	// update it
	u := &Feed{}
	u.Name = "two"
	err = orm.UpdateFeed(id, u)
	assert.NoError(err)

	// read again, to check
	r, err = orm.ReadFeed(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("two", r.Name)

	{
		// does create work again?
		c2 := &Feed{}
		c2.Name = "oneone"
		id2, err := orm.CreateFeed("uu", c2)
		assert.NoError(err)
		assert.NotEmpty(id2)

		// read again, to check
		r, err = orm.ReadFeed(id2)
		assert.NoError(err)
		assert.NotNil(r)
		assert.EqualValues(id2, r.Id)
		assert.EqualValues("oneone", r.Name)

		// check bulk read
		ary, cnt, err := orm.ReadFeeds(0, 2)
		assert.NoError(err)
		assert.EqualValues(2, cnt)
		assert.Len(ary, 2)
		assert.EqualValues("two", ary[0].Name)
	}

	// try delete
	err = orm.DeleteFeed(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadFeed(id)
	assert.Error(err)
}

func TestRuleCRUD(t *testing.T) {
	assert := assert.New(t)

	Reset(t)

	orm := &BtOrm{}
	err := orm.Open()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &Rule{}
	c.Name = "one"
	id, err := orm.CreateRule("u", c)
	assert.NoError(err)
	assert.NotEmpty(id)

	// does read work?
	r, err := orm.ReadRule(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("one", r.Name)

	// update it
	u := &Rule{}
	u.Name = "two"
	err = orm.UpdateRule(id, u)
	assert.NoError(err)

	// read again, to check
	r, err = orm.ReadRule(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("two", r.Name)

	{
		// does create work again?
		c2 := &Rule{}
		c2.Name = "oneone"
		id2, err := orm.CreateRule("uu", c2)
		assert.NoError(err)
		assert.NotEmpty(id2)

		// read again, to check
		r, err = orm.ReadRule(id2)
		assert.NoError(err)
		assert.NotNil(r)
		assert.EqualValues(id2, r.Id)
		assert.EqualValues("oneone", r.Name)

		// check bulk read
		ary, cnt, err := orm.ReadRules(0, 2)
		assert.NoError(err)
		assert.EqualValues(2, cnt)
		assert.Len(ary, 2)
		assert.EqualValues("two", ary[0].Name)
	}

	// try delete
	err = orm.DeleteRule(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadRule(id)
	assert.Error(err)
}

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)

	Reset(t)

	orm := &BtOrm{}
	err := orm.Open()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &User{}
	c.Name = "one"
	id, err := orm.CreateUser("u", c)
	assert.NoError(err)
	assert.NotEmpty(id)

	// does read work?
	r, err := orm.ReadUser(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("one", r.Name)

	// update it
	u := &User{}
	u.Name = "two"
	err = orm.UpdateUser(id, u)
	assert.NoError(err)

	// read again, to check
	r, err = orm.ReadUser(id)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.Id)
	assert.EqualValues("two", r.Name)

	{
		// does create work again?
		c2 := &User{}
		c2.Name = "oneone"
		id2, err := orm.CreateUser("uu", c2)
		assert.NoError(err)
		assert.NotEmpty(id2)

		// read again, to check
		r, err = orm.ReadUser(id2)
		assert.NoError(err)
		assert.NotNil(r)
		assert.EqualValues(id2, r.Id)
		assert.EqualValues("oneone", r.Name)

		// check bulk read
		ary, cnt, err := orm.ReadUsers(0, 2)
		assert.NoError(err)
		assert.EqualValues(2, cnt)
		assert.Len(ary, 2)
		assert.EqualValues("two", ary[0].Name)
	}

	// try delete
	err = orm.DeleteUser(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadUser(id)
	assert.Error(err)
}
