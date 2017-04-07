package esorm

import (
	"testing"
	"time"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

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

func TestMappingGeneration(t *testing.T) {
	assert := assert.New(t)

	demoMapping := `{
		"settings":{
		},
		"mappings":{
			"demo_type":{
				"dynamic":"strict",
				"properties":{
					"id":{
						"type":"keyword"
					},
					"name":{
						"type":"keyword"
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

	type Data struct {
		obj      Elasticable
		expected string
	}
	data := []Data{
		Data{&Demo{}, demoMapping},
	}

	for _, d := range data {

		m := NewMapping(d.obj)
		assert.NotNil(m)
		byts, err := json.MarshalIndent(m, "", "    ")
		assert.NoError(err)
		actualMapping := string(byts)
		//log.Printf("%s", actualMapping)
		assert.JSONEq(d.expected, actualMapping)
	}
}

func TestDocumentCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	orig := &Demo{Name: "Alice"}
	orig2 := &Demo{Name: "Zed"}

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

	// make a second one
	id2, err := orm.CreateDocument(orig2)
	assert.NoError(err)
	assert.NotEmpty(id2)

	// does read still work?
	tmp = &Demo{Id: id2}
	dup, err = orm.ReadDocument(tmp)
	assert.NoError(err)
	assert.NotNil(dup)
	assert.EqualValues(id2, dup.GetId())
	assert.EqualValues(orig2.Name, dup.(*Demo).Name)

	//	time.Sleep(5 * time.Second)

	// read all
	makearay := func() []Elasticable {
		ary := make([]Elasticable, 10)
		for i := range ary {
			ary[i] = &Demo{}
		}
		return ary
	}
	{
		// quick side trip to test pagination
		ary0 := makearay()
		ary1, tot, err := orm.ReadDocuments(ary0, 0, 1)
		assert.NoError(err)
		assert.NotNil(ary1)
		assert.Equal(int64(2), tot)
		assert.Len(ary1, 1)

		ary2 := makearay()
		ary3, tot, err := orm.ReadDocuments(ary2, 1, 1)
		assert.NoError(err)
		assert.NotNil(ary3)
		assert.Equal(int64(2), tot)
		assert.Len(ary3, 1)

		ary4 := makearay()
		ary5, tot, err := orm.ReadDocuments(ary4, 0, 10)
		assert.NoError(err)
		assert.NotNil(ary5)
		assert.Equal(int64(2), tot)
		assert.Len(ary5, 2)

		ary6 := makearay()
		ary7, tot, err := orm.ReadDocuments(ary6, 10, 10)
		assert.NoError(err)
		assert.NotNil(ary7)
		assert.Equal(int64(2), tot)
		assert.Len(ary7, 0)
	}
	ary8 := makearay()
	ary9, tot, err := orm.ReadDocuments(ary8, 0, 10)
	assert.NoError(err)
	assert.NotNil(ary9)
	assert.Equal(int64(2), tot)
	assert.Len(ary9, 2)
	/*ary2 := make([]*Demo, len(ary))
	for i, v := range ary {
		tmp := &Demo{}
		err = json.Unmarshal(v, tmp)
		ary2[i] = tmp
		assert.NoError(err)
	}*/

	ok1 := (id == ary9[0].GetId() /*&& "Bob" == ary2[0].Name*/ && id2 == ary9[1].GetId() /*&& "Zed" == ary2[1].Name*/)
	ok2 := (id2 == ary9[0].GetId() /*&& "Zed" == ary2[0].Name*/ && id == ary9[1].GetId() /*&& "Bob" == ary2[1].Name*/)
	assert.True(ok1 || ok2)

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
	assert.True(common.MapsAreEqualValues(feed.Object.(map[string]interface{}), g.(*Demo).Object.(map[string]interface{})))
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
	Id       common.Ident `json:"id"`
	Name     string       `json:"name"`
	Time     time.Time    `json:"time"`
	Bool     bool         `json:"bool"`
	Int      int          `json:"int"`
	Float    float64      `json:"float"`
	IntArray []int        `json:"int_array"`
	Object   interface{}  `json:"object"`
	Core     DemoCore     `json:"core"`
	CoreX    DemoCoreX    `json:"corex"`
	Nested   []DemoCoreX  `json:"nested"`
}

func (f *Demo) SetFieldsForCreate(ownerId common.Ident, fields interface{}) error { panic(1) }
func (f *Demo) GetFieldsForRead() (interface{}, error)                            { panic(1) }
func (f *Demo) SetFieldsForUpdate(fields interface{}) error                       { panic(1) }
func (F *Demo) String() string                                                    { panic(1) }

func (f *Demo) GetLoweredName() string { return "demo" }
func (f *Demo) GetIndexName() string   { return "demo_index" }
func (f *Demo) GetTypeName() string    { return "demo_type" }

func (f *Demo) GetMappingProperties() map[string]MappingPropertyFields {

	data := map[string]MappingPropertyFields{
		"id":        MappingPropertyFields{Type: "keyword"},
		"name":      MappingPropertyFields{Type: "keyword"},
		"time":      MappingPropertyFields{Type: "date"},
		"bool":      MappingPropertyFields{Type: "boolean"},
		"int":       MappingPropertyFields{Type: "integer"},
		"float":     MappingPropertyFields{Type: "double"},
		"int_array": MappingPropertyFields{Type: "integer"},
		"object":    MappingPropertyFields{Type: "object", Dynamic: "true"},

		"core": MappingPropertyFields{
			Type: "object",
			Properties: map[string]MappingPropertyFields{
				"a2": MappingPropertyFields{Type: "integer"},
				"b2": MappingPropertyFields{Type: "double"},
				"c2": MappingPropertyFields{
					Type: "object",
					Properties: map[string]MappingPropertyFields{
						"a1": MappingPropertyFields{Type: "integer"},
						"b1": MappingPropertyFields{Type: "double"},
					},
				},
			},
		},

		"corex": MappingPropertyFields{
			Type: "object",
			Properties: map[string]MappingPropertyFields{
				"a1": MappingPropertyFields{Type: "integer"},
				"b1": MappingPropertyFields{Type: "double"},
			},
		},

		"nested": MappingPropertyFields{
			Type: "nested",
			Properties: map[string]MappingPropertyFields{
				"a1": MappingPropertyFields{Type: "integer"},
				"b1": MappingPropertyFields{Type: "double"},
			},
		},
	}

	return data
}

func (f *Demo) GetId() common.Ident {
	return f.Id
}

func (f *Demo) SetId() common.Ident {
	f.Id = common.NewId()
	return f.Id
}
