package orm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

type Config struct {
	feedType   string
	name       string
	isEnabled  string
	attributes map[string]interface{}
}

type Feed struct {
	Core

	FeedType  string
	Name      string
	IsEnabled bool

	MessageDecl       map[string]interface{}
	MessageDeclAsJSON common.JSON // owned by system

	Config       map[string]interface{}
	ConfigAsJSON common.JSON // owned by system

	NumMessagesRecieved uint
	LastMessageAt       *time.Time
	Owner               User
	OwnerID             uint
}

func (f Feed) String() string {
	s := fmt.Sprintf("f.%d: %s", f.ID, f.Name)
	return s
}
