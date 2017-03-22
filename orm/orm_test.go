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

	requestorID := model.AdminID

	now := time.Now()

	createFields := &FeedFieldsForCreate{
		Name:      "TestFeed",
		FeedType:  "900",
		IsEnabled: false,
		Settings: map[string]interface{}{
			"alpha": "figgy",
		},
	}

	feed, err := CreateFeed(requestorID, createFields)
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

	requestorID := model.AdminID

	var itemID uint

	fields := &UserFieldsForCreate{
		Name:      "Bob",
		Role:      AdminRole,
		IsEnabled: true,
	}

	now := time.Now()

	{ // create the user
		itemID, err = model.CreateUser(requestorID, fields)
		assert.NoError(err)
		readFields, err := model.ReadUser(requestorID, itemID)
		assert.NoError(err)

		assert.True(readFields.IsEnabled)
		assert.Equal(AdminRole, readFields.Role)
		assert.EqualValues("Bob", readFields.Name)
		assert.Equal(itemID, readFields.ID)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.Equal(readFields.CreatedAt, readFields.UpdatedAt)
		assert.Zero(readFields.LastLoginAt)
	}

	{ // update it
		updateFields := &UserFieldsForUpdate{
			IsEnabled: false,
			Name:      "Alice",
			Role:      CreatorRole,
		}
		err = model.UpdateUser(requestorID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadUser(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.EqualValues("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.Zero(readFields.LastLoginAt)
		assert.False(readFields.IsEnabled)
		assert.Equal(CreatorRole, readFields.Role)
	}

	{ // update with default payload
		f := &UserFieldsForUpdate{
		// admin, isEnabled default to false
		}
		err = model.UpdateUser(requestorID, itemID, f)
		assert.NoError(err)

		readFields, err := model.ReadUser(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.EqualValues("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.Zero(readFields.LastLoginAt)
		assert.False(readFields.IsEnabled)
		assert.Equal(UserRole, readFields.Role)
	}

	{
		err = model.DeleteUser(model.AdminID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadUser(requestorID, itemID)
		assert.NoError(err)
		assert.Nil(readFields)
	}

	err = model.DeleteUser(requestorID, 20169)
	assert.Error(err)
}

func TestFeed(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	now := time.Now()

	requestorID := model.AdminID

	var itemID uint

	createFields := &FeedFieldsForCreate{
		Name:      "Bob",
		FeedType:  "901",
		IsEnabled: true,
		Settings:  map[string]interface{}{},
		IsPublic:  true,
	}

	{
		itemID, err = model.CreateFeed(requestorID, createFields)
		assert.NoError(err)
		readFields, err := model.ReadFeed(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.Equal("Bob", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.EqualValues("901", readFields.FeedType)
		assert.True(readFields.IsEnabled)
		// MessageDecl TODO
		// Settings // TODO
		assert.Zero(readFields.MessageCount)
		assert.Zero(readFields.LastMessageAt)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(true, readFields.IsPublic)
	}

	{
		updateFields := &FeedFieldsForUpdate{
			Name:      "Alice",
			IsEnabled: false,
			IsPublic:  false,
		}
		err = model.UpdateFeed(requestorID, itemID, updateFields)
		assert.NoError(err)
		readFields, err := model.ReadFeed(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.Equal("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.EqualValues("901", readFields.FeedType)
		assert.False(readFields.IsEnabled)
		assert.Equal(requestorID, readFields.OwnerID)
		// MessageDecl TODO
		// Settings // TODO
		assert.Zero(readFields.MessageCount)
		assert.Zero(readFields.LastMessageAt)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(false, readFields.IsPublic)
	}

	{ // update with empty object
		updateFields := &FeedFieldsForUpdate{}
		err = model.UpdateFeed(requestorID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadFeed(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(itemID, readFields.ID)
		assert.Equal("Alice", readFields.Name)
		assert.WithinDuration(now, readFields.CreatedAt, secs2)
		assert.WithinDuration(now, readFields.UpdatedAt, secs2)
		assert.EqualValues("901", readFields.FeedType)
		assert.False(readFields.IsEnabled)
		assert.Equal(requestorID, readFields.OwnerID)
		// MessageDecl TODO
		// Settings // TODO
		assert.Zero(readFields.MessageCount)
		assert.Zero(readFields.LastMessageAt)
		assert.Equal(false, readFields.IsPublic)
	}

	{
		err = model.DeleteFeed(requestorID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadFeed(requestorID, itemID)
		assert.Error(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteFeed(requestorID, 20169)
		assert.Error(err)
	}
}

func TestRule(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	requestorID := model.AdminID

	var itemID uint

	createFields := &RuleFieldsForCreate{
		Name:         "Bob",
		IsEnabled:    true,
		Expression:   "a+b",
		PollDuration: time.Duration(secs2),
		IsPublic:     true,
	}

	{
		itemID, err = model.CreateRule(requestorID, createFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(requestorID, itemID)
		assert.NoError(err)

		assert.Equal("Bob", readFields.Name)
		assert.True(readFields.IsEnabled)
		assert.Equal("a+b", readFields.Expression)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(secs2, readFields.PollDuration)
		assert.Equal(true, readFields.IsPublic)
	}

	{
		updateFields := &RuleFieldsForUpdate{
			Name:         "alice",
			PollDuration: time.Duration(secs10),
			IsEnabled:    false,
			Expression:   "x-y/z",
			IsPublic:     false,
		}
		err = model.UpdateRule(requestorID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(requestorID, itemID)
		assert.NoError(err)

		assert.Equal("alice", readFields.Name)
		assert.False(readFields.IsEnabled)
		assert.Equal("x-y/z", readFields.Expression)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(secs10, readFields.PollDuration)
		assert.Equal(false, readFields.IsPublic)
	}

	{
		updateFields := &RuleFieldsForUpdate{}
		err = model.UpdateRule(requestorID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadRule(requestorID, itemID)
		assert.NoError(err)

		assert.Equal("alice", readFields.Name)
		assert.False(readFields.IsEnabled)
		assert.Equal("x-y/z", readFields.Expression)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(secs10, readFields.PollDuration)
		assert.Equal(false, readFields.IsPublic)
	}

	{
		err = model.DeleteRule(requestorID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadRule(requestorID, itemID)
		assert.Error(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteRule(requestorID, 20169)
		assert.Error(err)
	}
}

func TestAction(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	requestorID := model.AdminID

	var itemID uint

	createFields := &ActionFieldsForCreate{
		Name:      "Bob",
		IsEnabled: true,
		Settings:  map[string]interface{}{"a": "b"},
		IsPublic:  true,
	}

	{
		itemID, err = model.CreateAction(requestorID, createFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(requestorID, itemID)
		assert.NoError(err)
		assert.Equal(true, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(true, readFields.IsEnabled)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(true, readFields.IsPublic)
	}

	{
		updateFields := &ActionFieldsForUpdate{
			Name:      "Alice",
			IsEnabled: false,
			IsPublic:  false,
		}
		err = model.UpdateAction(requestorID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(false, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(false, readFields.IsEnabled)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(false, readFields.IsPublic)
	}

	{
		updateFields := &ActionFieldsForUpdate{}
		err = model.UpdateAction(requestorID, itemID, updateFields)
		assert.NoError(err)

		readFields, err := model.ReadAction(requestorID, itemID)
		assert.NoError(err)

		assert.Equal(false, readFields.IsEnabled)
		assert.EqualValues(map[string]interface{}{"a": "b"}, readFields.Settings)
		assert.Equal(false, readFields.IsEnabled)
		assert.Equal(requestorID, readFields.OwnerID)
		assert.Equal(false, readFields.IsPublic)
	}

	{
		err = model.DeleteAction(requestorID, itemID)
		assert.NoError(err)
		readFields, err := model.ReadAction(requestorID, itemID)
		assert.Error(err)
		assert.Nil(readFields)
	}

	{
		err = model.DeleteAction(requestorID, 20169)
		assert.Error(err)
	}
}

func TestFeedToRuleToAction(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	requestorID := model.AdminID

	createFeedFields := &FeedFieldsForCreate{
		Name:      "Bob",
		FeedType:  "902",
		IsEnabled: true,
		Settings:  map[string]interface{}{},
		IsPublic:  true,
	}

	feedID, err := model.CreateFeed(requestorID, createFeedFields)
	assert.NoError(err)

	createRuleFields := &RuleFieldsForCreate{
		Name:         "Bob",
		IsEnabled:    true,
		Expression:   "a+b",
		PollDuration: time.Duration(secs2),
		IsPublic:     true,
	}

	ruleID, err := model.CreateRule(requestorID, createRuleFields)
	assert.NoError(err)

	createActionFields := &ActionFieldsForCreate{
		Name:      "Bob",
		IsEnabled: true,
		Settings:  map[string]interface{}{"a": "b"},
		IsPublic:  true,
	}

	actionID, err := model.CreateAction(requestorID, createActionFields)
	assert.NoError(err)

	feedToRuleID, err := model.CreateFeedToRule(requestorID, feedID, ruleID, false)
	assert.NoError(err)
	_ = feedToRuleID

	ruleToActionID, err := model.CreateRuleToAction(requestorID, ruleID, actionID, false)
	assert.NoError(err)
	_ = ruleToActionID

	readFeedToRuleFields, err := model.ReadFeedToRule(requestorID, feedToRuleID)
	assert.NoError(err)
	assert.Equal(requestorID, readFeedToRuleFields.OwnerID)
	assert.Equal(feedID, readFeedToRuleFields.FeedID)
	assert.Equal(ruleID, readFeedToRuleFields.RuleID)
	assert.Equal(false, readFeedToRuleFields.IsPublic)

	readRuleToActionFields, err := model.ReadRuleToAction(requestorID, ruleToActionID)
	assert.NoError(err)
	assert.Equal(requestorID, readRuleToActionFields.OwnerID)
	assert.Equal(ruleID, readRuleToActionFields.RuleID)
	assert.Equal(actionID, readRuleToActionFields.ActionID)
	assert.Equal(false, readRuleToActionFields.IsPublic)
}
