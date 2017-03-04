package common

type JSON string

const NilJSON = ""

func (json *JSON) ToObject(interface{}) error {
	return nil
}

func ToJson(interface{}) (JSON, error) {
	return NilJSON, nil
}
