/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
