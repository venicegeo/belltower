package feeders

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/orm"
)

type FileSysFeed struct {
	id    uint
	name  string
	sleep time.Duration
	path  string
	files map[string]bool
}

func NewFileSysFeed(feed *orm.Feed) (*FileSysFeed, error) {
	var _ FeedRunner = &FileSysFeed{}

	f := &FileSysFeed{
		id:   feed.ID,
		name: feed.Name,
	}

	err := f.setVars(feed.Config)
	if err != nil {
		return nil, err
	}

	f.files = map[string]bool{}
	currFiles, err := ioutil.ReadDir(f.path)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range currFiles {
		name := fileInfo.Name()
		f.files[name] = false
	}

	return f, nil
}

func (f *FileSysFeed) setVars(m map[string]interface{}) error {
	var err error

	f.sleep, err = common.GetMapValueAsDuration(m, "sleep")
	if err != nil {
		return err
	}

	f.path, err = common.GetMapValueAsString(m, "path")
	if err != nil {
		return err
	}

	return nil
}

func (rf *FileSysFeed) ID() uint {
	return rf.id
}

func (rf *FileSysFeed) Name() string {
	return rf.name
}

func (f *FileSysFeed) Run(statusF StatusF, mssgF MssgF) error {

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

			time.Sleep(f.sleep)
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
