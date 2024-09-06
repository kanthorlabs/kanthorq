package utils

import "reflect"

func NameOf(v interface{}) string {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
