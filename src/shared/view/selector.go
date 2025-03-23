package view

import (
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
)

func NewSelector(optLimit viewgen.OptInt, optOffset viewgen.OptInt) data.Selector {
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

	return data.Selector{
		PagingLimit:  limit,
		PagingOffset: offset,
	}
}
