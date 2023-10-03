package validate

func Struct(s any) error {
	return validate.Struct(s)
}

func Var(v any, tag string) error {
	return validate.Var(v, tag)
}
