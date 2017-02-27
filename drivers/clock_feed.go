package drivers

import (
	"fmt"
	"strconv"
	"time"
)

// ClockFeed checks the second hand of the clock and if it is zero when mod'd with the
// given value, it will give the server a message. It sleeps for the given time length between
// checks.
type ClockFeed struct {
	name          string
	sleepDuration time.Duration
	secondsMod    int
}

func (f *ClockFeed) init(config *Config) error {
	name, ok := (*config)["name"]
	if !ok {
		return fmt.Errorf("Missing config field: name")
	}
	f.name = name

	tim, ok := (*config)["sleep"]
	if !ok {
		return fmt.Errorf("Missing config field: sleep")
	}
	secs, err := strconv.Atoi(tim)
	if err != nil {
		return err
	}
	dur := int64(time.Second) * int64(secs)
	f.sleepDuration = time.Duration(dur)

	modstr, ok := (*config)["mod"]
	if !ok {
		return fmt.Errorf("Missing config field: mod")
	}
	mod, err := strconv.Atoi(modstr)
	if err != nil {
		return err
	}
	if mod > 59 || mod < 1 {
		return fmt.Errorf("Illegal config field for mod: %s", modstr)
	}
	f.secondsMod = mod

	return nil
}

func (f *ClockFeed) Run(config Config, statusF StatusF, mssgF MssgF) error {

	err := f.init(&config)
	if err != nil {
		return err
	}

	err = nil
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

			now := time.Now()
			secs := now.Second()
			secs %= f.secondsMod
			if secs == 0 {
				m := map[string]string{
					"mssg": now.Format(time.RFC3339),
				}

				err = mssgF(m)
				if err != nil {
					return
				}
			}

			time.Sleep(f.sleepDuration)
		}
	}()

	return nil
}
