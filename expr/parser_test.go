package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFails(t *testing.T) {
	assert := assert.New(t)

	texts := []string{
		"1+",
	}

	for _, text := range texts {
		e, err := NewExpression(text)
		assert.Error(err, text)
		assert.Nil(e, text)
	}
}

func TestEvalFails(t *testing.T) {
	assert := assert.New(t)

	texts := []string{
		"1+x",
	}

	for _, text := range texts {
		e, err := NewExpression(text)
		assert.NoError(err, text)
		assert.NotNil(e, text)

		x, err := e.Eval()
		assert.Error(err, text)
		assert.Nil(x, text)
	}
}

func TestEvals(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		text     string
		expected interface{}
	}{
		// misc
		{expected: 3.0, text: "1+2"},

		// constants
		{expected: 17.0, text: "17"},

		// binops
		{expected: 3.0, text: "1+2"},
		{expected: -1.0, text: "1-2"},
		{expected: 6.0, text: "2*3"},
		{expected: 2.0, text: "6/3"},
		{expected: 1.2, text: "6/5"},
		{expected: true, text: "41>5 && (4/5 < 1.1)"},
	}

	for _, item := range data {
		e, err := NewExpression(item.text)
		assert.NoError(err)
		assert.NotNil(e)

		x, err := e.Eval()
		assert.NoError(err)
		assert.NotNil(x)

		assert.Equal(item.expected, x)
	}
}

func TestConversion(t *testing.T) {
	assert := assert.New(t)

	{
		e, err := NewExpression("1.23")
		assert.NoError(err)
		x, err := e.Eval()
		assert.NoError(err)
		f := AsFloat(x)
		assert.NotNil(f)
		assert.Equal(1.23, *f)
	}
	{
		e, err := NewExpression(`"1,2,3"`)
		assert.NoError(err)
		x, err := e.Eval()
		assert.NoError(err)
		s := AsString(x)
		assert.NotNil(s)
		assert.Equal("1,2,3", *s)
	}
	{
		e, err := NewExpression("true")
		assert.NoError(err)
		x, err := e.Eval()
		assert.NoError(err)
		b := AsBool(x)
		assert.NotNil(b)
		assert.Equal(true, *b)
	}
}
