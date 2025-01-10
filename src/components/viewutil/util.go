package viewutil

import (
	"github.com/abc-valera/giggler-golang/gen/view"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
)

func NewString(optString view.OptString) *string {
	if optString.Set {
		return &optString.Value
	}
	return nil
}

func NewSelector(optLimit view.OptInt, optOffset view.OptInt) persistence.Selector {
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

	return persistence.Selector{
		PagingLimit:  limit,
		PagingOffset: offset,
	}
}
