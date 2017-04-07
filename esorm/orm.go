package esorm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/venicegeo/belltower/common"

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

func (orm *Orm) Close() error {
	return nil
}

func (orm *Orm) CreateDocument(obj Elasticable) (common.Ident, error) {

	if obj.GetId() != "" {
		return "", fmt.Errorf("ID already assigned prior to Create()")
	}

	resp, err := orm.esClient.Index().
		Index(GetIndexName(obj)).
		Type(GetTypeName(obj)).
		Id(obj.SetId().String()).
		BodyJson(obj).
		Do(orm.ctx)
	if err != nil {
		return "", err
	}
	if !resp.Created {
		return "", fmt.Errorf("Create() did not create")
	}

	return common.ToIdent(resp.Id), nil
}

func (orm *Orm) ReadDocument(obj Elasticable) (Elasticable, error) {

	result, err := orm.esClient.Get().
		Index(GetIndexName(obj)).
		Type(GetTypeName(obj)).
		Id(obj.GetId().String()).
		Do(orm.ctx)
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, fmt.Errorf("document not found")
	}

	err = json.Unmarshal(*result.Source, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// TODO: the passed-in array should be sufficient, ought not return it too
// TODO: for now, always return sorted by id (ascending)
func (orm *Orm) ReadDocuments(objs []Elasticable, from int, size int) ([]Elasticable, int64, error) {

	if len(objs) < size {
		return nil, 0, fmt.Errorf("array not long enough")
	}

	result, err := orm.esClient.Search().
		Index(GetIndexName(objs[0])).
		Type(GetTypeName(objs[0])).
		Query(elastic.NewMatchAllQuery()).
		From(from).Size(size).
		Sort("id", true).
		Do(orm.ctx)
	if err != nil {
		return nil, 0, err
	}

	if result.Hits.TotalHits <= 0 {
		return nil, 0, nil
	}

	i := 0
	for _, hit := range result.Hits.Hits {
		err = json.Unmarshal(*hit.Source, objs[i])
		if err != nil {
			return nil, 0, err
		}
		i++
	}

	objs = objs[0:i]

	return objs, result.Hits.TotalHits, nil
}

func (orm *Orm) UpdateDocument(obj Elasticable) error {
	_, err := orm.esClient.Update().
		Index(GetIndexName(obj)).
		Type(GetTypeName(obj)).
		Id(obj.GetId().String()).
		Doc(obj).
		Do(orm.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (orm *Orm) DeleteDocument(obj Elasticable) error {
	res, err := orm.esClient.Delete().
		Index(GetIndexName(obj)).
		Type(GetTypeName(obj)).
		Id(obj.GetId().String()).
		Do(orm.ctx)
	if err != nil {
		return err
	}
	if !res.Found {
		return fmt.Errorf("document not found")
	}
	return nil
}

func (orm *Orm) listIndexes(delete bool) {
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
	exists, err := orm.esClient.IndexExists(GetIndexName(e)).Do(orm.ctx)
	return exists, err
}

func (orm *Orm) DeleteIndex(e Elasticable) error {
	response, err := orm.esClient.DeleteIndex(GetIndexName(e)).Do(orm.ctx)
	if err != nil {
		return err
	}
	if !response.Acknowledged {
		return fmt.Errorf("DeleteIndex() not acknowledged")
	}
	return nil
}

func (orm *Orm) CreateIndex(e Elasticable) error {
	index := GetIndexName(e)

	exists, err := orm.IndexExists(e)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index already exists")
	}

	mapping := NewMapping(e)
	byts, err := json.Marshal(mapping)
	//log.Printf("%s", string(byts))
	if err != nil {
		return err
	}
	mappingString := string(byts)

	result, err := orm.esClient.CreateIndex(index).BodyString(mappingString).Do(orm.ctx)
	if err != nil {
		return err
	}

	if !result.Acknowledged {
		return fmt.Errorf("CreateIndex() not acknowledged")
	}

	return nil
}
