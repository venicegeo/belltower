package btorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseOpen(t *testing.T) {
	assert := assert.New(t)

	exists := func(orm *BtOrm) {
		exists, err := orm.databaseExists()
		assert.NoError(err)
		assert.True(exists)
	}

	delete := func(orm *BtOrm) {
		err := orm.deleteDatabase()
		assert.NoError(err)
	}

	create := func() *BtOrm {
		orm, err := NewBtOrm("", OrmOptionCreate)
		if err != nil {
			return nil
		}
		return orm
	}

	open := func() *BtOrm {
		orm, err := NewBtOrm("", OrmOptionOpen)
		if err != nil {
			return nil
		}
		return orm
	}

	openOrCreate := func() *BtOrm {
		orm, err := NewBtOrm("", OrmOptionOpenOrCreate)
		if err != nil {
			return nil
		}
		return orm
	}

	close := func(orm *BtOrm) {
		err := orm.Close()
		assert.NoError(err)
	}

	forceClean := func() {
		orm := openOrCreate()
		assert.NotNil(orm)
		delete(orm)
		close(orm)
	}

	var orm *BtOrm

	forceClean()

	// db doesn't exist: open() should fail
	orm = open()
	assert.Nil(orm)

	// db doesn't exist, create should succeed
	orm = create()
	assert.NotNil(orm)
	exists(orm)
	close(orm)

	// db does exist: create() should succeed (by overwriting)
	orm = create()
	assert.NotNil(orm)
	close(orm)

	// db does exist: open() should succeed
	orm = open()
	assert.NotNil(orm)
	exists(orm)
	close(orm)

	forceClean()

	// db doesn't exist, create() should succeed
	orm = openOrCreate()
	assert.NotNil(orm)
	exists(orm)
	close(orm)

	// db does exist, open() should succeed
	orm = openOrCreate()
	assert.NotNil(orm)
	exists(orm)
	delete(orm)
	close(orm)
}

func TestActionCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewBtOrm("", OrmOptionCreate)
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

	// try delete
	err = orm.DeleteAction(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadAction(id)
	assert.Error(err)
}

func TestFeedCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewBtOrm("", OrmOptionCreate)
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

	// try delete
	err = orm.DeleteFeed(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadFeed(id)
	assert.Error(err)
}

func TestRuleCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewBtOrm("", OrmOptionCreate)
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

	// try delete
	err = orm.DeleteRule(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadRule(id)
	assert.Error(err)
}

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewBtOrm("", OrmOptionCreate)
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

	// try delete
	err = orm.DeleteUser(id)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadUser(id)
	assert.Error(err)
}
