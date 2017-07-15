package common

//---------------------------------------------------------------------

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
