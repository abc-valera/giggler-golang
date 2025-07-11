// Package serviceLocator provides a simple implementation of the eponymous design pattern.
// It is used as a dependency injection mechanism.
//
// The so called services can be added and pulled by their type.
// This means that only one service can be stored per type.
//
// The package should be 'disabled' by using the Disable() function
// in the start of the main function to prevent undesirable access.
package serviceLocator

import (
	"reflect"
	"sync"
)

var (
	services   = make(map[string]any)
	isDisabled bool
	mu         sync.Mutex
)

func Set[T any](val T) {
	mu.Lock()
	defer mu.Unlock()

	if isDisabled {
		panic("serviceLocator is disabled, cannot set value")
	}

	services[reflect.TypeOf(val).String()] = val
}

func Get[T any]() T {
	mu.Lock()
	defer mu.Unlock()

	if isDisabled {
		panic("serviceLocator is disabled, cannot get value")
	}

	t := reflect.TypeFor[T]()

	val, exists := services[t.String()]
	if !exists {
		panic("value not found for the type " + t.String())
	}

	return val.(T)
}

func Disable() {
	mu.Lock()
	defer mu.Unlock()

	services = nil
	isDisabled = true
}
