package common

//---------------------------------------------------------------------

type Map map[string]interface{}

func (m Map) Has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Map) IsInt(key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	_, ok = v.(int)
	return ok
}

func (m Map) AsInt(key string) int {
	if !m.IsInt(key) {
		return 0
	}
	return m[key].(int)
}

func (m Map) IsString(key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	_, ok = v.(string)
	return ok
}

func (m Map) AsString(key string) string {
	if !m.IsString(key) {
		return ""
	}
	return m[key].(string)
}

func (m Map) AsValidString(key string) (string, bool) {
	s := m.AsString(key)
	if s == "" {
		return "", false
	}
	return s, true
}

//---------------------------------------------------------------------

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
