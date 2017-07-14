package esorm

import (
	"encoding/json"
)

type Filter map[string]interface{}

func Match(field string, value string) *Filter {
	f := Filter{}
	f["match"] = map[string]interface{}{
		field: value,
	}
	return &f
}

func Or(args ...*Filter) *Filter {
	f := Filter{}
	f["bool"] = map[string]interface{}{
		"should": args,
	}
	return &f
}

func And(args ...*Filter) *Filter {
	f := Filter{}
	f["bool"] = map[string]interface{}{
		"must": args,
	}
	return &f
}

func (f *Filter) toJson() string {

	type T map[string]interface{}
	m := T{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": f,
			},
		},
	}
	byts, err := json.Marshal(m)
	if err != nil {
		panic(99) // TODO
	}

	return string(byts)
}
