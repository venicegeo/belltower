package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		text     string
		expected bool
	}{
		{expected: true, text: "1+2"},
		{expected: false, text: "1+"},
	}

	for _, item := range data {
		e, err := NewExpression(item.text)
		assert.Equal(item.expected, err == nil)
		assert.Equal(item.expected, e != nil)
	}
}
