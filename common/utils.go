package common

import (
	"fmt"
	"reflect"
	"time"
)

func GetMapValueAsInt(m map[string]interface{}, key string) (int, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("Missing config field '%s'", key)
	}
	ret, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("Config field '%s' not legal integer", key)
	}
	return ret, nil
}

func GetMapValueAsBool(m map[string]interface{}, key string) (bool, error) {
	value, ok := m[key]
	if !ok {
		return false, fmt.Errorf("Missing config field '%s'", key)
	}
	ret, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("Config field '%s' not legal boolean", key)
	}
	return ret, nil
}

func GetMapValueAsString(m map[string]interface{}, key string) (string, error) {
	value, ok := m[key]
	if !ok {
		return "", fmt.Errorf("Missing config field '%s'", key)
	}
	ret, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("Config field '%s' not legal string", key)
	}
	return ret, nil
}

func GetMapValueAsFloat(m map[string]interface{}, key string) (float64, error) {
	value, ok := m[key]
	if !ok {
		return 0.0, fmt.Errorf("Missing config field '%s'", key)
	}
	ret, ok := value.(float64)
	if !ok {
		return 0.0, fmt.Errorf("Config field '%s' not legal float", key)
	}
	return ret, nil
}

func GetMapValueAsDuration(m map[string]interface{}, key string) (time.Duration, error) {
	value, ok := m[key]
	if !ok {
		return time.Duration(0), fmt.Errorf("Missing config field '%s'", key)
	}
	t, ok := value.(time.Duration)
	if !ok {
		return time.Duration(0), fmt.Errorf("Config field '%s' not legal time.Duration: %v", key, value)
	}
	//	ret, err := time.ParseDuration(t)
	//	if err != nil {
	//		return time.Duration(0), err
	//	}
	//	return ret, nil
	return t, nil
}

// ObjectsAreEqual determines if two objects are considered equal.
//
// taken from testify's assertions.go
func ObjectsAreEqual(expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}

	return reflect.DeepEqual(expected, actual)
}

func MapsAreEqualValues(expected, actual interface{}) bool {

	e, eok := expected.(map[string]interface{})
	a, aok := actual.(map[string]interface{})
	if !eok || !aok {
		return false
	}

	if len(e) != len(a) {
		return false
	}

	for ek, ev := range e {
		av, ok := a[ek]
		if !ok {
			return false
		}
		ok = ObjectsAreEqualValues(ev, av)
		if !ok {
			return false
		}
	}

	return true
}

// ObjectsAreEqualValues gets whether two objects are equal, or if their
// values are equal.
//
// taken from testify's assertions.go
func ObjectsAreEqualValues(expected, actual interface{}) bool {
	if ObjectsAreEqual(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		// Attempt comparison after type conversion
		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
	}

	return false
}
