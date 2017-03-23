package orm2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewBtOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &ActionFieldsForCreate{Name: "one"}
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
	u := &ActionFieldsForUpdate{Name: "two"}
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

	orm, err := NewBtOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &FeedFieldsForCreate{Name: "one"}
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
	u := &FeedFieldsForUpdate{Name: "two"}
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

	orm, err := NewBtOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &RuleFieldsForCreate{Name: "one"}
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
	u := &RuleFieldsForUpdate{Name: "two"}
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

	orm, err := NewBtOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	// does create work?
	c := &UserFieldsForCreate{Name: "one"}
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
	u := &UserFieldsForUpdate{Name: "two"}
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
