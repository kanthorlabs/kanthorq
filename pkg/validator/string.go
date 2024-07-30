package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	RegexCollectionName = regexp.MustCompile(`^[a-zA-Z0-9]+(_[a-zA-Z0-9]+)*$`)
	RegexTopicPart      = regexp.MustCompile(`^[a-zA-Z0-9]+(_-[a-zA-Z0-9]+)*$`)
	RegexTopicPartLast  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]*[a-zA-Z0-9*]$`)
)

func collection(fl validator.FieldLevel) bool {
	return RegexCollectionName.MatchString(fl.Field().String())
}

func topic(fl validator.FieldLevel) bool {
	parts := strings.Split(fl.Field().String(), ".")
	if len(parts) == 0 {
		return false
	}
	if len(parts) == 1 {
		return RegexTopicPartLast.MatchString(parts[0])
	}

	// match all parts except the last one
	for i := 0; i < len(parts)-1; i++ {
		if !RegexTopicPart.MatchString(parts[i]) {
			return false
		}
	}

	// if latest part is *, match it
	if parts[len(parts)-1] == "*" {
		return true
	}

	// match the last part
	return RegexTopicPartLast.MatchString(parts[len(parts)-1])
}
