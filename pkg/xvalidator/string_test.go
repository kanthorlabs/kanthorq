package xvalidator

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

type subjectstruct struct {
	Subject string `validate:"is_subject"`
}

func TestSubject(t *testing.T) {
	ok := []string{
		"event",
		"event-log",
		"event_log",
		"EVENT.log",
		"event.log",
		"event.log-1",
		"event.log_1",
	}
	for _, v := range ok {
		require.NoError(t, Validate.Struct(subjectstruct{Subject: v}))
	}

	ko := []string{
		"$event.log",
		"event.log+1",
		"_event",
		"-event",
		"event_",
		"event-",
		"event._log",
		"event.-log",
		"event.log_",
		"event.log-",
		"event..log",
		".event.log",
		"event.log.",
	}
	for _, v := range ko {
		require.Error(t, Validate.Struct(collectionstruct{Stream: v}))
	}
}
