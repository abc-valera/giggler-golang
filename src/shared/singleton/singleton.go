package singleton

func New[ValueType any, InitFuncType func() ValueType](InitializeFunc InitFuncType) func() ValueType {
	val := InitializeFunc()
	return func() ValueType {
		return val
	}
}
