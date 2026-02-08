package must

import (
	"fmt"
	"reflect"
)

// NoNilInterfaceFields panics if any of the fields of the input struct is a nil interface.
func NoNilInterfaceFields[T any](input T) T {
	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() != reflect.Struct {
		panic("input must be a struct")
	}

	for i := range inputValue.NumField() {
		field := inputValue.Field(i)
		if field.Kind() == reflect.Interface && field.IsNil() {
			panic(fmt.Sprintf("field %s on struct %s is a nil interface",
				inputValue.Type().Field(i).Name,
				inputValue.Type().Name()),
			)
		}
	}

	return input
}
