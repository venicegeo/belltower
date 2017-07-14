package esorm

import (
	"fmt"

	"github.com/venicegeo/belltower/common"
)

//---------------------------------------------------------------------

type FakeOrm struct {
	indexes map[string]*FakeIndex
	isOpen  bool
	//settings interface{}
}

type FakeIndex struct {
	name  string
	types map[string]*FakeType
}

type FakeType struct {
	name string

	// maps from id string to document body
	items   map[string]Elasticable
	mapping interface{}
}

// implements Elasticable
type FakeDocument struct {
	id   common.Ident
	data string
}

func (d *FakeDocument) GetId() common.Ident                              { return d.id }
func (d *FakeDocument) SetId() common.Ident                              { d.id = common.NewId(); return d.id }
func (d *FakeDocument) String() string                                   { return "string" }
func (d *FakeDocument) GetMappingProperties() map[string]MappingProperty { return nil }

func init() {
	var _ Ormer = (*FakeOrm)(nil)
}

func (orm *FakeOrm) Open() error {
	orm.indexes = map[string]*FakeIndex{}
	orm.isOpen = true
	return nil
}

func (orm *FakeOrm) Close() error {
	// make lists of things to delete
	idxs := map[string]map[string]bool{}
	for _, idx := range orm.indexes {
		idxs[idx.name] = map[string]bool{}
		for _, typ := range idx.types {
			idxs[idx.name][typ.name] = true
		}
	}
	// now delete them
	for i, idx := range idxs {
		for j, _ := range idx {
			delete(idx, j)
		}
		delete(idxs, i)
	}
	if len(idxs) != 0 {
		panic(99)
	}

	orm.isOpen = false

	return nil
}

func (orm *FakeOrm) CreateDocument(obj Elasticable) (common.Ident, error) {
	idx := obj.GetIndexName()
	typ := obj.GetTypeName()
	id := common.NewId()
	obj.SetId(id)

	destIdx, ok := orm.indexes[idx]
	if !ok {
		return common.NoIdent, fmt.Errorf("index %s does not exist", idx)
	}
	if destIdx.types == nil {
		destIdx.types = map[string]*FakeType{}
	}
	destTyp, ok := destIdx.types[typ]
	if !ok {
		destIdx.types[typ] = &FakeType{name: typ, items: map[string]Elasticable{}}
		destTyp = destIdx.types[typ]
	}
	destTyp.items[id.String()] = obj
	return id, nil
}

func (orm *FakeOrm) ReadDocument(obj Elasticable) (Elasticable, error) {
	if obj.GetId() == common.NoIdent {
		return nil, fmt.Errorf("object does not have Id set")
	}

	idx := obj.GetIndexName()
	typ := obj.GetTypeName()
	id := obj.GetId()

	destIdx, ok := orm.indexes[idx]
	if !ok {
		return nil, fmt.Errorf("index %s does not exist", idx)
	}
	destTyp, ok := destIdx.types[typ]
	if !ok {
		return nil, fmt.Errorf("type %s in index %s does not exist", typ, idx)
	}

	doc, ok := destTyp.items[id.String()]
	if !ok {
		return nil, fmt.Errorf("document not found")
	}

	return doc, nil
}

// TODO: for now, always return sorted by id (ascending)
func (orm *FakeOrm) ReadDocuments(obj Elasticable, from int, size int) ([]Elasticable, int64, error) {

	idx := obj.GetIndexName()
	typ := obj.GetTypeName()

	destIdx, ok := orm.indexes[idx]
	if !ok {
		return nil, 0, fmt.Errorf("index %s does not exist", idx)
	}
	destTyp, ok := destIdx.types[typ]
	if !ok {
		return nil, 0, fmt.Errorf("type %s in index %s does not exist", typ, idx)
	}

	totalHits := len(destTyp.items)

	ids := make([]string, totalHits)
	i := 0
	for k, _ := range destTyp.items {
		ids[i] = k
		i++
	}

	// TODO: sort the list

	s := from
	if s > totalHits {
		return []Elasticable{}, int64(totalHits), nil
	}
	e := from + size
	if e > totalHits {
		e = totalHits
	}

	slice := ids[s:e]
	if len(slice) == 0 {
		return []Elasticable{}, int64(totalHits), nil
	}

	ary := []Elasticable{}

	for _, hit := range slice {
		ary = append(ary, destTyp.items[hit])
	}

	return ary, int64(totalHits), nil
}

func (orm *FakeOrm) UpdateDocument(src Elasticable) error {
	idx := src.GetIndexName()
	typ := src.GetTypeName()
	id := src.GetId()

	destIdx, ok := orm.indexes[idx]
	if !ok {
		return fmt.Errorf("index %s does not exist", idx)
	}
	destTyp, ok := destIdx.types[typ]
	if !ok {
		return fmt.Errorf("type %s in index %s does not exist", typ, idx)
	}
	_, ok = destTyp.items[id.String()]
	if !ok {
		return fmt.Errorf("document %s.%s.%s does nto exist", idx, typ, id)
	}
	destTyp.items[id.String()] = src

	return nil
}

func (orm *FakeOrm) DeleteDocument(obj Elasticable) error {
	idx := obj.GetIndexName()
	typ := obj.GetTypeName()
	id := obj.GetId()

	destIdx, ok := orm.indexes[idx]
	if !ok {
		return fmt.Errorf("index %s does not exist", idx)
	}
	destTyp, ok := destIdx.types[typ]
	if !ok {
		return fmt.Errorf("type %s in index %s does not exist", typ, idx)
	}

	_, ok = destTyp.items[id.String()]
	if !ok {
		return fmt.Errorf("document not found")
	}
	delete(destTyp.items, id.String())
	return nil
}

func (orm *FakeOrm) GetIndexes() ([]string, error) {
	nams := []string{}
	for _, v := range orm.indexes {
		nams = append(nams, v.name)
	}
	return nams, nil
}

func (orm *FakeOrm) IndexExists(e Elasticable) (bool, error) {
	_, ok := orm.indexes[e.GetIndexName()]
	return ok, nil
}

func (orm *FakeOrm) DeleteIndex(e Elasticable) error {
	_, ok := orm.indexes[e.GetIndexName()]
	if !ok {
		return fmt.Errorf("index %s does not exists", e.GetIndexName())
	}
	delete(orm.indexes, e.GetIndexName())
	return nil
}

func (orm *FakeOrm) CreateIndex(e Elasticable, usePercolation bool) error {
	_, ok := orm.indexes[e.GetIndexName()]
	if ok {
		return fmt.Errorf("index %s already exists", e.GetIndexName())
	}
	orm.indexes[e.GetIndexName()] = &FakeIndex{name: e.GetIndexName()}
	return nil
}
