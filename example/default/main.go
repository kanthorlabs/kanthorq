package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// start sending an event, it will be stored int the stream entities.DefaultStreamName
	events := []*entities.Event{
		entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")),
		entities.NewEvent("system.say_hello", []byte("{\"msg\": \"I'm comming!\"}")),
	}
	NoError(pub.Send(ctx, events))

	go func() {
		// wait a few seconds to send another event
		time.Sleep(time.Second * 3)

		// start sending another event
		events := []*entities.Event{
			entities.NewEvent("system.say_goodbye", []byte("{\"msg\": \"See you!\"}")),
		}
		NoError(pub.Send(ctx, events))
	}()

	// Initialize a subscriber that will process events that has subject that match with the filter "system.>"
	// so both system.say_hello and system.say_goodbye will be processed
	err := kanthorq.Sub(ctx, &subscriber.Options{
		Connection:                DATABASE_URI,
		StreamName:                entities.DefaultStreamName,
		ConsumerName:              entities.DefaultConsumerName,
		ConsumerSubjectFilter:     []string{"system.>"},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: &puller.PullerIn{
			// Size is how many events you want to pull at one batch
			Size: 100,
			// WaitingTime is how long you want to wait before pulling again
			// if you didn't get enough events in current batch
			WaitingTime: 1000,
		},
	},
		func(ctx context.Context, msg *subscriber.Message) error {
			ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
			// print out recevied event
			log.Printf("RECEIVED: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)
			return nil
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
