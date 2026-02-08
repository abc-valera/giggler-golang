package must

// Do panics if err is not nil
func Do[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// NoErr panics if err is not nil
func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

// NotZeroValue panics if the provided value is zero
func NotZeroValue[T comparable](val T) T {
	var zeroValue T
	if val == zeroValue {
		panic("empty value")
	}
	return val
}
