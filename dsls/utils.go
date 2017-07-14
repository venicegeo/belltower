package dsls

//---------------------------------------------------------------------

type Map map[string]interface{}

func (m Map) has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Map) isString(key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	_, ok = v.(string)
	if !ok {
		return false
	}
	return true
}

func (m Map) asString(key string) string {
	if is
	s, _ := m.asValidString(key)
	return s
}

func (m Map) asValidString(key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
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
