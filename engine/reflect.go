package engine

import (
	"reflect"
)

func NewViaReflection(any interface{}) interface{} {
	typ := reflect.TypeOf(any)
	ptr := typ.Kind() == reflect.Ptr

	var val reflect.Value
	var ifc interface{}

	if ptr {
		// typ is *T
		//log.Printf("typ: %s", typ)
		typ = typ.Elem()
		val = reflect.New(typ)
		//log.Printf("val: %s", val.Type())
		ifc = val.Interface()
		//log.Printf("ifc: %s", reflect.TypeOf(iface))
	} else {
		// typ is T
		//log.Printf("typ: %s", typ)
		val = reflect.New(typ)
		//log.Printf("val: %s", val.Type())
		ifc = val.Elem().Interface()
		//log.Printf("ifc: %s", reflect.TypeOf(iface))
	}

	return ifc
}
