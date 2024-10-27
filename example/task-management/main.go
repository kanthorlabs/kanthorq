package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
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

	event := entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}"))
	NoError(pub.Send(ctx, []*entities.Event{event}))

	// register a consumer
	conn, err := pgx.Connect(ctx, DATABASE_URI)
	NoError(err)
	registry, err := core.Do(ctx, &core.ConsumerRegisterReq{
		StreamName:                entities.DefaultStreamName,
		ConsumerName:              entities.DefaultConsumerName,
		ConsumerSubjectIncludes:   []string{"system.>"},
		ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
	}, conn)
	NoError(err)

	// convert an event to a task
	results, err := core.Do(ctx, &core.TaskConvertReq{
		Consumer:     registry.ConsumerRegistry,
		EventIds:     []string{event.Id},
		InitialState: entities.StatePending,
	}, conn)
	NoError(err)
	eventId := results.EventIds[0]
	// primary key of task is event id
	task := results.Tasks[eventId]

	// ----- MAIN ------
	cancellation, err := core.Do(ctx, &core.TaskMarkCancelledReq{
		Consumer: registry.ConsumerRegistry,
		Tasks:    []*entities.Task{task},
	}, conn)
	NoError(err)

	fmt.Printf("Canceling task %v \n", cancellation.Updated)
	fmt.Println("----- END OF EXAMPLE ------")
}

func NoError(err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
