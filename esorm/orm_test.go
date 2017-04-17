package esorm

import (
	"context"
	"log"
	"testing"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"

	"encoding/json"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/venicegeo/belltower/common"
)

const mock = false

func TestIndexOperations(t *testing.T) {
	assert := assert.New(t)

	var orm Ormer = &Orm{}
	if mock {
		orm = &FakeOrm{}
	}
	err := orm.Open()
	assert.NoError(err)

	e := &Demo{}
	assert.Equal("demo_index", GetIndexName(e))

	// to set things up correctly (ignore errors)
	_ = orm.DeleteIndex(e)

	err = orm.CreateIndex(e)
	assert.NoError(err)

	exists, err := orm.IndexExists(e)
	assert.NoError(err)
	assert.True(exists)

	names, err := orm.GetIndexes()
	assert.NoError(err)
	assert.True(len(names) >= 1)
	assert.Contains(names, GetIndexName(e))

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

	err = orm.Close()
	assert.NoError(err)
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
					"str_map":{
						"dynamic":"true",
						"type":"object"
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

	var orm Ormer = &Orm{}
	if mock {
		orm = &FakeOrm{}
	}
	err := orm.Open()
	assert.NoError(err)

	orig := &Demo{Name: "Alice"}
	orig2 := &Demo{Name: "Zed"}

	err = orm.CreateIndex(orig)
	assert.NoError(err)

	// does create work?
	id, err := orm.CreateDocument(orig)
	assert.NoError(err)
	assert.NotEmpty(id)

	// does read work?
	dup, err := orm.ReadDocument(&Demo{Id: id})
	assert.NoError(err)
	assert.NotNil(dup)
	assert.EqualValues(id, dup.(*Demo).GetId())
	assert.EqualValues(orig.Name, dup.(*Demo).Name)

	// update it
	src := &Demo{Id: id, Name: "Bob"}
	err = orm.UpdateDocument(src)
	assert.NoError(err)

	// read again, to check
	dup, err = orm.ReadDocument(&Demo{Id: id})
	assert.NoError(err)
	assert.NotNil(dup)
	assert.EqualValues(id, dup.(*Demo).GetId())
	assert.EqualValues("Bob", dup.(*Demo).Name)

	// not allowed to update for invalid id
	src = &Demo{Id: "3241234124", Name: "Bob"}
	err = orm.UpdateDocument(src)
	assert.Error(err)

	// not allowed to read from an invalid id
	_, err = orm.ReadDocument(&Demo{Id: "99999"})
	assert.Error(err)

	// make a second one
	id2, err := orm.CreateDocument(orig2)
	assert.NoError(err)
	assert.NotEmpty(id2)

	// does read still work?
	dup, err = orm.ReadDocument(&Demo{Id: id2})
	assert.NoError(err)
	assert.NotNil(dup)
	assert.EqualValues(id2, dup.(*Demo).GetId())
	assert.EqualValues(orig2.Name, dup.(*Demo).Name)

	//	time.Sleep(5 * time.Second)

	// read all
	{
		tmp := &Demo{}

		// quick side trip to test pagination
		ary1, tot, err := orm.ReadDocuments(tmp, 0, 1)
		assert.NoError(err)
		assert.NotNil(ary1)
		assert.Equal(int64(2), tot)
		assert.Len(ary1, 1)

		ary3, tot, err := orm.ReadDocuments(tmp, 1, 1)
		assert.NoError(err)
		assert.NotNil(ary3)
		assert.Equal(int64(2), tot)
		assert.Len(ary3, 1)

		ary5, tot, err := orm.ReadDocuments(tmp, 0, 10)
		assert.NoError(err)
		assert.NotNil(ary5)
		assert.Equal(int64(2), tot)
		assert.Len(ary5, 2)

		ary7, tot, err := orm.ReadDocuments(tmp, 10, 10)
		assert.NoError(err)
		assert.NotNil(ary7)
		assert.Equal(int64(2), tot)
		assert.Len(ary7, 0)
	}

	ary9, tot, err := orm.ReadDocuments(&Demo{}, 0, 10)
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
	tmp := &Demo{Id: id}
	err = orm.DeleteDocument(tmp)
	assert.NoError(err)

	// not allowed to delete if doesn't exist
	err = orm.DeleteDocument(tmp)
	assert.Error(err)
}

func TestDemoMappings(t *testing.T) {
	assert := assert.New(t)

	var orm Ormer = &Orm{}
	if mock {
		orm = &FakeOrm{}
	}
	err := orm.Open()
	assert.NoError(err)

	// to set things up correctly (ignore errors)
	_ = orm.DeleteIndex(&Demo{})

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
		StrMap:   map[string]string{"a": "b", "c": "d", "e": "f"},
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

	g, err := orm.ReadDocument(&Demo{Id: id})
	assert.NoError(err)
	assert.NotNil(g)
	assert.EqualValues(id, g.(*Demo).GetId())

	assert.EqualValues("Bob", feed.Name)
	assert.EqualValues(feed.Name, g.(*Demo).Name)
	assert.EqualValues(feed.Bool, g.(*Demo).Bool)
	assert.EqualValues(feed.Time, g.(*Demo).Time)
	assert.EqualValues(feed.Int, g.(*Demo).Int)
	assert.EqualValues(feed.Float, g.(*Demo).Float)
	assert.EqualValues(feed.IntArray, g.(*Demo).IntArray)
	//log.Printf("%#v   %#v", feed.StrMap, g.(*Demo).StrMap)
	assert.EqualValues(feed.StrMap, g.(*Demo).StrMap)
	//assert.True(common.MapsAreEqualValues(feed.StrMap, g.(*Demo).StrMap))
	assert.True(common.MapsAreEqualValues(feed.Object.(map[string]interface{}), g.(*Demo).Object.(map[string]interface{})))
	assert.EqualValues(feed.Core, g.(*Demo).Core)
	assert.EqualValues(feed.Nested, g.(*Demo).Nested)
	assert.EqualValues(feed.Nested[1].B1, g.(*Demo).Nested[1].B1)
}

//---------------------------------------------------------------------

type Queries struct {
	Id    common.Ident `json:"id"`
	Query string       `json:"query"`
}

type Doc struct {
	Id      common.Ident `json:"id"`
	Message string       `json:"message"`
}

func (d *Doc) GetId() common.Ident {
	return d.Id
}

func (d *Doc) SetId() common.Ident {
	d.Id = common.NewId()
	return d.Id
}

func (d *Doc) String() string { return fmt.Sprintf("%#v", d) }

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

	// to set things up correctly (ignore errors)
	_ = orm.DeleteIndex(&Doc{})

	obj := &Doc{}
	index := GetIndexName(obj)

	index = "docindex"

	res0, err := orm.(*Orm).esClient.DeleteIndex(index).Do(orm.(*Orm).ctx)
	assert.NoError(err)
	assert.True(res0.Acknowledged)

	// create index
	mappingString := `
	{
	   	"settings":{},
	   	"mappings":{
	   		"doctype":{
	   			"properties":{
	   				"message":{
						"type":"text",
						"store": true,
						"fielddata": true
					}
				}
	   		},
	   		"queries":{
	   			"properties":{
	   				"query":{"type":"percolator"}
	   			}
	   		}
	   	}
	}`
	log.Printf(">> %#v", mappingString)
	res1, err := orm.(*Orm).esClient.CreateIndex("docindex").BodyString(mappingString).Do(orm.(*Orm).ctx)
	assert.NoError(err)
	assert.True(res1.Acknowledged)
	log.Printf(">> %#v", res1)

	// Add a document
	res2, err := orm.(*Orm).esClient.Index().
		Index(index).
		Type("queries").
		Id("1").
		BodyJson(`{"query":{"match":{"message":"bonsai tree"}}}`).
		Refresh("wait_for").
		Do(context.TODO())
	assert.NoError(err)
	log.Printf("\n>>>>>> %#v", res2)

	// Percolate should return our registered query
	pq := elastic.NewPercolatorQuery().
		Field("query").
		DocumentType("doctype").
		Document(Doc{Message: "A new bonsai tree in the office"})
	res3, err := orm.(*Orm).esClient.Search(index).Query(pq).Do(orm.(*Orm).ctx)
	assert.NoError(err)
	log.Printf("\n>>>>>>>>>>>>> %#v", res3)
	log.Printf("\n>>>>>>>>>>>>> %#v", res3.Hits.Hits[0])
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
	Id       common.Ident      `json:"id"        crud:"r"`
	Name     string            `json:"name"      crud:"cru"`
	Time     time.Time         `json:"time"      crud:"cru"`
	Bool     bool              `json:"bool"      crud:"r"`
	Int      int               `json:"int"`
	Float    float64           `json:"float"`
	IntArray []int             `json:"int_array"`
	StrMap   map[string]string `json:"str_map"`
	Object   interface{}       `json:"object"`
	Core     DemoCore          `json:"core"`
	CoreX    DemoCoreX         `json:"corex"`
	Nested   []DemoCoreX       `json:"nested"`
}

func (d *Demo) String() string { return fmt.Sprintf("%#v", d) }

func (d *Demo) GetMappingProperties() map[string]MappingProperty {

	data := map[string]MappingProperty{
		"id":        MappingProperty{Type: "keyword"},
		"name":      MappingProperty{Type: "keyword"},
		"time":      MappingProperty{Type: "date"},
		"bool":      MappingProperty{Type: "boolean"},
		"int":       MappingProperty{Type: "integer"},
		"float":     MappingProperty{Type: "double"},
		"int_array": MappingProperty{Type: "integer"},
		"str_map":   MappingProperty{Type: "object", Dynamic: "true"},
		"object":    MappingProperty{Type: "object", Dynamic: "true"},

		"core": MappingProperty{
			Type: "object",
			Properties: map[string]MappingProperty{
				"a2": MappingProperty{Type: "integer"},
				"b2": MappingProperty{Type: "double"},
				"c2": MappingProperty{
					Type: "object",
					Properties: map[string]MappingProperty{
						"a1": MappingProperty{Type: "integer"},
						"b1": MappingProperty{Type: "double"},
					},
				},
			},
		},

		"corex": MappingProperty{
			Type: "object",
			Properties: map[string]MappingProperty{
				"a1": MappingProperty{Type: "integer"},
				"b1": MappingProperty{Type: "double"},
			},
		},

		"nested": MappingProperty{
			Type: "nested",
			Properties: map[string]MappingProperty{
				"a1": MappingProperty{Type: "integer"},
				"b1": MappingProperty{Type: "double"},
			},
		},
	}

	return data
}

func (d *Demo) GetId() common.Ident {
	return d.Id
}

func (d *Demo) SetId() common.Ident {
	d.Id = common.NewId()
	return d.Id
}
