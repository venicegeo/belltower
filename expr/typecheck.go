package expr

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
func TypeCheck(a interface{}, b map[string]interface{}) (interface{}, []string) {
	errs := []string{}

	s := structs.New(a)

	keyUsed := map[string]bool{}
	for k, _ := range b {
		keyUsed[k] = false
	}

	for _, name := range s.Names() {
		_, ok := b[name]
		if !ok {
			errs = append(errs, fmt.Sprintf("struct field %s: not found in map", name))
			continue
		}

		keyUsed[name] = true

		ok, errMsg := sameType(s.Field(name).Value(), b[name])
		if !ok {
			errs = append(errs, fmt.Sprintf("struct field %s: type error, %s", name, errMsg))
			continue
		}

		// success!
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

type typeCode int

const (
	typeInt typeCode = iota
	typeFloat
	typeString
	typeBool
)

var typeEnum map[typeCode]string = map[typeCode]string{
	typeInt:    "INT",
	typeFloat:  "FLOAT",
	typeString: "STRING",
	typeBool:   "BOOL",
}

func getType(a interface{}) typeCode {
	switch a.(type) {
	case int:
		return typeInt
	case float64:
		return typeFloat
	case float32:
		return typeFloat
	case string:
		return typeString
	case bool:
		return typeBool
	}
	panic(12)
}

func sameType(a interface{}, b interface{}) (bool, string) {
	aType := getType(a)
	bType := getType(b)

	if aType != bType {
		return false, fmt.Sprintf("expected %s but received %s", typeEnum[aType], typeEnum[bType])
	}
	return true, ""
}
