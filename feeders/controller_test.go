package feeders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/btorm"
)

func TestController(t *testing.T) {
	assert := assert.New(t)

	{
		err := btorm.DatabaseInit()
		assert.NoError(err)
	}

	{
		orm := &btorm.BtOrm{}
		err := orm.Open()
		assert.NoError(err)
		assert.NotNil(orm)

		feed := &btorm.Feed{}
		feed.Name = "mysimplefeed"
		feed.FeederId = SimpleFeederId
		feed.PollingInterval = 1
		feed.PollingEndAt = time.Now().Add(3 * time.Second)
		feed.Settings = map[string]string{
			"Value": "7",
		}

		id, err := orm.CreateFeed("root", feed)
		assert.NoError(err)
		assert.NotEmpty(id)

		{
			// does read work?
			r, err := orm.ReadFeed(id)
			assert.NoError(err)
			assert.NotNil(r)
			assert.EqualValues(id, r.Id)
			assert.EqualValues("mysimplefeed", r.Name)
		}

		err = orm.Close()
		assert.NoError(err)
	}

	ctl := &Controller{}

	err := ctl.Start()
	assert.NoError(err)
}

need to test error handling in run loop
