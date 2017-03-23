package orm2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mapsEqual(assert *assert.Assertions, a map[string]interface{}, b map[string]interface{}) {
	assert.Equal(len(a), len(b))

	for k, v := range a {
		vv, ok := b[k]
		assert.True(ok)
		assert.EqualValues(v, vv)
	}
}

func TestIndexOperations(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	e := &Demo{}

	// to set things up correctly (ignore errors)
	orm.DeleteIndex(e)

	err = orm.CreateIndex(e)
	assert.NoError(err)

	exists, err := orm.IndexExists(e)
	assert.NoError(err)
	assert.True(exists)

	// not allowed to create an index that already exists
	err = orm.CreateIndex(e)
	assert.Error(err)

	err = orm.DeleteIndex(e)
	assert.NoError(err)

	exists, err = orm.IndexExists(e)
	assert.NoError(err)
	assert.False(exists)

	// not allowed to delete an index that doesn't exist
	err = orm.DeleteIndex(e)
	assert.Error(err)
}

func TestDocumentCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	orig := &Demo{Name: "Alice"}

	err = orm.CreateIndex(orig)
	assert.NoError(err)

	// does create work?
	id, err := orm.CreateDocument(orig)
	assert.NoError(err)
	assert.NotEmpty(id)

	// does read work?
	tmp := &Demo{Id: id}
	dup, err := orm.ReadDocument(tmp)
	assert.NoError(err)
	assert.NotNil(dup)
	assert.EqualValues(id, dup.GetId())
	assert.EqualValues(orig.Name, dup.(*Demo).Name)

	// update it
	tmp = &Demo{Id: id, Name: "Bob"}
	err = orm.UpdateDocument(tmp)
	assert.NoError(err)

	// read again, to check
	tmp = &Demo{Id: id}
	dup, err = orm.ReadDocument(tmp)
	assert.NoError(err)
	assert.NotNil(dup)
	assert.EqualValues(id, dup.GetId())
	assert.EqualValues("Bob", dup.(*Demo).Name)

	// not allowed to update for invalid id
	tmp = &Demo{Id: "3241234124", Name: "Bob"}
	err = orm.UpdateDocument(tmp)
	assert.Error(err)

	// not allowed to read from an invalid id
	tmp = &Demo{Id: "99999"}
	_, err = orm.ReadDocument(tmp)
	assert.Error(err)

	// try delete
	tmp = &Demo{Id: id}
	err = orm.DeleteDocument(tmp)
	assert.NoError(err)

	// not allowed to delete if doesn't exist
	err = orm.DeleteDocument(tmp)
	assert.Error(err)
}

func TestDemoMappings(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	// to set things up correctly (ignore errors)
	orm.DeleteIndex(&Demo{})

	err = orm.CreateIndex(&Demo{})
	assert.NoError(err)

	now := time.Now()

	feed := &Demo{
		Name:     "Bob",
		Time:     now,
		Bool:     true,
		Int:      17,
		Float:    17.19,
		IntArray: []int{2, 4, 8},
		Object:   map[string]interface{}{"x5": 5, "x10": false},
		Core: DemoCore{
			A2: 5,
			B2: 17.81,
			C2: DemoCoreX{A1: 55, B1: 81.71},
		},
		Nested: []DemoCoreX{
			DemoCoreX{A1: 22, B1: 2.2},
			DemoCoreX{A1: 33, B1: 3.3},
		},
	}

	id, err := orm.CreateDocument(feed)
	assert.NoError(err)
	assert.NotEmpty(id)

	f := &Demo{Id: id}
	g, err := orm.ReadDocument(f)
	assert.NoError(err)
	assert.NotNil(g)
	//log.Printf("%#v", g)
	assert.EqualValues(id, g.GetId())

	assert.EqualValues("Bob", feed.Name)
	assert.EqualValues(feed.Name, g.(*Demo).Name)
	assert.EqualValues(feed.Bool, g.(*Demo).Bool)
	assert.EqualValues(feed.Time, g.(*Demo).Time)
	assert.EqualValues(feed.Int, g.(*Demo).Int)
	assert.EqualValues(feed.Float, g.(*Demo).Float)
	assert.EqualValues(feed.IntArray, g.(*Demo).IntArray)
	mapsEqual(assert, feed.Object.(map[string]interface{}), g.(*Demo).Object.(map[string]interface{}))
	assert.EqualValues(feed.Core, g.(*Demo).Core)
	assert.EqualValues(feed.Nested, g.(*Demo).Nested)
	assert.EqualValues(feed.Nested[1].B1, g.(*Demo).Nested[1].B1)
}

//---------------------------------------------------------------------

type DemoCoreX struct {
	A1 int     `json:"a1"`
	B1 float32 `json:"b1"`
}
type DemoCore struct {
	A2 int       `json:"a2"`
	B2 float32   `json:"b2"`
	C2 DemoCoreX `json:"c2"`
}

type Demo struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Time     time.Time   `json:"time"`
	Bool     bool        `json:"bool"`
	Int      int         `json:"int"`
	Float    float64     `json:"float"`
	IntArray []int       `json:"int_array"`
	Object   interface{} `json:"object"`
	Core     DemoCore    `json:"core"`
	CoreX    DemoCoreX   `json:"corex"`
	Nested   []DemoCoreX `json:"nested"`
}

func (f *Demo) GetIndexName() string {
	return "demo_index"
}

func (f *Demo) GetTypeName() string {
	return "demo_type"
}

func (f *Demo) GetMapping() string {

	mapping := `{
	"settings":{
	},
	"mappings":{
		"demo_type":{
			"dynamic":"strict",
			"properties":{
				"id":{
					"type":"string"
				},
				"name":{
					"type":"string"
				},
				"time":{
					"type":"date"
				},
				"bool":{
					"type":"boolean"
				},
				"int":{
					"type":"integer"
				},
				"float":{
					"type":"double"
				},
				"int_array":{
					"type":"integer"
				},
				"object":{
					"dynamic":"true",
					"type":"object"
				},
				"core":{
					"type":"object",
					"properties":{
						"a2":{ "type":"integer" },
						"b2":{ "type":"double" },
						"c2":{
							"type":"object",
							"properties":{
								"a1":{ "type":"integer" },
								"b1":{ "type":"double" }
							}
						}
					}
				},
				"corex":{
					"type":"object",
					"properties":{
						"a1":{ "type":"integer" },
						"b1":{ "type":"double" }
					}
				},
				"nested":{
					"type":"nested",
					"properties":{
						"a1":{ "type":"integer" },
						"b1":{ "type":"double" }
					}
				}
			}
		}
	}
}`

	return mapping
}

func (f *Demo) GetId() string {
	return f.Id
}

func (f *Demo) SetId() string {
	f.Id = NewId()
	return f.Id
}
