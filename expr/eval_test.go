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
		{expected: 0, text: "0"},
		{expected: 6, text: "1+2+3"},
	}

	for _, item := range data {
		e, err := NewExpression(item.text)
		assert.NoError(err)
		v, err := e.Eval()
		assert.NoError(err)
		assert.Equal(item.expected, v)
	}
}
