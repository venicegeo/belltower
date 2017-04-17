package esorm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/venicegeo/belltower/common"

	"golang.org/x/net/context"

	elastic "gopkg.in/olivere/elastic.v5"
)

//---------------------------------------------------------------------

type Ormer interface {
	Open() error
	Close() error

	CreateDocument(obj Elasticable) (common.Ident, error)
	ReadDocument(typ Elasticable) (Elasticable, error)
	ReadDocuments(typ Elasticable, from int, size int) ([]Elasticable, int64, error)
	UpdateDocument(src Elasticable) error
	DeleteDocument(obj Elasticable) error

	GetIndexes() ([]string, error)
	IndexExists(e Elasticable) (bool, error)
	DeleteIndex(e Elasticable) error
	CreateIndex(e Elasticable, usePercolation bool) error
}

type Orm struct {
	esClient *elastic.Client
	ctx      context.Context
}

func init() {
	var _ Ormer = (*Orm)(nil)
}

func (orm *Orm) Open() error {

	ctx := context.Background()

	// defaults to 127.0.0.1:9200
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}

	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		return err
	}
	if !strings.HasPrefix(esversion, "5.2") {
		return fmt.Errorf("unsupported elasticsearch version: %s", esversion)
	}

	orm.esClient = client
	orm.ctx = ctx

	return nil
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

func (orm *Orm) CreatePercolatorDocument(obj Elasticable, jsonQuery string) (common.Ident, error) {
	if obj.GetId() != "" {
		return "", fmt.Errorf("ID already assigned prior to Create()")
	}

	resp, err := orm.esClient.Index().
		Index(GetIndexName(obj)).
		Type("queries").
		Id(obj.SetId().String()).
		BodyJson(jsonQuery).
		Refresh("wait_for").
		Do(orm.ctx)

	if err != nil {
		return "", err
	}
	if !resp.Created {
		return "", fmt.Errorf("Create() did not create")
	}

	return common.ToIdent(resp.Id), nil

}

func (orm *Orm) CreatePercolatorQuery(obj Elasticable) ([]Elasticable, int64, error) {
	pq := elastic.NewPercolatorQuery().
		Field("query").
		DocumentType(GetTypeName(obj)).
		Document(obj)
	result, err := orm.esClient.Search(GetIndexName(obj)).Query(pq).Do(orm.ctx)

	if result.Hits.TotalHits <= 0 {
		return nil, 0, nil
	}

	ary := []Elasticable{}

	i := 0
	for _, hit := range result.Hits.Hits {
		tmp := common.NewViaReflection(obj)

		err = json.Unmarshal(*hit.Source, tmp)
		if err != nil {
			return nil, 0, err
		}
		ary = append(ary, tmp.(Elasticable))
		i++
	}

	return ary, result.Hits.TotalHits, nil
}

func (orm *Orm) ReadDocument(typ Elasticable) (Elasticable, error) {
	if typ.GetId() == common.NoIdent {
		return nil, fmt.Errorf("object does not have Id set")
	}

	result, err := orm.esClient.Get().
		Index(GetIndexName(typ)).
		Type(GetTypeName(typ)).
		Id(typ.GetId().String()).
		Do(orm.ctx)
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, fmt.Errorf("document not found")
	}

	typ2 := common.NewViaReflection(typ)
	err = json.Unmarshal(*result.Source, typ2)
	if err != nil {
		return nil, err
	}

	return typ2.(Elasticable), nil
}

// TODO: for now, always return sorted by id (ascending)
func (orm *Orm) ReadDocuments(typ Elasticable, from int, size int) ([]Elasticable, int64, error) {

	result, err := orm.esClient.Search().
		Index(GetIndexName(typ)).
		Type(GetTypeName(typ)).
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

	ary := []Elasticable{}

	i := 0
	for _, hit := range result.Hits.Hits {
		tmp := common.NewViaReflection(typ)

		err = json.Unmarshal(*hit.Source, tmp)
		if err != nil {
			return nil, 0, err
		}
		ary = append(ary, tmp.(Elasticable))
		i++
	}

	return ary, result.Hits.TotalHits, nil
}

func (orm *Orm) UpdateDocument(src Elasticable) error {
	dest, err := orm.ReadDocument(src)
	if err != nil {
		return err
	}

	err = CrudMerge(src, dest, CrudFieldUpdate)
	if err != nil {
		return err
	}

	_, err = orm.esClient.Update().
		Index(GetIndexName(src)).
		Type(GetTypeName(src)).
		Id(src.GetId().String()).
		Doc(dest).
		Do(orm.ctx)
	return err
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

func (orm *Orm) GetIndexes() ([]string, error) {
	names, err := orm.esClient.IndexNames()
	if err != nil {
		return nil, err
	}
	return names, nil
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

func (orm *Orm) CreateIndex(e Elasticable, usePercolation bool) error {
	index := GetIndexName(e)

	exists, err := orm.IndexExists(e)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index already exists")
	}

	mapping := NewMapping(e, usePercolation)
	byts, err := json.Marshal(mapping)

	if err != nil {
		return err
	}
	mappingString := string(byts)

	log.Printf("[[[ %s ]]]", mappingString)

	if strings.Contains(mappingString, "string") {
		panic("internal error: obselete datatype \"string\" in mapping for index " + index)
	}

	result, err := orm.esClient.CreateIndex(index).BodyString(mappingString).Do(orm.ctx)
	if err != nil {
		return err
	}

	if !result.Acknowledged {
		return fmt.Errorf("CreateIndex() not acknowledged")
	}

	return nil
}

func (orm *Orm) CreateIndexWithPercolation(e Elasticable) error {
	index := GetIndexName(e)

	exists, err := orm.IndexExists(e)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index already exists")
	}

	/*{
		"settings":{},
		"mappings":{
			"doc_type":{
				"dynamic":"strict",
				"properties":{
					"id":{"type":"keyword"},
					"message":{"type":"text"}
				}
			},
			"queries":{
				"properties":{
					"query":{"type":"percolator"}
				}
			}
		}
	}*/

	mapping := NewMapping(e, true)
	//	mapping.Mappings["queries"] = MappingProperty{
	////		Properties: map[string]MappingProperty{
	///		"query": MappingProperty{
	//		Type: "percolator",
	//	},
	//	},
	//}
	byts, err := json.Marshal(mapping)

	if err != nil {
		return err
	}
	mappingString := string(byts)
	log.Printf(">> %s <<", mappingString)
	if strings.Contains(mappingString, "string") {
		panic("internal error: obselete datatype \"string\" in mapping for index " + index)
	}

	result, err := orm.esClient.CreateIndex(index).BodyString(mappingString).Do(orm.ctx)
	if err != nil {
		return err
	}

	if !result.Acknowledged {
		return fmt.Errorf("CreateIndex() not acknowledged")
	}

	return nil
}

func (orm *Orm) AddPercolationQuery(obj Elasticable) error {
	index := GetIndexName(obj)

	_, err := orm.esClient.Index().
		Index(index).
		Type("query").
		Id(obj.SetId().String()).
		BodyJson(`{"query":{"match":{"message":"bonsai tree"}}}`).
		Refresh("wait_for").
		Do(orm.ctx)

	return err
}

func (orm *Orm) AddPercolationDocument(obj Elasticable) error {
	index := GetIndexName(obj)
	typ := GetTypeName(obj)

	pq := elastic.NewPercolatorQuery().
		Field("query").
		DocumentType(typ).
		Document(obj)
	res, err := orm.esClient.Search(index).Query(pq).Do(orm.ctx)
	if err != nil {
		return err
	}
	log.Printf("[[ %#v ]]", res)
	return nil
}
