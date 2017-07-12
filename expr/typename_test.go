package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeEnum(t *testing.T) {
	assert := assert.New(t)

	items := []struct {
		src string
		ok  bool
	}{
		{src: "int", ok: true},
		{src: "float", ok: true},
		{src: "bool", ok: true},
		{src: "string", ok: true},
		{src: "map", ok: true},
		{src: "array", ok: true},
		{src: "struct", ok: true},
		{src: "", ok: false},
		{src: "record", ok: false},
	}

	for _, item := range items {
		t := FromString(item.src)
		if item.ok {
			assert.NotEqual(TypeNameInvalid, t)
			s := t.String()
			assert.Equal(item.src, s)
		} else {
			assert.Equal(TypeNameInvalid, t)
		}
	}
}
