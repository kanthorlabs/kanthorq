package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/kanthorlabs/kanthorq/subscriber"
)

func main() {
	var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"
	if uri := os.Getenv("KANTHORQ_POSTGRES_URI"); uri != "" {
		DATABASE_URI = uri
	}

	// listen for SIGINT and SIGTERM so if you press Ctrl-C you can stop the program
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize a publisher
	pub, cleanup := kanthorq.Pub(ctx, &publisher.Options{
		Connection: DATABASE_URI,
		StreamName: entities.DefaultStreamName,
	})
	defer cleanup()

	events := []*entities.Event{
		entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")),
		entities.NewEvent("system.say_goodbye", []byte("{\"msg\": \"I'm comming!\"}")),
	}
	NoError(pub.Send(ctx, events))

	err := kanthorq.Sub(ctx, &subscriber.Options{
		Connection:                DATABASE_URI,
		StreamName:                entities.DefaultStreamName,
		ConsumerName:              entities.DefaultConsumerName,
		ConsumerSubjectIncludes:   []string{"system.>"},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			// Size is how many events you want to pull at one batch
			Size: 100,
			// WaitingTime is how long you want to wait before pulling again
			// if you didn't get enough events in current batch
			WaitingTime: 1000,
		},
	},
		func(ctx context.Context, msg *subscriber.Message) error {
			// acknowledge the event
			if msg.Event.Subject == "system.say_hello" {
				return msg.Ack(ctx)
			}
			// we don't want to say goodby, not acknowledging it
			return msg.Nack(ctx, errors.New("not saying goodbye"))
		},
	)
	NoError(err)
	fmt.Println("----- END OF EXAMPLE ------")
}

func NoError(err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
