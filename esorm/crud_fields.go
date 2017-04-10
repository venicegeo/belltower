package esorm

import (
	"fmt"
	"reflect"
	"strings"
)

type CrudFieldMode string

const (
	CrudFieldCreate CrudFieldMode = "c"
	CrudFieldRead                 = "r"
	CrudFieldUpdate               = "u"
	CrudFieldAll                  = CrudFieldCreate + CrudFieldRead + CrudFieldUpdate

	CrudTag string = "crud"
)

func hasMode(modes, mode CrudFieldMode) bool {
	return strings.Contains(string(modes), string(mode))
}

func IsCrudField(obj interface{}, fieldName string, fieldMode CrudFieldMode) (bool, error) {

	// bad mode
	if len(fieldMode) != 1 ||
		!hasMode(CrudFieldAll, fieldMode) {
		return false, fmt.Errorf("Invalid field mode: %s", fieldMode)
	}

	srcPtrToStruct := reflect.ValueOf(obj)
	srcStruct := srcPtrToStruct.Elem()
	// bad type
	if srcStruct.Kind() != reflect.Struct {
		return false, fmt.Errorf("Type is not a valid struct: %s", srcStruct.Kind())
	}

	// bad field name
	if fieldName == "" {
		return false, fmt.Errorf("Invalid field name: %s", fieldName)
	}

	field, ok := srcStruct.Type().FieldByName(fieldName)
	if !ok {
		return false, fmt.Errorf("Invalid field name: %s", fieldName)
	}

	modes, ok := field.Tag.Lookup(CrudTag)
	if !ok {
		return false, nil
	}
	if modes == "" {
		return false, nil
	}
	if !hasMode(CrudFieldMode(modes), fieldMode) {
		return false, nil
	}

	return true, nil
}

func CrudMerge(src interface{}, dest interface{}, mode CrudFieldMode) error {

	srcPtrToStruct := reflect.ValueOf(src)
	destPtrToStruct := reflect.ValueOf(dest)
	if srcPtrToStruct.Type() != destPtrToStruct.Type() {
		return fmt.Errorf("incompatible types: %#v %#v", srcPtrToStruct, destPtrToStruct)
	}

	srcStruct := srcPtrToStruct.Elem()
	destStruct := destPtrToStruct.Elem()
	if srcStruct.Kind() != destStruct.Kind() {
		return fmt.Errorf("incompatible types: %#v %#v", srcStruct.Kind(), destStruct.Kind())
	}
	if srcStruct.Kind() != reflect.Struct {
		return fmt.Errorf("incorrect types")
	}

	numField := srcStruct.NumField()
	for i := 0; i < numField; i++ {
		field := srcStruct.Type().Field(i)
		fieldName := field.Name
		fieldValue := srcStruct.Field(i)

		xfieldValue := destStruct.Field(i)

		if !fieldValue.IsValid() {
			return fmt.Errorf("invalid")
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("can't set")
		}

		ok, err := IsCrudField(src, fieldName, mode)
		if err != nil {
			return err
		}
		if ok {
			xfieldValue.Set(fieldValue)
		}
	}
	return nil
}
