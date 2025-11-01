package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"testing"
)

type Inner struct {
	Value string `validate:"max=5"`
}

type Outer1 struct {
	Inner map[string]Inner `validate:"dive,keys,max=5,endkeys"`
}

type Outer2 struct {
	Inner map[string]Inner `validate:"dive"`
}

type Outer3 struct {
	Inner map[string]Inner `validate:"dive,keys,max=5,endkeys,required"`
}

func TestValidateStruct(t *testing.T) {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Should fail, but passes
	fmt.Println("Outer1:", v.Struct(Outer1{Inner: map[string]Inner{
		"a": {Value: "1234567890"},
	}}))

	// Should fail, and fails (but there is no validation for keys)
	fmt.Println("Outer2:", v.Struct(Outer2{Inner: map[string]Inner{
		"a": {Value: "1234567890"},
	}}))

	// Should fail, and fails (but has extra required validation, which is not intended)
	fmt.Println("Outer3:", v.Struct(Outer3{Inner: map[string]Inner{
		"a": {Value: "1234567890"},
	}}))

	// Should fail (due to required), but I don't want it to fail
	fmt.Println("Outer3 (a):", v.Struct(Outer3{Inner: map[string]Inner{
		"a": {Value: ""},
	}}))
}
