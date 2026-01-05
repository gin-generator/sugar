package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"testing"
)

type Inner struct {
	Value string            `validate:"max=5"`
	Name  map[string]Nested `validate:"dive,keys,endkeys,required"`
}

type Nested struct {
	Name string `validate:"required,max=1"`
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
	v := validator.New()

	// Should fail, and fails (but there is no validation for keys)
	fmt.Println("Outer2:", v.Struct(Outer2{Inner: map[string]Inner{
		"a": {
			Value: "123",
			Name: map[string]Nested{
				"x": {Name: "validname"},
			},
		},
	}}))
}
