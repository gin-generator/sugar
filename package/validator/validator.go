package validator

import (
	"errors"
	_validator "github.com/go-playground/validator/v10"
	"regexp"
)

// validate
/**
 * @description: validator instance
 */
var validate *_validator.Validate

func init() {
	// 启用 RequiredStructEnabled 以自动验证嵌套结构体
	validate = _validator.New(_validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("phone", validatePhone)
	if err != nil {
		panic("Unable to register validator, error: " + err.Error())
	}
}

// validatePhone
func validatePhone(f _validator.FieldLevel) bool {
	phone := f.Field().String()
	regx := `^1[3-9]\d{9}$`
	return regexp.MustCompile(regx).MatchString(phone)
}

// ValidateStruct
/**
 * @description: validate struct
 */
func ValidateStruct(s interface{}) (err error) {
	err = validate.Struct(s)
	if err != nil {
		var validationErrors _validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				return e
			}
		}
	}
	return err
}
