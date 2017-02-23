package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//---------------------------------------------------------------------------

func TestExprParser(t *testing.T) {
	assert := assert.New(t)

	for _, item := range exprTestData {
		ep := &ExprParser{}
		node, err := ep.Parse(item.token)
		assert.NoError(err)
		assert.NotNil(node)

		//log.Printf("%v", tc.node)
		//log.Printf("%v", node)
		assert.Equal(item.node, node)
	}
}
