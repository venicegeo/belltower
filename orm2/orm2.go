package orm2

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	elastic "gopkg.in/olivere/elastic.v5"
)

//---------------------------------------------------------------------

type Orm struct {
	esClient *elastic.Client
	ctx      context.Context
}

func NewOrm() (*Orm, error) {

	ctx := context.Background()

	// defaults to 127.0.0.1:9200
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		return nil, err
	}
	log.Printf("Elasticsearch version %s", esversion)

	orm := &Orm{
		esClient: client,
		ctx:      ctx,
	}

	//orm.listAll(true)

	clean := true
	err = orm.register(&Feed{}, clean)
	if err != nil {
		return nil, err
	}

	return orm, nil
}

func (orm *Orm) Create(obj Elasticable) (string, error) {

	if obj.GetID() != "" {
		return "", fmt.Errorf("ID already assigned prior to Create()")
	}
	obj.SetID()

	resp, err := orm.esClient.Index().
		Index(obj.GetIndex()).
		Type(obj.GetType()).
		Id(obj.GetID()).
		BodyJson(obj).
		Do(orm.ctx)
	if err != nil {
		return "", err
	}
	if !resp.Created {
		return "", fmt.Errorf("Create() did not create")
	}

	return resp.Id, nil
}

func (orm *Orm) Read(obj Elasticable) (Elasticable, error) {

	result, err := orm.esClient.Get().
		Index(obj.GetIndex()).
		Type(obj.GetType()).
		Id(obj.GetID()).
		Do(orm.ctx)
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, nil
	}

	src := result.Source

	err = json.Unmarshal(*src, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (orm *Orm) listAll(delete bool) {
	names, err := orm.esClient.IndexNames()
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		fmt.Printf("%s\n", name)
		if delete && strings.HasPrefix(name, "f") {
			resp, err := orm.esClient.DeleteIndex(name).Do(orm.ctx)
			if err != nil {
				panic(err)
			}
			if !resp.Acknowledged {
				panic("DeleteIndex() not acknowledged")
			}
		}
	}
}

func (orm *Orm) register(e Elasticable, clean bool) error {
	index := e.GetIndex()

	exists, err := orm.esClient.IndexExists(index).Do(orm.ctx)
	if err != nil {
		return err
	}

	if clean && exists {
		response, err := orm.esClient.DeleteIndex(index).Do(orm.ctx)
		if err != nil {
			return err
		}
		if !response.Acknowledged {
			return fmt.Errorf("CreateIndex() not acknowledged")

		}
	}

	if !exists {
		result, err := orm.esClient.CreateIndex(index).BodyString(e.GetMapping()).Do(orm.ctx)
		if err != nil {
			return err
		}

		if !result.Acknowledged {
			return fmt.Errorf("CreateIndex() not acknowledged")
		}
	}

	return nil
}

//---------------------------------------------------------------------

var globalID int = 0

func NewID() string {
	globalID++
	return strconv.Itoa(globalID)
}

//---------------------------------------------------------------------

// every object type wil be stored in its own type in its own index
type Elasticable interface {
	GetIndex() string
	GetType() string
	GetMapping() string
	GetID() string
	SetID()
}

type CoreX struct {
	A int
	B float32
}
type Core struct {
	A int
	B float32
	C CoreX
}

type Feed struct {
	Id       string
	Name     string
	Time     time.Time
	Bool     bool
	Int      int
	Float    float64
	IntArray []int
	Object   interface{}
	Core     Core
	Nested   []CoreX
}

func (f *Feed) GetIndex() string {
	return "feed_index"
}

func (f *Feed) GetType() string {
	return "feed_type"
}

func (f *Feed) GetMapping() string {

	mapping := `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"feed_type":{
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
					"type":"bool"
				},
				"int":{
					"type":"integer"
				},
				"float":{
					"type":"double"
				},
				"intArray":{
					"type":"integer"
				},
				"object":{
					"type":"object"
				},
				"core":{
					"type":"object"
				},
				"nested":{
					"type":"nested"
				}
			}
		}
	}
}`

	return mapping
}

func (f *Feed) GetID() string {
	return f.Id
}

func (f *Feed) SetID() {
	f.Id = NewID()
}
