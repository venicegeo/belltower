package common

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

//---------------------------------------------------------------------

// TODO: remove these

/*func AsInt(x interface{}) *int {
	v, ok := x.(int)
	if ok {
		return &v
	}
	return nil
}*/

func AsFloat(x interface{}) *float64 {
	v, ok := x.(float64)
	if ok {
		return &v
	}
	return nil
}

func AsBool(x interface{}) *bool {
	v, ok := x.(bool)
	if ok {
		return &v
	}
	return nil
}

func AsString(x interface{}) *string {
	v, ok := x.(string)
	if ok {
		return &v
	}
	return nil
}
