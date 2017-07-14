package orm

import (
	"fmt"
	"time"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type Feed struct {
	Core
	FeederId        common.Ident      `json:"feeder_id"        crud:"cr"`
	PollingInterval uint              `json:"polling_interval" crud:"cr"` // in seconds
	PollingEndAt    time.Time         `json:"polling_end_at"   crud:"cr"`
	MessageCount    uint              `json:"message_count"    crud:"r"`
	LastMessageAt   time.Time         `json:"last_message_at"  crud:"r"`
	Settings        map[string]string `json:"settings"         crud:"cr"`
}

func (feed *Feed) GetIndexName() string { return "feed_index" }
func (feed *Feed) GetTypeName() string  { return "feed_type" }

func (feed *Feed) String() string { return fmt.Sprintf("%#v", feed) }

//---------------------------------------------------------------------

func (feed *Feed) GetMappingProperties() map[string]esorm.MappingProperty {
	properties := map[string]esorm.MappingProperty{}

	properties["feeder_id"] = esorm.MappingProperty{Type: "keyword"}
	properties["message_count"] = esorm.MappingProperty{Type: "integer"}
	properties["polling_interval"] = esorm.MappingProperty{Type: "integer"}
	properties["polling_end_at"] = esorm.MappingProperty{Type: "date"}
	properties["last_message_at"] = esorm.MappingProperty{Type: "date"}
	properties["settings"] = esorm.MappingProperty{Type: "object", Dynamic: "true"}

	for k, v := range feed.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------
