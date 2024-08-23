package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/subscriber"
)

func main() {
	var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"

	// Initialize a publisher
	pub, _ := publisher.New(DATABASE_URI, &publisher.Options{
		StreamName: entities.DefaultStreamName,
	})
	NoError(pub.Start(context.Background()))
	defer pub.Stop(context.Background())

	// publish some events, it will be stored inside our stream with name kanthorq.DefaultStreamName
	NoError(pub.Send(
		context.Background(),
		entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")),
	))

	// Initialize a subscriber that will process events that has subject that match with the filter "system.>"
	// so both system.say_hello and system.say_goodbye will be processed
	sub, _ := subscriber.New(DATABASE_URI, &subscriber.Options{
		StreamName:            entities.DefaultStreamName,
		ConsumerName:          entities.DefaultConsumerName,
		ConsumerSubjectFilter: []string{"system.>"},
		ConsumerAttemptMax:    entities.DefaultConsumerAttemptMax,
	})
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	NoError(sub.Start(ctx))
	defer sub.Stop(ctx)

	go sub.Receive(ctx, func(ctx context.Context, event *entities.Event) error {
		log.Print(string(event.Body))
		return nil
	})

	<-ctx.Done()
}

func NoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
