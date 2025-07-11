package validate

import (
	"github.com/go-playground/validator/v10"

	"giggler-golang/src/shared/enum"
)

var validate = func() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	{
		// IEnum validation
		validateIEnum := func(fl validator.FieldLevel) bool {
			value, ok := fl.Field().Interface().(enum.Interface)
			if !ok {
				panic("enum validation must be used on a field that implements enum.IEnum")
			}

			return value.IsValid()
		}
		validate.RegisterValidation("enum", validateIEnum)
	}

	{
		// Other custom validations can be defined here...
	}

	return validate
}()

func Struct(s any) error {
	return validate.Struct(s)
}

func Var(v any, tag string) error {
	return validate.Var(v, tag)
}
