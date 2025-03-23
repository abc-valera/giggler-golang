package errutil

// NoErr stops program execution if err is not nil
func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Must stops program execution if err is not nil
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// NoEmpty stops program execution if the provided value is empty/nil
func NoEmpty[T comparable](val T) T {
	var nullValue T
	if val == nullValue {
		panic("empty value")
	}
	return val
}
