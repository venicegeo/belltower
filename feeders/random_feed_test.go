package drivers

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomFeed(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	tf := RandomFeed{}

	limit := 5
	hitCount := 0
	sleep := 1
	countLimit := 5
	count := 0

	config := Config{
		"limit": strconv.Itoa(limit),
		"name":  "random1",
		"sleep": strconv.Itoa(sleep),
	}

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
		log.Printf("event ==> [0..%d)", limit)

		hitCount++
		return nil
	}

	err := tf.Run(config, statusF, mssgF)
	assert.NoError(err)

	dur := time.Duration(float64(sleep*countLimit)*1.25) * time.Second
	time.Sleep(dur)

	assert.Equal(countLimit, count)
	assert.True(hitCount >= 0 && hitCount < countLimit)
}
