package common

import (
	"encoding/json"
)

// Json objects are immutable/frozen/const.
type Json struct {
	bytes []byte
	imap  map[string]interface{}
}

func NewJsonFromString(s string) (*Json, error) {
	b := []byte(s)
	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	j := &Json{
		bytes: b,
		imap:  m,
	}
	return j, nil
}

func NewJsonFromMap(m map[string]interface{}) (*Json, error) {
	b, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}
	j := &Json{
		bytes: b,
		imap:  m,
	}
	return j, nil
}

func NewJsonFromObject(obj interface{}) (*Json, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return NewJsonFromString(string(b))
}

func (j *Json) AsMap() map[string]interface{} {
	return j.imap
}

func (j *Json) AsString() string {
	return string(j.bytes)
}

func (j *Json) ToObject(obj interface{}) error {
	return json.Unmarshal(j.bytes, obj)
}

func ValidateJsonString(s string) error {
	//log.Printf("== %s ==", s)

	obj := &map[string]interface{}{}
	return json.Unmarshal([]byte(s), obj)
}