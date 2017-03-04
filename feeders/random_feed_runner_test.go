package feeders

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/orm"
)

func TestRandomFeed(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	feed := &orm.Feed{
		ID:   17,
		Name: "randomtest",
		Settings: map[string]interface{}{
			"limit": 5,
			"sleep": 1,
		},
	}
	runner, err := NewRandomFeedRunner(feed)
	assert.NoError(err)

	hitCount := 0
	countLimit := 5
	count := 0

	statusF := func(s string) (bool, error) {
		assert.Equal("good", s)
		count++
		if count == countLimit {
			return false, nil
		}
		return true, nil
	}

	mssgF := func(data map[string]string) error {
		//tim, err := time.Parse(time.RFC3339, data["mssg"])
		//if err != nil {
		//	return err
		//}
		limit, err := common.GetMapValueAsInt(feed.Settings, "limit")
		assert.NoError(err)
		log.Printf("event ==> [0..%d)", limit)

		hitCount++
		return nil
	}

	err = runner.Run(statusF, mssgF)
	assert.NoError(err)

	sleep, err := common.GetMapValueAsInt(feed.Settings, "sleep")
	assert.NoError(err)
	dur := time.Duration(float64(sleep*countLimit)*1.25) * time.Second
	time.Sleep(dur)

	assert.Equal(countLimit, count)
	assert.True(hitCount >= 0 && hitCount < countLimit)
}
