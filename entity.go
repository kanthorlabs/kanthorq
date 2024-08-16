package kanthorq

import (
	"fmt"
	"reflect"
)

func Collection(name string) string {
	return fmt.Sprintf("kanthorq_%s", name)
}

func Properties(entity any) []string {
	var props []string
	eventType := reflect.TypeOf(entity)

	for i := 0; i < eventType.NumField(); i++ {
		field := eventType.Field(i)
		props = append(props, field.Tag.Get("json"))
	}

	return props
}
