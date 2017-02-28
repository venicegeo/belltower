package drivers

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type FileSysFeed struct {
	name          string
	sleepDuration time.Duration
	path          string
	files         map[string]bool
}

func (f *FileSysFeed) init(config *Config) error {
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

	path, ok := (*config)["path"]
	if !ok {
		return fmt.Errorf("Missing config field: path")
	}
	f.path = path

	f.files = map[string]bool{}
	currFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fileInfo := range currFiles {
		name := fileInfo.Name()
		f.files[name] = false
	}

	return nil
}

func (f *FileSysFeed) Run(config Config, statusF StatusF, mssgF MssgF) error {

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

			added, err := f.checkFileSys(f.path)
			if err != nil {
				return
			}
			if len(added) > 0 {
				s := strings.Join(added, " ")
				m := map[string]string{
					"mssg":  "Hi!",
					"added": s,
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

func (f *FileSysFeed) checkFileSys(path string) ([]string, error) {
	for k, _ := range f.files {
		f.files[k] = false
	}

	currFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	added := []string{}

	for _, fileInfo := range currFiles {
		name := fileInfo.Name()
		v, ok := f.files[name]
		switch {
		case !ok:
			// not in list from last time, so we have a new file
			f.files[name] = true
			added = append(added, name)
		case v:
			// is in list from last time, and its flag is true: should never happen
			panic(name)
		case !v:
			// is in list from last time, and its flag is false: no change
			f.files[name] = true
		}
	}

	// if any files left from last time with flag still false, remove their entries
	delList := []string{}
	for k, v := range f.files {
		if !v {
			delList = append(delList, k)
		}
	}
	for _, name := range delList {
		delete(f.files, name)
	}

	return added, nil
}
