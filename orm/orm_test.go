package orm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDBOperations(t *testing.T) {
	assert := assert.New(t)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	db := model.db

	feedType := &FeedType{}
	feedType.Name = "TestFeedType"
	feedType.ConfigInfo = "figgy"

	err = db.Create(feedType).Error
	assert.NoError(err)

	err = db.First(&feedType, 1).Error
	assert.NoError(err)
	assert.Equal("TestFeedType", feedType.Name)

	err = db.First(&feedType, "name = ?", "xyzzy").Error
	assert.Equal("record not found", err.Error())
	assert.Equal("TestFeedType", feedType.Name)

	err = db.Model(&feedType).Update("Name", "TFT").Error
	assert.NoError(err)
	assert.Equal("TFT", feedType.Name)

	err = db.First(&feedType, 1).Error
	assert.NoError(err)
	assert.Equal("TFT", feedType.Name)

	err = db.Delete(&feedType).Error
	assert.NoError(err)
	assert.Equal("TFT", feedType.Name)

	err = db.First(&feedType, 1).Error
	assert.Error(err)
}

func TestUser(t *testing.T) {
	assert := assert.New(t)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	user := &UserAttributes{
		Name:      "Bob",
		IsAdmin:   true,
		IsEnabled: true,
	}

	id, err := model.AddUser(user)
	assert.NoError(err)
	u, err := model.GetUser(id)
	assert.NoError(err)
	assert.True(u.IsEnabled)
	assert.EqualValues(u.UserAttributes, *user)

	user.IsEnabled = false
	err = model.UpdateUser(id, user)
	assert.NoError(err)

	u, err = model.GetUser(id)
	assert.NoError(err)
	assert.False(u.IsEnabled)
	assert.EqualValues(u.UserAttributes, *user)

	err = model.DeleteUser(id)
	assert.NoError(err)
	u, err = model.GetUser(id)
	assert.NoError(err)
	assert.Nil(u)

	err = model.DeleteUser(20169)
	assert.Error(err)
}

func TestFeedType(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	item := &FeedTypeAttributes{
		Name:       "Bob",
		ConfigInfo: "...",
		IsEnabled:  true,
	}

	id, err := model.AddFeedType(item)
	assert.NoError(err)
	t, err := model.GetFeedType(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(t.FeedTypeAttributes, *item)

	item.IsEnabled = false
	err = model.UpdateFeedType(id, item)
	assert.NoError(err)

	t, err = model.GetFeedType(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)
	assert.EqualValues(t.FeedTypeAttributes, *item)

	err = model.DeleteFeedType(id)
	assert.NoError(err)
	t, err = model.GetFeedType(id)
	assert.NoError(err)
	assert.Nil(t)

	err = model.DeleteFeedType(20169)
	assert.Error(err)
}

func TestFeed(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	item := &FeedAttributes{
		Name:                "Bob",
		PersistenceDuration: time.Duration(0),
		IsEnabled:           true,
	}

	id, err := model.AddFeed(item)
	assert.NoError(err)
	t, err := model.GetFeed(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(t.FeedAttributes, *item)

	item.IsEnabled = false
	err = model.UpdateFeed(id, item)
	assert.NoError(err)

	t, err = model.GetFeed(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)
	assert.EqualValues(t.FeedAttributes, *item)

	err = model.DeleteFeed(id)
	assert.NoError(err)
	t, err = model.GetFeed(id)
	assert.NoError(err)
	assert.Nil(t)

	err = model.DeleteFeed(20169)
	assert.Error(err)
}

func TestRule(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	item := &RuleAttributes{
		Name:         "Bob",
		PollDuration: time.Duration(0),
		IsEnabled:    true,
		Expression:   "expr expr",
	}

	id, err := model.AddRule(item)
	assert.NoError(err)
	t, err := model.GetRule(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(t.RuleAttributes, *item)

	item.IsEnabled = false
	err = model.UpdateRule(id, item)
	assert.NoError(err)

	t, err = model.GetRule(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)
	assert.EqualValues(t.RuleAttributes, *item)

	err = model.DeleteRule(id)
	assert.NoError(err)
	t, err = model.GetRule(id)
	assert.NoError(err)
	assert.Nil(t)

	err = model.DeleteRule(20169)
	assert.Error(err)
}

func TestAction(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	item := &ActionAttributes{
		Name:      "Bob",
		IsEnabled: true,
	}

	id, err := model.AddAction(item)
	assert.NoError(err)
	t, err := model.GetAction(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(t.ActionAttributes, *item)

	item.IsEnabled = false
	err = model.UpdateAction(id, item)
	assert.NoError(err)

	t, err = model.GetAction(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)
	assert.EqualValues(t.ActionAttributes, *item)

	err = model.DeleteAction(id)
	assert.NoError(err)
	t, err = model.GetAction(id)
	assert.NoError(err)
	assert.Nil(t)

	err = model.DeleteAction(20169)
	assert.Error(err)
}

func TestActionType(tst *testing.T) {
	assert := assert.New(tst)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	item := &ActionTypeAttributes{
		Name:      "Bob",
		IsEnabled: true,
	}

	id, err := model.AddActionType(item)
	assert.NoError(err)
	t, err := model.GetActionType(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(t.ActionTypeAttributes, *item)

	item.IsEnabled = false
	err = model.UpdateActionType(id, item)
	assert.NoError(err)

	t, err = model.GetActionType(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)
	assert.EqualValues(t.ActionTypeAttributes, *item)

	err = model.DeleteActionType(id)
	assert.NoError(err)
	t, err = model.GetActionType(id)
	assert.NoError(err)
	assert.Nil(t)

	err = model.DeleteActionType(20169)
	assert.Error(err)
}

/*
	db := model.db

	ft1 := FeedType{
		Name: "S3FeedType",
	}

	ft2 := FeedType{
		Name: "FileSysFeedType",
	}
	{
		err = db.Create(&ft1).Error
		assert.NoError(err)
		err = db.Create(&ft2).Error
		assert.NoError(err)
	}

	feed1 := Feed{
		Name:    "BeachfrontFeed",
		OwnerID: u0.ID,
	}

	feed2 := Feed{
		Name:    "NetsFeed",
		OwnerID: u0.ID,
	}

	{
		err = db.Create(&feed1).Error
		assert.NoError(err)
		err = db.Create(&feed2).Error
		assert.NoError(err)
	}

	{
		err = db.Create(&FeedAccessList{UserID: u1.ID, FeedID: ft1.ID}).Error
		assert.NoError(err)
		err = db.Create(&FeedAccessList{UserID: u2.ID, FeedID: ft2.ID}).Error
		assert.NoError(err)
	}
*/

/*
	   	{
	   		var fts []FeedType
	   		err = db.Find(&fts).Error
	   		assert.NoError(err)
	   		for _, v := range fts {
	   			fmt.Printf("%s", v)
	   		}
	   	}

	   	{
	   		var us []User
	   		err = db.Find(&us).Error
	   		assert.NoError(err)
	   		for _, v := range us {
	   			fmt.Printf("%s", v)
	   		}
	   	}

	   	ft := FeedType{}
	   	err = db.First(&ft, 1).Error
	   	assert.NoError(err)
	   	assert.Equal("S3FeedType", ft.Name)

	   	var u User
	   	err = db.Model(&ft).Related(&u, "Owner").Error
	   	assert.NoError(err)
	   	assert.Equal(u0.ID, u.ID)
	   	assert.Equal("Admin", u.Name)

	   	{
	   		var acls []FeedAccessList
	   		err = db.Find(&acls).Error
	   		assert.NoError(err)
	   		for _, v := range acls {
	   			fmt.Printf("%s", v)
	   		}
	   	}

}
*/
