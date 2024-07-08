package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type stringv struct {
	Topic  string `validate:"is_topic"`
	Stream string `validate:"is_collection_name"`
}

func TestStream(t *testing.T) {
	value := stringv{Topic: "system.ping", Stream: "default"}
	require.NoError(t, Validate.Struct(value))

	value.Stream = ""
	value.Topic = "system.ping"
	cerr := Validate.Struct(value).(validator.ValidationErrors)[0]
	require.Equal(t, "is_collection_name", cerr.Tag(), "should have only have collection name error")

	value.Stream = "default"
	value.Topic = ""
	terr := Validate.Struct(value).(validator.ValidationErrors)[0]
	require.Equal(t, "is_topic", terr.Tag(), "should have only have collection name error")

	value.Stream = ""
	value.Topic = ""
	errs := Validate.Struct(value).(validator.ValidationErrors)
	require.Equal(t, 2, len(errs), "should have both topic and collection name error")
}
