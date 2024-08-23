package xvalidator

import (
	"github.com/go-playground/validator/v10"
)

func enum(fl validator.FieldLevel) bool {
	// no String method, should return error
	method := fl.Field().MethodByName("String")
	if !method.IsValid() {
		return false
	}

	out := method.Call(nil)

	// no output, should return error
	if len(out) == 0 {
		return false
	}

	// don't accept default value
	return out[0].String() != ""
}
