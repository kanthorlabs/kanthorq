package entities

import (
	"fmt"
	"reflect"
)

func Collection(name string) string {
	return fmt.Sprintf("kanthorq_%s", name)
}

func Properties(entity any) []string {
	typeof := reflect.TypeOf(entity)

	// Check if the entity is a pointer, and get the underlying element if it is
	if typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
	}

	var props []string

	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		props = append(props, field.Tag.Get("json"))
	}

	return props
}
