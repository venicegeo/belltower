package dsl

import "testing"

//---------------------------------------------------------------------------

func TestNodes(t *testing.T) {
	/*
		assert := assert.New(t)

		t1 := NewExprNodeArrayRef("arrayo", NewExprNodeIntConstant(17))
		assert.Equal("ARRAY(17, INT)", t1.String())

		t2 := NewExprNodeStructRef("structo", "aa")
		t2s := t2.String()
		assert.Contains(t2s, "STRUCT(aa)")
		assert.Contains(t2s, ")")

		t3 := NewExprNodeMapRef("mappo", NewExprNodeBoolConstant(true))

		assert.Equal("MAP[STRING]BOOL", t3.String())
	*/
}

func TestNodeEquality(t *testing.T) {
	/*
		assert := assert.New(t)

		//eq := func(a, b Node) bool {
		//	return a == b
		//}

		a := NewExprNodeMapRef("mappy", NewExprNodeStringConstant("asdf"))
		b := NewExprNodeMapRef("mappy", NewExprNodeStringConstant("asdf"))
		c := NewExprNodeStringConstant("qwerty")
		d := NewExprNodeFloatConstant(2.1)
		assert.Equal(a, b)
		assert.NotEqual(a, c)
		assert.NotEqual(a, d)
		assert.NotEqual(c, d)
		assert.EqualValues(a, b)
	*/
}
