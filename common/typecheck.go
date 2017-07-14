package common

import (
	"fmt"

	"github.com/fatih/structs"
)

// TypeCheck tries to determine if a map will fit into a struct neatly.
//
// a is a object of struct type T declared at compile time
//    T{
//        w string
//        x int
//        y string
//    }
// b is a map of values made at runtime
//    b := map[string]interface{}{
//        "w": 12.34
//        "x": 5
//        "z": true
//    }
// then TypeCheck(a,b) should return a T object like this
//    t.w = ""
//    t.x = 5
//    t.y = ""
// with these error messages:
//    t.w expected string but got float
//    t.y not set
//    b.z not used
func typeCheck(a interface{}, b map[string]interface{}) (interface{}, []string) {
	errs := []string{}

	s := structs.New(a)

	// to track if any of b's keys were not used
	keyUsed := map[string]bool{}
	for k, _ := range b {
		keyUsed[k] = false
	}

	for _, name := range s.Names() {

		// do we have a's required field?
		_, ok := b[name]
		if !ok {
			errs = append(errs, fmt.Sprintf("struct field %s: not found in map", name))
			continue
		}

		// yes, we used b's key
		keyUsed[name] = true

		// but, do the fields from a and b have the same type?
		ok, errMsg := sameType(s.Field(name).Value(), b[name])
		if !ok {
			errs = append(errs, fmt.Sprintf("struct field %s: type error, %s", name, errMsg))
			continue
		}

		// success, set the value into a
		err := s.Field(name).Set(b[name])
		if err != nil {
			errs = append(errs, fmt.Sprintf("struct field %s: can't set (%s)", name, err))
			continue
		}
	}

	for k, v := range keyUsed {
		if !v {
			errs = append(errs, fmt.Sprintf("map field %s: not found in struct", k))
		}
	}

	return a, errs
}

func getType(a interface{}) TypeName {
	switch a.(type) {
	case int:
		return TypeNameInt
	case float32, float64:
		return TypeNameFloat
	case string:
		return TypeNameString
	case bool:
		return TypeNameBool
	}
	panic(12)
}

func sameType(a interface{}, b interface{}) (bool, string) {
	aType := getType(a)
	bType := getType(b)

	if aType != bType {
		return false, fmt.Sprintf("expected %s but received %s", aType, bType)
	}
	return true, ""
}
