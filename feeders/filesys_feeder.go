package feeders

import (
	"io/ioutil"
	"strings"

	"github.com/venicegeo/belltower/btorm"
	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

const FileSysFeederId common.Ident = "fce4cb70-11f0-4c70-98ed-5ec0d5273fda"

// implements Feeder
type FileSysFeeder struct {
	settings map[string]string
	path     string
	files    map[string]bool
}

type FileSysEventData struct {
	Added   string
	Deleted string
}

func (f *FileSysFeeder) GetId() common.Ident {
	return FileSysFeederId
}

func (f *FileSysFeeder) GetName() string {
	return "FileSysFeeder"
}

func (f *FileSysFeeder) GetSettingsSchema() map[string]string {
	return map[string]string{
		"Path": "string",
	}
}

func (f *FileSysFeeder) GetEventSchema() map[string]string {
	return map[string]string{
		"Added":   "string",
		"Deleted": "string",
	}
}

//---------------------------------------------------------------------

func init() {
	info := &FeederInfo{
		FeederId:    FileSysFeederId,
		Description: "file system feeder",
		Create:      FileSysFeederCreate,
	}
	feederRegistry.register(info)
}

func FileSysFeederCreate(feed *btorm.Feed) (Feeder, error) {
	path := feed.Settings["Path"]

	f := &FileSysFeeder{
		settings: feed.Settings,
		path:     path,
		files:    nil,
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

func (f *FileSysFeeder) Poll() (interface{}, error) {

	added, deleted, err := f.checkFileSys(f.path)
	if err != nil {
		return nil, err
	}

	e := &FileSysEventData{
		Added:   strings.Join(added, " "),
		Deleted: strings.Join(deleted, " "),
	}

	return e, nil
}

func (f *FileSysFeeder) checkFileSys(path string) (added []string, deleted []string, err error) {
	for k := range f.files {
		f.files[k] = false // not yet visited
	}

	currFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}

	added = []string{}
	deleted = []string{}

	for _, fileInfo := range currFiles {
		name := fileInfo.Name()
		v, ok := f.files[name]
		switch {
		case !ok:
			// not in list from last time, so we have a new file; put on "added"" list, mark as visited
			f.files[name] = true
			added = append(added, name)
		case v:
			// is in list from last time, and its flag is true: shouldn't happen, internal error
			panic(name)
		case !v:
			// is in list from last time, and its flag is false: no change to file, just mark visited
			f.files[name] = true
		}
	}

	// if any files left from last time with flag still false, remove their entries
	for k, v := range f.files {
		if !v {
			deleted = append(deleted, k)
		}
	}

	for _, name := range deleted {
		delete(f.files, name)
	}

	return added, deleted, nil
}
