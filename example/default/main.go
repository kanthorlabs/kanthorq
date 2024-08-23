package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanthorlabs/kanthorq"
)

func main() {
	var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"

	// Initialize a publisher
	publisher, _ := kanthorq.NewPublisher(DATABASE_URI, &kanthorq.PublisherOptions{
		StreamName: kanthorq.DefaultStreamName,
	})
	kanthorq.NoError(publisher.Start(context.Background()))
	defer publisher.Stop(context.Background())

	// publish some events, it will be stored inside our stream with name kanthorq.DefaultStreamName
	kanthorq.NoError(publisher.Send(
		context.Background(),
		kanthorq.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")),
	))

	// Initialize a subscriber that will process events that has subject that match with the filter "system.>"
	// so both system.say_hello and system.say_goodbye will be processed
	subscriber, _ := kanthorq.NewSubscriber(DATABASE_URI, &kanthorq.SubscriberOptions{
		StreamName:            kanthorq.DefaultStreamName,
		ConsumerName:          kanthorq.DefaultConsumerName,
		ConsumerSubjectFilter: []string{"system.>"},
		ConsumerAttemptMax:    kanthorq.DefaultConsumerAttemptMax,
	})
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	kanthorq.NoError(subscriber.Start(ctx))
	defer subscriber.Stop(ctx)

	go subscriber.Receive(ctx, func(ctx context.Context, event *kanthorq.Event) error {
		log.Print(string(event.Body))
		return nil
	})

	<-ctx.Done()
}
