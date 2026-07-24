package validate

import (
	"github.com/go-playground/validator/v10"

	"giggler-golang/src/shared/enum"
	"giggler-golang/src/shared/singleton"
)

var getValidator = singleton.New(func() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// IEnum validation
	validate.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(enum.Interface)
		if !ok {
			panic("enum validation must be used on a field that implements enum.IEnum")
		}

		return value.IsValid()
	})

	// Other custom validations can be defined here...

	return validate
})

func Struct(s any) error {
	return getValidator().Struct(s)
}

func StructPartial(s any, fields ...string) error {
	return getValidator().StructPartial(s, fields...)
}

func Var(v any, tag string) error {
	return getValidator().Var(v, tag)
}
