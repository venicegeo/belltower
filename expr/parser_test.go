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
		{expected: false, text: "1+"},

		{expected: true, text: "1+2"},
		{expected: true, text: "a+b"},
		{expected: true, text: "f(8)"},
		{expected: true, text: "f((x))"},
	}

	for _, item := range data {
		e, err := NewExpression(item.text)
		assert.Equal(item.expected, err == nil, item.text)
		assert.Equal(item.expected, e != nil, item.text)
	}
}
