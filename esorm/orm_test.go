package esorm

import (
	"testing"
	"time"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
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

func TestMappingGeneration(t *testing.T) {
	assert := assert.New(t)

	xyzMapping := `{
		"settings":{
		},
		"mappings":{
			"xyz_type":{
				"dynamic":"strict",
				"properties":{
					"id":{
						"type":"string"
					},
					"name":{
						"type":"string"
					}
				}
			}
		}
	}`

	demoMapping := `{
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

	type Data struct {
		obj      Elasticable
		expected string
	}
	data := []Data{
		Data{&Xyz{}, xyzMapping},
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

func TestThingCRUD(t *testing.T) {
	assert := assert.New(t)

	orm, err := NewOrm()
	assert.NoError(err)
	assert.NotNil(orm)

	e := &Xyz{}

	exists, err := orm.IndexExists(e)
	assert.NoError(err)
	if exists {
		err = orm.DeleteIndex(e)
		assert.NoError(err)
	}

	err = orm.CreateIndex(e)
	assert.NoError(err)

	// does create work?
	c := &XyzCreateFields{Name: "one"}
	id, err := orm.CreateThing("u", e, c)
	assert.NoError(err)
	assert.NotEmpty(id)

	tmp := &Xyz{Id: id}

	// does read work?
	r, err := orm.ReadThing(tmp)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.(*XyzReadFields).Id)
	assert.EqualValues("one", r.(*XyzReadFields).Name)

	// update it
	u := &XyzUpdateFields{Name: "two"}
	err = orm.UpdateThing(tmp, u)
	assert.NoError(err)

	// read again, to check
	r, err = orm.ReadThing(tmp)
	assert.NoError(err)
	assert.NotNil(r)
	assert.EqualValues(id, r.(*XyzReadFields).Id)
	assert.EqualValues("two", r.(*XyzReadFields).Name)

	// try delete
	err = orm.DeleteThing(tmp)
	assert.NoError(err)

	// read again, to make sure got deleted
	_, err = orm.ReadThing(tmp)
	assert.Error(err)
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
		"id":        MappingPropertyFields{Type: "string"},
		"name":      MappingPropertyFields{Type: "string"},
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

//---------------------------------------------------------------------

type Xyz struct {
	Id   common.Ident `json:"id"`
	Name string       `json:"name"`
}

type XyzCreateFields struct {
	Name string
}

type XyzReadFields struct {
	Id   common.Ident
	Name string
}

type XyzUpdateFields struct {
	Name string
}

func (xyz *Xyz) SetFieldsForCreate(ownerId common.Ident, fields interface{}) error {
	xyz.Name = fields.(*XyzCreateFields).Name
	return nil
}

func (xyz *Xyz) GetFieldsForRead() (interface{}, error) {
	fields := &XyzReadFields{}
	fields.Id = xyz.Id
	fields.Name = xyz.Name
	return fields, nil
}

func (xyz *Xyz) SetFieldsForUpdate(fields interface{}) error {
	xyz.Name = fields.(*XyzUpdateFields).Name
	return nil
}

func (F *Xyz) String() string { panic(1) }

func (f *Xyz) GetLoweredName() string { return "xyz" }
func (f *Xyz) GetIndexName() string   { return "xyz_index" }
func (f *Xyz) GetTypeName() string    { return "xyz_type" }

func (f *Xyz) GetMappingProperties() map[string]MappingPropertyFields {

	data := map[string]MappingPropertyFields{
		"id":   MappingPropertyFields{Type: "string"},
		"name": MappingPropertyFields{Type: "string"},
	}

	return data
}

func (f *Xyz) GetId() common.Ident {
	return f.Id
}

func (f *Xyz) SetId() common.Ident {
	f.Id = common.NewId()
	return f.Id
}
