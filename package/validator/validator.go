package validator

import (
	"errors"
	"fmt"
	_validator "github.com/go-playground/validator/v10"
	"regexp"
)

// validate
/**
 * @description: validator instance
 */
var validate *_validator.Validate

func init() {
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
		var invalidValidationError *_validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return
		}
		for _, obj := range err.(_validator.ValidationErrors) {
			return errors.New(
				fmt.Sprintf("In %s, Field '%s' validation failed on the '%s' tag",
					obj.StructNamespace(), obj.Field(), obj.Tag()))
		}
	}
	return
}
