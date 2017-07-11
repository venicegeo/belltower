package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		text     string
		expected int
	}{
		// constants
		{expected: 17, text: "17"},

		// binops
		{expected: 3, text: "1+2"},
		{expected: -1, text: "1-2"},
		{expected: 6, text: "2*3"},
		{expected: 2, text: "6/3"},
		{expected: 1, text: "6/5"},
	}

	for _, item := range data {
		e, err := NewExpression(item.text)
		assert.NoError(err)
		v, err := e.Eval()
		assert.NoError(err)
		assert.Equal(item.expected, v)
	}
}
