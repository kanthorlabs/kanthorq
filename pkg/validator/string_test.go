package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type collectionstruct struct {
	Stream string `validate:"is_collection_name"`
}

func TestCollection(t *testing.T) {
	ok := []string{
		"event",
		"event_log",
	}
	for _, v := range ok {
		require.NoError(t, Validate.Struct(collectionstruct{Stream: v}))
	}

	ko := []string{
		"_event",
		"event_log_",
		"event-log",
	}
	for _, v := range ko {
		require.Error(t, Validate.Struct(collectionstruct{Stream: v}))
	}
}

type topicstruct struct {
	Topic string `validate:"is_topic"`
}

func TestTopic(t *testing.T) {
	ok := []string{
		"*",
		"event",
		"event-log",
		"event_log",
		"event.log",
		"event.log.*",
		"event.log.created-*",
		"event.log.created_*",
		"event.log.created-o*",
		"event.log.created_o*",
	}
	for _, v := range ok {
		require.NoError(t, Validate.Struct(topicstruct{Topic: v}))
	}

	ko := []string{
		".*",
		"_event",
		"event-log-",
		"event_log_",
		"event.log_",
		"event._*",
		"event._created_*",
		"event.-created-*",
	}
	for _, v := range ko {
		require.Error(t, Validate.Struct(collectionstruct{Stream: v}))
	}
}
