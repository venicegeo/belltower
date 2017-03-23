package orm2

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	if !strings.HasPrefix(esversion, "5.2") {
		return nil, fmt.Errorf("unsupported elasticsearch version: %s", esversion)
	}

	orm := &Orm{
		esClient: client,
		ctx:      ctx,
	}

	//orm.listAll(true)

	return orm, nil
}

func (orm *Orm) CreateDocument(obj Elasticable) (string, error) {

	if obj.GetID() != "" {
		return "", fmt.Errorf("ID already assigned prior to Create()")
	}

	resp, err := orm.esClient.Index().
		Index(obj.GetIndexName()).
		Type(obj.GetTypeName()).
		Id(obj.SetID()).
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

func (orm *Orm) ReadDocument(obj Elasticable) (Elasticable, error) {

	result, err := orm.esClient.Get().
		Index(obj.GetIndexName()).
		Type(obj.GetTypeName()).
		Id(obj.GetID()).
		Do(orm.ctx)
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, fmt.Errorf("document not found")
	}

	src := result.Source

	//log.Printf("%s", *src)
	err = json.Unmarshal(*src, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (orm *Orm) DeleteDocument(obj Elasticable) error {
	res, err := orm.esClient.Delete().
		Index(obj.GetIndexName()).
		Type(obj.GetTypeName()).
		Id(obj.GetID()).
		Do(orm.ctx)
	if err != nil {
		return err
	}
	if !res.Found {
		return fmt.Errorf("document not found")
	}
	return nil
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

func (orm *Orm) IndexExists(e Elasticable) (bool, error) {
	exists, err := orm.esClient.IndexExists(e.GetIndexName()).Do(orm.ctx)
	return exists, err
}

func (orm *Orm) DeleteIndex(e Elasticable) error {
	response, err := orm.esClient.DeleteIndex(e.GetIndexName()).Do(orm.ctx)
	if err != nil {
		return err
	}
	if !response.Acknowledged {
		return fmt.Errorf("CreateIndex() not acknowledged")
	}
	return nil
}

func (orm *Orm) CreateIndex(e Elasticable) error {
	index := e.GetIndexName()

	exists, err := orm.IndexExists(e)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index already exists")
	}

	err = validateJson(e.GetMapping())
	if err != nil {
		return err
	}

	result, err := orm.esClient.CreateIndex(index).BodyString(e.GetMapping()).Do(orm.ctx)
	if err != nil {
		return err
	}

	if !result.Acknowledged {
		return fmt.Errorf("CreateIndex() not acknowledged")
	}

	return nil
}

func validateJson(s string) error {
	//log.Printf("== %s ==", s)

	obj := &map[string]interface{}{}
	err := json.Unmarshal([]byte(s), obj)
	if err != nil {
		return err
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
	GetIndexName() string
	GetTypeName() string
	GetMapping() string
	GetID() string
	SetID() string
}
