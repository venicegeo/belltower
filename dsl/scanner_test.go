package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//--------------------------

func TestScannerPrefixMatch(t *testing.T) {
	assert := assert.New(t)

	var siz int
	var ok bool

	ok, _ = matchArrayTypePrefix("[]")
	assert.False(ok)
	ok, _ = matchArrayTypePrefix(" []")
	assert.False(ok)
	ok, _ = matchArrayTypePrefix("[}")
	assert.False(ok)
	ok, _ = matchArrayTypePrefix("[3a]")
	assert.False(ok)
	ok, siz = matchArrayTypePrefix("[3]")
	assert.True(ok)
	assert.Equal(3, siz)
	ok, siz = matchArrayTypePrefix("[32]")
	assert.True(ok)
	assert.Equal(32, siz)
}

func TestScanner(t *testing.T) {
	assert := assert.New(t)

	s := &Scanner{}

	type data struct {
		source string
		tokens []string
	}

	table := []data{
		{
			source: "  as > 10 | b +3 < c",
			tokens: []string{"as", ">", "10", "|", "b", "+", "3", "<", "c"},
		},
		{
			source: `[map] int`,
			tokens: []string{"[map]", "int"},
		},
	}

	for _, testcase := range table {
		tokens, err := s.Scan(testcase.source)
		assert.NoError(err)
		assert.NotNil(tokens)

		assert.Len(tokens, len(testcase.tokens))
		for i := range testcase.tokens {
			assert.Equal(testcase.tokens[i], tokens[i].Text)
		}
	}
}
