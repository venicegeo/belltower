package dsls

//---------------------------------------------------------------------

type Map map[string]interface{}

func (m Map) has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Map) isInt(key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	_, ok = v.(int)
	return ok
}

func (m Map) asInt(key string) int {
	if !m.isInt(key) {
		return 0
	}
	return m[key].(int)
}

func (m Map) isString(key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	_, ok = v.(string)
	return ok
}

func (m Map) asString(key string) string {
	if !m.isString(key) {
		return ""
	}
	return m[key].(string)
}

func (m Map) asValidString(key string) (string, bool) {
	s := m.asString(key)
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
