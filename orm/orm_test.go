package orm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var secs2 time.Duration = time.Second * 2
var secs10 time.Duration = time.Second * 10

func TestDBOperations(t *testing.T) {
	assert := assert.New(t)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	db := model.db

	now := time.Now()

	createFields := &FeedFieldsForCreate{
		Name:        "TestFeed",
		FeedType:    "myfeedtype",
		IsEnabled:   false,
		MessageDecl: map[string]interface{}{},
		Settings: map[string]interface{}{
			"alpha": "figgy",
		},
	}

	feed, err := CreateFeed(createFields)
	assert.NoError(err)

	err = db.Create(feed).Error
	assert.NoError(err)

	err = db.First(&feed, 1).Error
	assert.NoError(err)
	assert.Equal("TestFeed", feed.Name)

	assert.True(now.Before(feed.CreatedAt))

	err = db.First(&feed, "name = ?", "xyzzy").Error
	assert.Equal("record not found", err.Error())
	assert.Equal("TestFeed", feed.Name)

	err = db.Model(&feed).Update("Name", "TFT").Error
	assert.NoError(err)
	assert.Equal("TFT", feed.Name)

	err = db.First(&feed, 1).Error
	assert.NoError(err)
	assert.Equal("TFT", feed.Name)

	err = db.Delete(&feed).Error
	assert.NoError(err)
	assert.Equal("TFT", feed.Name)

	err = db.First(&feed, 1).Error
	assert.Error(err)
}

func TestUser(t *testing.T) {
	assert := assert.New(t)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	const ownerID uint = 1999
	var itemID uint

	fields := &UserFieldsForCreate{
		Name:           "Bob",
		IsAdmin:        false,
		IsEnabled:      true,
		OwnerID:        ownerID,
		PublicCanRead:  true,
		PublicCanWrite: false,
	}

	now := time.Now()

	{ // create the user
		itemID, err = model.CreateUser(ownerID, fields)
		assert.NoError(err)
		readFields, err := model.ReadUser(ownerID, itemID)
		assert.NoError(err)

		assert.True(readFields.IsEnabled)
		assert.False(readFields.IsAdmin)
		assert.EqualValues("Bob", readFields.Name)
		assert.Equal(itemID, readFields.ID)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.Equal(readFields.CreatedAt, readFields.UpdatedAt)
		assert.Zero(readFields.LastLoginAt)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(true, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{ // update it
		updateFields := &UserFieldsForUpdate{
			IsEnabled:      false,
			IsAdmin:        true,
			Name:           "Alice",
			PublicCanRead:  false,
			PublicCanWrite: true,
		}
		err = model.UpdateUser(ownerID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadUser(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.EqualValues("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.Zero(readFields.LastLoginAt)
		assert.False(readFields.IsEnabled)
		assert.True(readFields.IsAdmin)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(true, readFields.PublicCanWrite)
	}

	{ // update with default payload
		f := &UserFieldsForUpdate{
		// admin, isEnabled default to false
		}
		err = model.UpdateUser(ownerID, itemID, f)
		assert.NoError(err)

		readFields, err := model.ReadUser(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.EqualValues("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.Zero(readFields.LastLoginAt)
		assert.False(readFields.IsEnabled)
		assert.False(readFields.IsAdmin)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		err = model.DeleteUser(ownerID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadUser(ownerID, itemID)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	err = model.DeleteUser(ownerID, 20169)
	assert.Error(err)
}

func TestFeed(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	now := time.Now()
	const ownerID uint = 1888
	var itemID uint

	createFields := &FeedFieldsForCreate{
		Name:           "Bob",
		FeedType:       "myfeedtype",
		IsEnabled:      true,
		MessageDecl:    map[string]interface{}{},
		Settings:       map[string]interface{}{},
		OwnerID:        ownerID,
		PublicCanRead:  true,
		PublicCanWrite: false,
	}

	{
		itemID, err = model.CreateFeed(ownerID, createFields)
		assert.NoError(err)
		readFields, err := model.ReadFeed(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.Equal("Bob", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.EqualValues("myfeedtype", readFields.FeedType)
		assert.True(readFields.IsEnabled)
		// MessageDecl TODO
		// Settings // TODO
		assert.Zero(readFields.MessageCount)
		assert.Zero(readFields.LastMessageAt)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(true, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		updateFields := &FeedFieldsForUpdate{
			Name:           "Alice",
			IsEnabled:      false,
			PublicCanRead:  false,
			PublicCanWrite: true,
		}
		err = model.UpdateFeed(ownerID, itemID, updateFields)
		assert.NoError(err)
		readFields, err := model.ReadFeed(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.Equal("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.EqualValues("myfeedtype", readFields.FeedType)
		assert.False(readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
		// MessageDecl TODO
		// Settings // TODO
		assert.Zero(readFields.MessageCount)
		assert.Zero(readFields.LastMessageAt)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(true, readFields.PublicCanWrite)
	}

	{ // update with empty object
		updateFields := &FeedFieldsForUpdate{}
		err = model.UpdateFeed(ownerID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadFeed(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.Equal("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.EqualValues("myfeedtype", readFields.FeedType)
		assert.False(readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
		// MessageDecl TODO
		// Settings // TODO
		assert.Zero(readFields.MessageCount)
		assert.Zero(readFields.LastMessageAt)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		err = model.DeleteFeed(ownerID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadFeed(ownerID, itemID)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteFeed(ownerID, 20169)
		assert.Error(err)
	}
}

func TestRule(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	const ownerID uint = 1777
	var itemID uint

	createFields := &RuleFieldsForCreate{
		Name:           "Bob",
		IsEnabled:      true,
		Expression:     "a+b",
		OwnerID:        ownerID,
		PollDuration:   time.Duration(secs2),
		PublicCanRead:  true,
		PublicCanWrite: false,
	}

	{
		itemID, err = model.CreateRule(ownerID, createFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(ownerID, itemID)
		assert.NoError(err)

		assert.Equal("Bob", readFields.Name)
		assert.True(readFields.IsEnabled)
		assert.Equal("a+b", readFields.Expression)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(secs2, readFields.PollDuration)
		assert.Equal(true, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		updateFields := &RuleFieldsForUpdate{
			Name:           "alice",
			PollDuration:   time.Duration(secs10),
			IsEnabled:      false,
			Expression:     "x-y/z",
			PublicCanRead:  false,
			PublicCanWrite: true,
		}
		err = model.UpdateRule(ownerID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(ownerID, itemID)
		assert.NoError(err)

		assert.Equal("alice", readFields.Name)
		assert.False(readFields.IsEnabled)
		assert.Equal("x-y/z", readFields.Expression)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(secs10, readFields.PollDuration)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(true, readFields.PublicCanWrite)
	}

	{
		updateFields := &RuleFieldsForUpdate{}
		err = model.UpdateRule(ownerID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(ownerID, itemID)
		assert.NoError(err)

		assert.Equal("alice", readFields.Name)
		assert.False(readFields.IsEnabled)
		assert.Equal("x-y/z", readFields.Expression)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(secs10, readFields.PollDuration)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		err = model.DeleteRule(ownerID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadRule(ownerID, itemID)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteRule(ownerID, 20169)
		assert.Error(err)
	}
}

func TestAction(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	const ownerID uint = 1666
	var itemID uint

	createFields := &ActionFieldsForCreate{
		Name:           "Bob",
		IsEnabled:      true,
		Settings:       map[string]interface{}{"a": "b"},
		OwnerID:        ownerID,
		PublicCanRead:  true,
		PublicCanWrite: false,
	}

	{
		itemID, err = model.CreateAction(ownerID, createFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(ownerID, itemID)
		assert.NoError(err)
		assert.Equal(true, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(true, readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(true, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		updateFields := &ActionFieldsForUpdate{
			Name:           "Alice",
			IsEnabled:      false,
			PublicCanRead:  false,
			PublicCanWrite: true,
		}
		err = model.UpdateAction(ownerID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(false, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(false, readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(true, readFields.PublicCanWrite)
	}

	{
		updateFields := &ActionFieldsForUpdate{}
		err = model.UpdateAction(ownerID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(ownerID, itemID)
		assert.NoError(err)

		assert.Equal(false, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(false, readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
		assert.Equal(false, readFields.PublicCanRead)
		assert.Equal(false, readFields.PublicCanWrite)
	}

	{
		err = model.DeleteAction(ownerID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadAction(ownerID, itemID)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteAction(ownerID, 20169)
		assert.Error(err)
	}
}
