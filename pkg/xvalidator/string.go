package xvalidator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	RegexCollectionName = regexp.MustCompile(`^[a-zA-Z0-9]+(_[a-zA-Z0-9]+)*$`)
	RegexSubject        = regexp.MustCompile(`^[a-zA-Z0-9]+([._-]?[a-zA-Z0-9]+)*$`)
)

func collection(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return RegexCollectionName.MatchString(value)
}

func subject(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return false
	}

	return RegexSubject.MatchString(value)
}

// @TODO: implement me
func subjectFilter(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value != ""
}
