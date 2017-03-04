package drivers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// RandFeed checks to see if a random number from the range [0..limit) is equal to zero. If so,
// it will give the server a message. It sleeps for the given number of seconds between
// checks.
type RandomFeed struct {
	name  string
	sleep time.Duration
	limit int
}

func NewRandomFeed(config *Config) (*RandomFeed, error) {
	f := &RandomFeed{}

	err := f.init(config)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *RandomFeed) init(config *Config) error {

	name, ok := (*config)["name"]
	if !ok {
		return fmt.Errorf("Missing config field: name")
	}
	f.name = name

	sleepStr, ok := (*config)["sleep"]
	if !ok {
		return fmt.Errorf("Missing config field: sleep")
	}
	sleepSecs, err := strconv.Atoi(sleepStr)
	if err != nil {
		return err
	}
	sleep := int64(time.Second) * int64(sleepSecs)
	f.sleep = time.Duration(sleep)

	limitStr, ok := (*config)["limit"]
	if !ok {
		return fmt.Errorf("Missing config field: limit")
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return err
	}
	f.limit = limit

	rand.Seed(17)

	return nil
}

func (f *RandomFeed) Run(statusF StatusF, mssgF MssgF) error {

	var err error
	ok := true

	go func() {
		for {
			ok, err = statusF("good")
			if err != nil {
				return
			}
			if !ok {
				return
			}

			x := rand.Intn(f.limit)
			if x == 0 {
				m := map[string]string{
					"mssg":  time.Now().Format(time.RFC3339),
					"limit": strconv.Itoa(x),
				}

				err = mssgF(m)
				if err != nil {
					return
				}
			}

			time.Sleep(f.sleep)
		}
	}()

	return nil
}
