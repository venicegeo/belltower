/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package engine

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type ArgMap map[string]interface{}

const errFieldNotFound = "argmap: required field '%s' not found"
const errFieldWrongType = "argmap: field '%s' is of type '%s', but received '%s'"

func NewArgMap(jsn string) (ArgMap, error) {

	m := ArgMap{}
	err := json.Unmarshal([]byte(jsn), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m ArgMap) ToJSON() (string, error) {

	buf := []byte{}
	buf, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (m ArgMap) ToStruct(result interface{}) error {
	_, err := SetStructFromMap(m, result, true)
	return err
}

//---------------------------------------------------------------------

func (m ArgMap) GetInterfaceOrDefault(field string, defalt interface{}) (interface{}, error) {
	v, ok := m[field]
	if !ok {
		return defalt, nil
	}
	return v, nil
}

//---------------------------------------------------------------------

func (m ArgMap) GetFloat(field string) (float64, error) {
	v, ok := m[field]
	if !ok {
		return 0.0, fmt.Errorf(errFieldNotFound, field)
	}
	vf, ok := v.(float64)
	if !ok {
		vi, ok := v.(int)
		if ok {
			return float64(vi), nil
		}
		return 0.0, fmt.Errorf(errFieldWrongType, field, "float64", reflect.TypeOf(v).String())
	}
	return vf, nil
}

func (m ArgMap) GetString(field string) (string, error) {
	v, ok := m[field]
	if !ok {
		return "", fmt.Errorf(errFieldNotFound, field)
	}
	vs, ok := v.(string)
	if !ok {
		return "", fmt.Errorf(errFieldWrongType, field, "string", reflect.TypeOf(v).String())
	}
	return vs, nil
}

func (m ArgMap) GetBool(field string) (bool, error) {
	v, ok := m[field]
	if !ok {
		return false, fmt.Errorf(errFieldNotFound, field)
	}
	vb, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf(errFieldWrongType, field, "bool", reflect.TypeOf(v).String())
	}
	return vb, nil
}

//---------------------------------------------------------------------

func (m ArgMap) GetFloatOrDefault(field string, defalt float64) (float64, error) {
	v, ok := m[field]
	if !ok {
		return defalt, nil
		//return 0.0, fmt.Errorf(errFieldNotFound, field)
	}
	vf, ok := v.(float64)
	if !ok {
		vi, ok := v.(int)
		if ok {
			return float64(vi), nil
		}
		return 0.0, fmt.Errorf(errFieldWrongType, field, "float64", reflect.TypeOf(v).String())
	}
	return vf, nil
}

func (m ArgMap) GetStringOrDefault(field string, defalt string) (string, error) {
	v, ok := m[field]
	if !ok {
		return defalt, nil
		//return "", fmt.Errorf(errFieldNotFound, field)
	}
	vs, ok := v.(string)
	if !ok {
		return "", fmt.Errorf(errFieldWrongType, field, "string", reflect.TypeOf(v).String())
	}
	return vs, nil
}

func (m ArgMap) GetBoolOrDefault(field string, defalt bool) (bool, error) {
	v, ok := m[field]
	if !ok {
		return defalt, nil
		//return false, fmt.Errorf(errFieldNotFound, field)
	}
	vb, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf(errFieldWrongType, field, "bool", reflect.TypeOf(v).String())
	}
	return vb, nil
}
