package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	// RegexTopic matches valid topic names that is
	// - contains only alphanumeric characters and hyphens
	// - connected by dots
	// - does not have a dot at the beginning or end
	RegexTopic          = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*(\.[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*)*$`)
	RegexCollectionName = regexp.MustCompile(`^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`)
)

func topic(fl validator.FieldLevel) bool {
	return RegexTopic.MatchString(fl.Field().String())
}

func collection(fl validator.FieldLevel) bool {
	return RegexCollectionName.MatchString(fl.Field().String())
}
