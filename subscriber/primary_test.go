package subscriber

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/stretchr/testify/require"
)

func TestPrimary_Receive(t *testing.T) {
	options := &Options{
		Connection:                os.Getenv("KANTHORQ_POSTGRES_URI"),
		StreamName:                xfaker.StreamName(),
		ConsumerName:              xfaker.ConsumerName(),
		ConsumerSubjectFilter:     []string{"system.>"},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			Size:        100,
			WaitingTime: 1000,
		},
	}

	var bodyc = make(chan *entities.Event)
	go func() {
		sub, err := New(options)
		require.NoError(t, err)

		sub.Start(context.Background())
		defer func() {
			require.NoError(t, sub.Stop(context.Background()))
		}()

		err = sub.Receive(context.Background(), func(ctx context.Context, msg *Message) error {
			if msg.Event.Subject == "system.test.panic" {
				panic(errors.New(string(msg.Event.Body)))
			}
			if msg.Event.Subject == "system.test.error" {
				return errors.New(string(msg.Event.Body))
			}
			bodyc <- msg.Event
			return nil
		})
		require.NoError(t, err)
	}()

	go func() {
		pub, err := publisher.New(
			&publisher.Options{
				Connection: os.Getenv("KANTHORQ_POSTGRES_URI"),
				StreamName: options.StreamName,
			},
		)
		require.NoError(t, err)

		require.NoError(t, pub.Start(context.Background()))
		defer func() {
			require.NoError(t, pub.Stop(context.Background()))
		}()

		// first message, it's ok
		err = pub.Send(context.Background(), []*entities.Event{
			entities.NewEvent("system.test.ok", []byte("ok")),
		})
		require.NoError(t, err)
		<-time.After(time.Millisecond * time.Duration(2*options.Puller.WaitingTime))

		// second message, it's panic
		err = pub.Send(context.Background(), []*entities.Event{
			entities.NewEvent("system.test.panic", []byte("panic")),
		})
		require.NoError(t, err)
		<-time.After(time.Millisecond * time.Duration(2*options.Puller.WaitingTime))

		// third message, it's error
		err = pub.Send(context.Background(), []*entities.Event{
			entities.NewEvent("system.test.error", []byte("error")),
		})
		require.NoError(t, err)

		// final message, it's done
		err = pub.Send(context.Background(), []*entities.Event{
			entities.NewEvent("system.test.done", []byte("done")),
		})
		require.NoError(t, err)
	}()

	for e := range bodyc {
		if e.Subject == "system.test.done" {
			// wait for a while to let message to be acked
			time.Sleep(time.Millisecond * time.Duration(2*options.Puller.WaitingTime))
			return
		}
	}
}
