package orm

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func newDatabase() *gorm.DB {
	// ignore errors
	os.Remove("test.db")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func TestFeedType(t *testing.T) {
	assert := assert.New(t)
	var err error

	db := newDatabase()
	defer db.Close()

	err = db.CreateTable(&FeedType{}).Error
	assert.NoError(err)

	err = db.Create(&FeedType{Name: "ft1"}).Error
	assert.NoError(err)

	feedType := FeedType{}

	err = db.First(&feedType, 1).Error
	assert.NoError(err)
	assert.Equal("ft1", feedType.Name)

	err = db.First(&feedType, "name = ?", "xyzzy").Error
	assert.Equal("record not found", err.Error())
	assert.Equal("ft1", feedType.Name)

	err = db.Model(&feedType).Update("Name", "ft99").Error
	assert.NoError(err)
	assert.Equal("ft99", feedType.Name)

	err = db.First(&feedType, 1).Error
	assert.NoError(err)
	assert.Equal("ft99", feedType.Name)

	err = db.Delete(&feedType).Error
	assert.NoError(err)
	assert.Equal("ft99", feedType.Name)

	err = db.First(&feedType, 1).Error
	assert.Error(err)
}

func TestFeedTypeAcls(t *testing.T) {
	assert := assert.New(t)
	var err error

	db := newDatabase()
	defer db.Close()

	u0 := User{Name: "Admin"}
	u1 := User{Name: "Alpha"}
	u2 := User{Name: "Beta"}
	u3 := User{Name: "Gamma"}
	{
		err = db.CreateTable(&User{}).Error
		assert.NoError(err)
		err = db.Create(&u1).Error
		assert.NoError(err)
		err = db.Create(&u2).Error
		assert.NoError(err)
		err = db.Create(&u3).Error
		assert.NoError(err)
		err = db.Create(&u0).Error
		assert.NoError(err)
	}

	ft1 := FeedType{
		Name:    "S3FeedType",
		OwnerID: u0.ID,
	}

	ft2 := FeedType{
		Name:    "FileSysFeedType",
		OwnerID: u0.ID,
	}
	{
		err = db.CreateTable(&FeedType{}).Error
		assert.NoError(err)

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
		err = db.CreateTable(&Feed{}).Error
		assert.NoError(err)

		err = db.Create(&feed1).Error
		assert.NoError(err)
		err = db.Create(&feed2).Error
		assert.NoError(err)
	}

	{
		err = db.CreateTable(&FeedTypeAcl{}).Error
		assert.NoError(err)

		err = db.Create(&FeedTypeAcl{UserID: u1.ID, FeedTypeID: ft1.ID}).Error
		assert.NoError(err)
		err = db.Create(&FeedTypeAcl{UserID: u2.ID, FeedTypeID: ft2.ID}).Error
		assert.NoError(err)
		err = db.Create(&FeedTypeAcl{UserID: u3.ID, FeedTypeID: ft2.ID}).Error
		assert.NoError(err)
	}

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
		var acls []FeedTypeAcl
		err = db.Find(&acls).Error
		assert.NoError(err)
		for _, v := range acls {
			fmt.Printf("%s", v)
		}
	}
}
