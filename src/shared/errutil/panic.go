package errutil

// NoErr panics if err is not nil
func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Must panics if err is not nil
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// NoEmpty panics if the provided value is zero
func NoEmpty[T comparable](val T) T {
	var zeroValue T
	if val == zeroValue {
		panic("empty value")
	}
	return val
}
