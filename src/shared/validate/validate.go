package validate

import (
	"github.com/abc-valera/template-golang/src/components/enum"
	"github.com/go-playground/validator/v10"
)

var validate = func() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Define and register custom validation functions here:

	// IEnum validation
	validateIEnum := func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(enum.Interface)
		if !ok {
			panic("enum validation must be used on a field that implements enum.IEnum")
		}

		return value.IsValid()
	}
	validate.RegisterValidation("enum", validateIEnum)

	return validate
}()

func Struct(s any) error {
	return validate.Struct(s)
}

func Var(v any, tag string) error {
	return validate.Var(v, tag)
}
