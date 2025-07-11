package viewDTO

import (
	"giggler-golang/src/shared/data/dataModel"
	"giggler-golang/src/shared/view/viewgen"
)

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

func NewSelector(optLimit viewgen.OptInt, optOffset viewgen.OptInt) dataModel.Selector {
	var limit uint
	if optLimit.Set {
		limit = uint(optLimit.Value)
	} else {
		limit = 5
	}

	var offset uint
	if optOffset.Set {
		offset = uint(optOffset.Value)
	} else {
		offset = 0
	}

	return dataModel.Selector{
		Limit:  limit,
		Offset: offset,
	}
}
