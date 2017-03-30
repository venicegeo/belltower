package feeders

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomFeeder(t *testing.T) {
	assert := assert.New(t)

	feed := Feed{
		Id:       "20",
		Name:     "myrandomfeed",
		FeederId: RandomFeederId,
		Interval: 500 * time.Millisecond,
		EndDate:  time.Now().Add(3 * time.Second),
		SettingsValues: RandomSettings{
			Target: 50,
			Seed:   19, // 61, 52, 88, 23, 30, 70
		},
	}

	feeder, err := feederFactory.create(feed)
	assert.NoError(err)

	assert.Equal("integer", feeder.SettingsSchema()["Target"])
	assert.Equal("integer", feeder.SettingsSchema()["Seed"])
	assert.Equal("integer", feeder.EventSchema()["Value"])

	// define the post-event function, which takes a Feeder's event object
	hits := 0
	poster := func(event *Event) error {
		hits++
		return nil
	}

	err = RunFeed(feed, poster)
	assert.Error(err)

	// verify the event was posted
	assert.Equal(2, hits)
}
