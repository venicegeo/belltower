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

	"github.com/mitchellh/mapstructure"
)

//---------------------------------------------------------------------

type Serializer interface {
	ReadFromJSON(string) error
	WriteToJSON() (string, error)
	Validate() error
}

func ReadFromJSON(jsn string, obj interface{}) error {

	// TODO: zero out obj (if is struct, use Structs pkg)

	err := json.Unmarshal([]byte(jsn), obj)
	if err != nil {
		return err
	}

	return nil
}

func WriteToJSON(obj interface{}) (string, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

//---------------------------------------------------------------------

func SetStructFromMap(input map[string]interface{}, result interface{}, weakly bool) (interface{}, error) {

	var md mapstructure.Metadata
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   result,
	}
	if weakly {
		config.WeaklyTypedInput = true
		config.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
