package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type collectionstruct struct {
	Stream string `validate:"is_collection_name"`
}

func TestCollection(t *testing.T) {
	value := collectionstruct{Stream: "default"}
	require.NoError(t, Validate.Struct(value))

	value.Stream = ""
	cerr := Validate.Struct(value).(validator.ValidationErrors)[0]
	require.Equal(t, "is_collection_name", cerr.Tag(), "should have only have collection name error")
}

type topicstruct struct {
	Topic string `validate:"is_topic"`
}

func TestTopic(t *testing.T) {
	value := topicstruct{Topic: "system.ping"}
	require.NoError(t, Validate.Struct(value))

	value.Topic = ""
	terr := Validate.Struct(value).(validator.ValidationErrors)[0]
	require.Equal(t, "is_topic", terr.Tag(), "should have only have collection name error")
}
