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

	now := time.Now()
	feed := &Feed{
		ID:   17,
		Name: "TestFeed",
		Settings: map[string]interface{}{
			"alpha": "figgy",
		},
		LastMessageAt: &now,
	}

	err = db.Create(feed).Error
	assert.NoError(err)

	err = db.First(&feed, 17).Error
	assert.NoError(err)
	assert.Equal("TestFeed", feed.Name)

	err = db.First(&feed, "name = ?", "xyzzy").Error
	assert.Equal("record not found", err.Error())
	assert.Equal("TestFeed", feed.Name)

	err = db.Model(&feed).Update("Name", "TFT").Error
	assert.NoError(err)
	assert.Equal("TFT", feed.Name)

	err = db.First(&feed, 17).Error
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

	now := time.Now()

	user := &User{
		ID:          17,
		Name:        "Bob",
		IsAdmin:     true,
		IsEnabled:   true,
		LastLoginAt: now,
		CreatedAt:   now,
	}

	id, err := model.AddUser(user)
	assert.NoError(err)
	u, err := model.GetUser(id)
	assert.NoError(err)
	assert.True(u.IsEnabled)

	user.IsEnabled = false
	err = model.UpdateUser(id, user)
	assert.NoError(err)

	u, err = model.GetUser(id)
	assert.NoError(err)
	assert.False(u.IsEnabled)

	err = model.DeleteUser(id)
	assert.NoError(err)
	u, err = model.GetUser(id)
	assert.NoError(err)
	assert.Nil(u)

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

	item := &Feed{
		ID:            19,
		Name:          "Bob",
		IsEnabled:     true,
		Settings:      map[string]interface{}{},
		LastMessageAt: &now,
	}

	id, err := model.AddFeed(item)
	assert.NoError(err)
	t, err := model.GetFeed(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(t.Name, item.Name)
	assert.EqualValues(t.IsEnabled, item.IsEnabled)

	item.IsEnabled = false
	err = model.UpdateFeed(id, item)
	assert.NoError(err)

	t, err = model.GetFeed(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)
	assert.EqualValues(t.Name, item.Name)
	assert.EqualValues(t.IsEnabled, item.IsEnabled)

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

	now := time.Now()

	rule := &Rule{
		ID:           123,
		Name:         "Bob",
		PollDuration: time.Duration(0),
		IsEnabled:    true,
		Expression:   "expr expr",
		CreatedAt:    now,
	}

	id, err := model.AddRule(rule)
	assert.NoError(err)
	assert.Equal(uint(123), id)
	r, err := model.GetRule(id)
	assert.NoError(err)
	assert.Equal(uint(123), r.ID)
	assert.True(r.IsEnabled)

	rule.IsEnabled = false
	err = model.UpdateRule(id, rule)
	assert.NoError(err)

	t, err := model.GetRule(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)

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

	item := &Action{
		ID:        65535,
		Name:      "Bob",
		IsEnabled: true,
	}

	id, err := model.AddAction(item)
	assert.NoError(err)
	t, err := model.GetAction(id)
	assert.NoError(err)
	assert.True(t.IsEnabled)
	assert.EqualValues(*t, *item)

	item.IsEnabled = false
	err = model.UpdateAction(id, item)
	assert.NoError(err)

	t, err = model.GetAction(id)
	assert.NoError(err)
	assert.False(t.IsEnabled)

	err = model.DeleteAction(id)
	assert.NoError(err)
	t, err = model.GetAction(id)
	assert.NoError(err)
	assert.Nil(t)

	err = model.DeleteAction(20169)
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
