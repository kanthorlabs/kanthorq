package subscriber

import (
	"os"
	"testing"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	options := &Options{
		Connection:                os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName:                entities.DefaultStreamName,
		ConsumerName:              entities.DefaultConsumerName,
		ConsumerSubjectIncludes:   []string{"*.>"},
		ConsumerSubjectExcludes:   []string{},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			Size:        100,
			WaitingTime: 5000,
		},
	}

	_, err := New(options)
	require.NoError(t, err)
}

func TestNew_Validate(t *testing.T) {
	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
	}

	_, err := New(options)
	require.Error(t, err)
}

func TestNewRetry(t *testing.T) {
	options := &Options{
		Connection:                os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName:                entities.DefaultStreamName,
		ConsumerName:              entities.DefaultConsumerName,
		ConsumerSubjectIncludes:   []string{"*.>"},
		ConsumerSubjectExcludes:   []string{},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			Size:        100,
			WaitingTime: 5000,
		},
	}

	_, err := NewRetry(options)
	require.NoError(t, err)
}

func TestNewRetry_Validate(t *testing.T) {
	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
	}

	_, err := NewRetry(options)
	require.Error(t, err)
}

func TestNewVisibility(t *testing.T) {
	options := &Options{
		Connection:                os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName:                entities.DefaultStreamName,
		ConsumerName:              entities.DefaultConsumerName,
		ConsumerSubjectIncludes:   []string{"*.>"},
		ConsumerSubjectExcludes:   []string{},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			Size:        100,
			WaitingTime: 5000,
		},
	}

	_, err := NewVisibility(options)
	require.NoError(t, err)
}

func TestNewVisibility_Validate(t *testing.T) {
	options := &Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
	}

	_, err := NewVisibility(options)
	require.Error(t, err)
}
