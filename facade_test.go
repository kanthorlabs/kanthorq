package kanthorq

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/kanthorlabs/kanthorq/subscriber"
	"github.com/stretchr/testify/require"
)

func TestPub(t *testing.T) {
	options := &publisher.Options{
		Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName: entities.DefaultStreamName,
	}

	_, cleanup := Pub(context.Background(), options)
	defer cleanup()
}

func TestPub_NewValidateError(t *testing.T) {
	options := &publisher.Options{}

	defer func() {
		if r := recover(); r != nil {
			require.ErrorContains(t, r.(error), "Field validation")
		}
	}()
	_, _ = Pub(context.Background(), options)
}

func TestPub_NewStartError(t *testing.T) {
	options := &publisher.Options{
		Connection: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		StreamName: entities.DefaultStreamName,
	}

	defer func() {
		if r := recover(); r != nil {
			require.ErrorContains(t, r.(error), "password authentication failed")
		}
	}()
	_, _ = Pub(context.Background(), options)
}

func TestSub(t *testing.T) {
	options := &subscriber.Options{
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

	handler := func(ctx context.Context, msg *subscriber.Message) error { return nil }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		require.ErrorIs(t, Sub(ctx, options, handler), context.DeadlineExceeded)
	}()

	<-ctx.Done()
}
