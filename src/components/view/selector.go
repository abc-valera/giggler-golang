package view

import (
	"github.com/abc-valera/giggler-golang/src/components/viewgen"
	"github.com/abc-valera/giggler-golang/src/shared/ds"
)

func NewSelector(optLimit viewgen.OptInt, optOffset viewgen.OptInt) ds.Selector {
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

	return ds.Selector{
		PagingLimit:  limit,
		PagingOffset: offset,
	}
}
