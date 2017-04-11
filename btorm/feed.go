package btorm

import (
	"time"

	"github.com/venicegeo/belltower/common"
	"github.com/venicegeo/belltower/esorm"
)

//---------------------------------------------------------------------

type Feed struct {
	Core
	FeederId        common.Ident           `json:"feeder_id"        crud:"cr"`
	PollingInterval uint                   `json:"polling_interval" crud:"cr"` // in seconds
	PollingEndAt    time.Time              `json:"polling_end_at"   crud:"cr"`
	MessageCount    uint                   `json:"message_count"    crud:"r"`
	LastMessageAt   time.Time              `json:"last_message_at"  crud:"r"`
	Settings        map[string]interface{} `json:"settings"         crud:"cr"`
}

//---------------------------------------------------------------------

func (feed *Feed) GetMappingProperties() map[string]esorm.MappingPropertyFields {
	properties := map[string]esorm.MappingPropertyFields{}

	properties["feeder_id"] = esorm.MappingPropertyFields{Type: "keyword"}
	properties["message_count"] = esorm.MappingPropertyFields{Type: "integer"}
	properties["polling_interval"] = esorm.MappingPropertyFields{Type: "integer"}
	properties["polling_end_at"] = esorm.MappingPropertyFields{Type: "date"}
	properties["last_message_at"] = esorm.MappingPropertyFields{Type: "date"}
	properties["settings"] = esorm.MappingPropertyFields{Type: "object", Dynamic: "true"}

	for k, v := range feed.Core.GetCoreMappingProperties() {
		properties[k] = v
	}

	return properties
}

//---------------------------------------------------------------------
