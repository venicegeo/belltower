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
	fields := &UserFieldsForCreate{
		Name:      "Bob",
		IsAdmin:   false,
		IsEnabled: true,
		OwnerID:   ownerID,
	}

	var id uint

	now := time.Now()

	{ // create the user
		id, err = model.CreateUser(fields)
		assert.NoError(err)
		readFields, err := model.ReadUser(id)
		assert.NoError(err)

		assert.True(readFields.IsEnabled)
		assert.False(readFields.IsAdmin)
		assert.EqualValues("Bob", readFields.Name)
		assert.Equal(id, readFields.ID)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.Equal(readFields.CreatedAt, readFields.UpdatedAt)
		assert.Zero(readFields.LastLoginAt)
		assert.Equal(ownerID, readFields.OwnerID)
	}

	{ // update it
		updateFields := &UserFieldsForUpdate{
			IsEnabled: false,
			IsAdmin:   true,
			Name:      "Alice",
		}
		err = model.UpdateUser(id, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadUser(id)
		assert.NoError(err)

		assert.Equal(id, readFields.ID)
		assert.EqualValues("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.Zero(readFields.LastLoginAt)
		assert.False(readFields.IsEnabled)
		assert.True(readFields.IsAdmin)
		assert.Equal(ownerID, readFields.OwnerID)
	}

	{ // update with default payload
		f := &UserFieldsForUpdate{
		// admin, isEnabled default to false
		}
		err = model.UpdateUser(id, f)
		assert.NoError(err)

		readFields, err := model.ReadUser(id)
		assert.NoError(err)

		assert.Equal(id, readFields.ID)
		assert.EqualValues("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.Zero(readFields.LastLoginAt)
		assert.False(readFields.IsEnabled)
		assert.False(readFields.IsAdmin)
		assert.Equal(ownerID, readFields.OwnerID)
	}

	{
		err = model.DeleteUser(id)
		assert.NoError(err)
		readFields, err := model.ReadUser(id)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	err = model.DeleteUser(20169)
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

	createFields := &FeedFieldsForCreate{
		Name:        "Bob",
		FeedType:    "myfeedtype",
		IsEnabled:   true,
		MessageDecl: map[string]interface{}{},
		Settings:    map[string]interface{}{},
		OwnerID:     ownerID,
	}

	var id uint

	{
		id, err = model.CreateFeed(createFields)
		assert.NoError(err)
		readFields, err := model.ReadFeed(id)
		assert.NoError(err)

		assert.Equal(id, readFields.ID)
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
	}

	{
		updateFields := &FeedFieldsForUpdate{
			Name:      "Alice",
			IsEnabled: false,
		}
		err = model.UpdateFeed(id, updateFields)
		assert.NoError(err)
		readFields, err := model.ReadFeed(id)
		assert.NoError(err)

		assert.Equal(id, readFields.ID)
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
	}

	{ // update with empty object
		updateFields := &FeedFieldsForUpdate{}
		err = model.UpdateFeed(id, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadFeed(id)
		assert.NoError(err)

		assert.Equal(id, readFields.ID)
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
	}

	{
		err = model.DeleteFeed(id)
		assert.NoError(err)
		readFields, err := model.ReadFeed(id)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteFeed(20169)
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

	//now := time.Now()

	createFields := &RuleFieldsForCreate{
		Name:         "Bob",
		IsEnabled:    true,
		Expression:   "a+b",
		OwnerID:      777,
		PollDuration: time.Duration(secs2),
	}

	var id uint

	{
		id, err = model.CreateRule(createFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(id)
		assert.NoError(err)

		assert.Equal("Bob", readFields.Name)
		assert.True(readFields.IsEnabled)
		assert.Equal("a+b", readFields.Expression)
		assert.Equal(uint(777), readFields.OwnerID)
		assert.Equal(secs2, readFields.PollDuration)
	}

	{
		updateFields := &RuleFieldsForUpdate{
			Name:         "alice",
			PollDuration: time.Duration(secs10),
			IsEnabled:    false,
			Expression:   "x-y/z",
		}
		err = model.UpdateRule(id, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(id)
		assert.NoError(err)

		assert.Equal("alice", readFields.Name)
		assert.False(readFields.IsEnabled)
		assert.Equal("x-y/z", readFields.Expression)
		assert.Equal(uint(777), readFields.OwnerID)
		assert.Equal(secs10, readFields.PollDuration)
	}

	{
		updateFields := &RuleFieldsForUpdate{}
		err = model.UpdateRule(id, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(id)
		assert.NoError(err)

		assert.Equal("alice", readFields.Name)
		assert.False(readFields.IsEnabled)
		assert.Equal("x-y/z", readFields.Expression)
		assert.Equal(uint(777), readFields.OwnerID)
		assert.Equal(secs10, readFields.PollDuration)
	}

	{
		err = model.DeleteRule(id)
		assert.NoError(err)
		readFields, err := model.ReadRule(id)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteRule(20169)
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

	//now := time.Now()

	createFields := &ActionFieldsForCreate{
		Name:      "Bob",
		IsEnabled: true,
		Settings:  map[string]interface{}{"a": "b"},
		OwnerID:   ownerID,
	}

	var id uint

	{
		id, err = model.CreateAction(createFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(id)
		assert.NoError(err)
		assert.Equal(true, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(true, readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
	}

	{
		updateFields := &ActionFieldsForUpdate{
			Name:      "Alice",
			IsEnabled: false,
		}
		err = model.UpdateAction(id, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(id)
		assert.NoError(err)

		assert.Equal(false, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(false, readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
	}

	{
		updateFields := &ActionFieldsForUpdate{}
		err = model.UpdateAction(id, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(id)
		assert.NoError(err)

		assert.Equal(false, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(false, readFields.IsEnabled)
		assert.Equal(ownerID, readFields.OwnerID)
	}

	{
		err = model.DeleteAction(id)
		assert.NoError(err)
		readFields, err := model.ReadAction(id)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteAction(20169)
		assert.Error(err)
	}
}
