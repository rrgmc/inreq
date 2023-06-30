package inreq

import "reflect"

// reflectElem returns the first non-pointer type from the [reflect.Type].
func reflectElem(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
