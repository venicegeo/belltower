package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBOperations(t *testing.T) {
	assert := assert.New(t)
	var err error

	model, err := NewOrm()
	assert.NoError(err)
	defer model.Close()

	db := model.db

	feedType := &FeedType{Name: "TestFeedType", ConfigInfo: "figgy"}

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

func TestFeedTypeAcls(t *testing.T) {
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

	{
		err = model.AddUser("Bob", user)
		assert.NoError(err)
		u, err := model.GetUserByName("Bob")
		assert.NoError(err)
		assert.True(u.IsEnabled)
		assert.EqualValues(u.UserAttributes, *user)

		user.IsEnabled = false
		err = model.UpdateUser("Bob", user)
		assert.NoError(err)

		u, err = model.GetUserByName("Bob")
		assert.NoError(err)
		assert.False(u.IsEnabled)
		assert.EqualValues(u.UserAttributes, *user)
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
}

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
