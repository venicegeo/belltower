package drivers

import (
	"testing"
	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestClockFeed(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	tf := ClockFeed{}

	statusCount := 5
	hitCount := 0

	mod := 2
	config := Config{
		"mod":   fmt.Sprintf("%d", mod),
		"name":  "clock1",
		"sleep": "1",
	}

	statusF := func(s string) (bool, error) {
		assert.Equal("good", s)
		statusCount--
		if statusCount == 0 {
			return false, nil
		}
		return true, nil
	}

	mssgF := func(data map[string]string) error {
		tim, err := time.Parse(time.RFC3339, data["mssg"])
		if err != nil {
			return err
		}
		assert.Zero(tim.Second() % mod)
		//log.Printf("event==> %s", tim)

		delta := time.Now().Sub(tim).Seconds()
		assert.True(delta < 1.5)

		hitCount++
		return nil
	}

	err := tf.Run(config, statusF, mssgF)
	assert.NoError(err)

	time.Sleep(7 * time.Second)

	assert.Zero(statusCount)
	assert.True(hitCount > 1 && hitCount < 4) // 2 or 3
}
