package esorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrudFields(t *testing.T) {
	assert := assert.New(t)

	type T struct {
		F1 int `crud:"cru"`
		F2 int `crud:"cr"`
		F3 int `crud:""`
		F4 int `bar:"foo" crud:"" foo:"bar"`
		F5 int ``
		F6 int
	}

	type testdata struct {
		obj      *T
		field    string
		mode     CrudFieldMode
		resultOk bool
		errorOk  bool
	}

	//nosuchstruct := 9

	data := []testdata{
		// negative tests

		// no field name
		{&T{}, "", CrudFieldCreate, false, false},
		// bad field name
		{&T{}, "nosuchfield", CrudFieldCreate, false, false},
		// bad type
		{nil, "F1", CrudFieldCreate, false, false},

		// TODO
		////// bad type
		//////{nosuchstruct, "F1", CrudFieldCreate, false, false},

		// bad mode
		{&T{}, "F1", "nosuchmode", false, false},
		// no mode
		{&T{}, "F1", "", false, false},

		// positive tests

		// passes
		{&T{}, "F1", CrudFieldCreate, true, true},

		// passes
		{&T{}, "F2", CrudFieldCreate, true, true},

		// fails: no such field
		{&T{}, "F3", CrudFieldCreate, false, true},

		// fails: no such field
		{&T{}, "F4", CrudFieldCreate, false, true},

		// fails: no such field
		{&T{}, "F5", CrudFieldCreate, false, true},

		// fails: no such field
		{&T{}, "F6", CrudFieldCreate, false, true},
	}

	for _, d := range data {
		ok, err := IsCrudField(d.obj, d.field, d.mode)
		if d.errorOk {
			assert.NoError(err)
		} else {
			assert.Error(err)
		}
		if d.resultOk {
			assert.True(ok)
		} else {
			assert.False(ok)
		}
	}
}

func TestCrudMerge(t *testing.T) {
	assert := assert.New(t)

	type T struct {
		F1 int    `crud:"cr"`
		F2 int    `crud:"u"`
		F3 int    `crud:""`
		F4 []bool `crud:"u"`
	}

	src := T{1, 2, 3, []bool{false, true, false}}
	dest := T{4, 5, 6, []bool{true, true}}

	type testdata struct {
		mode CrudFieldMode
		f1   int
		f2   int
		f3   int
		f4   []bool
	}
	data := []testdata{
		{CrudFieldCreate, 1, 5, 6, []bool{true, true}},
		{CrudFieldUpdate, 4, 2, 6, []bool{false, true, false}},
		{CrudFieldRead, 1, 5, 6, []bool{true, true}},
	}

	for _, e := range data {
		s := src
		d := dest

		err := CrudMerge(&s, &d, e.mode)
		assert.NoError(err)

		assert.Equal(e.f1, d.F1)
		assert.Equal(e.f2, d.F2)
		assert.Equal(e.f3, d.F3)
		assert.Equal(e.f4, d.F4)
	}
}
