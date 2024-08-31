package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq"
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

	// simulate a transaction success
	go func() {
		// start a different connection
		conn, err := pgx.Connect(ctx, DATABASE_URI)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close(ctx)
		// initialize a transaction
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// start sending an event, it will be stored int the stream entities.DefaultStreamName
		events := []*entities.Event{
			entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")),
			entities.NewEvent("system.say_hello", []byte("{\"msg\": \"I'm comming!\"}")),
		}
		NoError(pub.SendTx(ctx, events, tx))

		// commit with done context will throw an error
		if err := tx.Commit(ctx); err != nil {
			log.Fatal(err)
		}

		log.Println("first transaction committed")
	}()

	// simulate a transaction cancelation
	go func() {
		// start a different connection
		conn, err := pgx.Connect(ctx, DATABASE_URI)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close(ctx)
		// initialize a transaction
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// start sending an event, it will be stored int the stream entities.DefaultStreamName
		events := []*entities.Event{
			entities.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")),
			entities.NewEvent("system.say_hello", []byte("{\"msg\": \"I'm comming!\"}")),
		}
		NoError(pub.SendTx(ctx, events, tx))

		<-ctx.Done()
		// commit with done context will throw an error
		if err := tx.Commit(ctx); err != nil {
			log.Println("second transaction committed error because of cancelation", err)
		}
	}()
	// wait for a while to get the goroutine started
	time.Sleep(time.Millisecond * 500)
	log.Println("press Ctrl-C to see the error and rollback the second transaction")
	<-ctx.Done()
}

func NoError(err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
