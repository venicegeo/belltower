package feeders

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

func TestRandomFeed(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	settings := map[string]interface{}{
		"limit": 5,
		"sleep": time.Duration(time.Second * 1),
	}
	runner, err := NewRandomFeedRunner(settings)
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
		limit, err := common.GetMapValueAsInt(settings, "limit")
		assert.NoError(err)
		log.Printf("event ==> [0..%d)", limit)

		hitCount++
		return nil
	}

	err = runner.Run(statusF, mssgF)
	assert.NoError(err)

	time.Sleep(8 * time.Second)

	assert.Equal(countLimit, count)
	assert.True(hitCount >= 0 && hitCount < countLimit)
}
