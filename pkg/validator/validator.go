package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate
var once sync.Once

func init() {
	once.Do(func() {
		Validate = validator.New()
		Validate.RegisterValidation("is_topic", topic)
		Validate.RegisterValidation("is_collection_name", collection)
	})
}
