package esorm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------------

func TestQuery(t *testing.T) {
	assert := assert.New(t)

	type T struct {
		filter   *Filter
		expected string
	}
	data := []T{
		T{
			filter: Match("message", "foo"),
			expected: `{
				"query": {
					"bool": {
						"filter": {
							"match": {
								"message": "foo"
							}
						}
					}
				}
			}`,
		}, T{
			filter: Or(Match("a", "b"), Match("a", "b")),
			expected: `{
				"query": {
					"bool": {
						"filter": {
							"bool": {
								"should": [
									{
										"match": {
											"a": "b"
										}
									},{
										"match": {
											"a": "b"
										}
									}
								]	
							}
						}
					}
				}
			}`,
		}, T{
			filter: And(Match("c", "d"), Match("c", "d")),
			expected: `{
				"query": {
					"bool": {
						"filter": {
							"bool": {
								"must": [
									{
										"match": {
											"c": "d"
										}
									},{
										"match": {
											"c": "d"
										}
									}
								]	
							}
						}
					}
				}
			}`,
		},
	}

	for _, d := range data {
		assert.JSONEq(d.expected, d.filter.toJson())
	}
}

//---------------------------------------------------------------------------

type Doc struct {
	Id      common.Ident `json:"id"`
	Message string       `json:"message"`
}

func (d *Doc) GetId() common.Ident   { return d.Id }
func (d *Doc) SetId(id common.Ident) { d.Id = id }
func (d *Doc) GetIndexName() string  { return "doc_index" }
func (d *Doc) GetTypeName() string   { return "doc_type" }
func (d *Doc) String() string        { return fmt.Sprintf("%#v", d) }

func (d *Doc) GetMappingProperties() map[string]MappingProperty {
	data := map[string]MappingProperty{
		"id":      MappingProperty{Type: "keyword"},
		"message": MappingProperty{Type: "text"},
	}
	return data
}

func TestPercolation(t *testing.T) {
	assert := assert.New(t)

	var orm Ormer = &Orm{}
	if mock {
		orm = &FakeOrm{}
	}
	err := orm.Open()
	assert.NoError(err)

	_ = orm.DeleteIndex(&Doc{})

	model := &Doc{}

	// create index
	err = orm.CreateIndex(model, true)
	assert.NoError(err)

	filters := []*Filter{
		Match("message", "foo"),                               // 0
		Match("message", "foo"),                               // 1
		Match("message", "bar"),                               // 2
		Match("message", "baz"),                               // 3
		Or(Match("message", "baz"), Match("message", "foo")),  // 4
		And(Match("message", "aaa"), Match("message", "bbb")), // 5
		And(Match("message", "aaa"), Match("message", "aaa")), // 6
	}
	filterIds := make([]common.Ident, len(filters))

	for i, filter := range filters {
		id, err := orm.(*Orm).CreatePercolatorDocument(model, filter)
		assert.NoError(err)
		filterIds[i] = id
	}

	type T struct {
		doc *Doc
		qs  []common.Ident
	}
	data := []T{
		T{
			doc: &Doc{Message: "bar"},
			qs:  []common.Ident{filterIds[2]},
		},
		T{
			doc: &Doc{Message: "ba"},
			qs:  []common.Ident{},
		},
		T{
			doc: &Doc{Message: "foo"},
			qs:  []common.Ident{filterIds[0], filterIds[1], filterIds[4]},
		},
		T{
			doc: &Doc{Message: "baz"},
			qs:  []common.Ident{filterIds[3], filterIds[4]},
		},
		T{
			doc: &Doc{Message: "yow"},
			qs:  []common.Ident{},
		},
		T{
			doc: &Doc{Message: "bbb"},
			qs:  []common.Ident{},
		},
		T{
			doc: &Doc{Message: "aaa"},
			qs:  []common.Ident{filterIds[6]},
		},
	}

	for _, item := range data {
		ary, cnt, err := orm.(*Orm).CreatePercolatorQuery(item.doc)
		assert.NoError(err)
		assert.EqualValues(len(item.qs), cnt)
		assert.Len(ary, len(item.qs))
		assert.EqualValues(item.qs, ary)
	}
}

//---------------------------------------------------------------------
