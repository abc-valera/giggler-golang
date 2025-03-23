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
