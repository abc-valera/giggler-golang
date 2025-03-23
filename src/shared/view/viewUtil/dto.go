package viewUtil

import "giggler-golang/src/shared/view/viewgen"

func NewDomainPointer[T any](value interface{ Get() (v T, ok bool) }) *T {
	if value == nil {
		return nil
	}
	v, ok := value.Get()
	if !ok {
		return nil
	}
	return &v
}

func NewOptString(value *string) viewgen.OptString {
	if value == nil {
		return viewgen.OptString{}
	}
	return viewgen.OptString{Value: *value, Set: true}
}
