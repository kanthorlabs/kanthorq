package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "embed"

	"github.com/kanthorlabs/kanthorq/cmd/base"
	"github.com/kanthorlabs/kanthorq/cmd/benchmark/publisher"
	"github.com/kanthorlabs/kanthorq/cmd/benchmark/subscriber"
	"github.com/kanthorlabs/kanthorq/telemetry"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := telemetry.Telemetry.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := telemetry.Telemetry.Stop(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	command := base.New()
	command.AddCommand(publisher.New())
	command.AddCommand(subscriber.New())

	command.PersistentFlags().StringP("connection", "c", os.Getenv("KANTHORQ_POSTGRES_URI"), "name of the stream")

	s := 5
	if x, err := strconv.Atoi(os.Getenv("KANTHORQ_STREAM_COUNT")); err == nil && x > 0 {
		s = x
	}
	streams := []string{}
	for i := 0; i < s; i++ {
		streams = append(streams, fmt.Sprintf("testing_%d", i))
	}
	command.PersistentFlags().StringSliceP("streams", "s", streams, "name of the stream")

	t := 5
	if x, err := strconv.Atoi(os.Getenv("KANTHORQ_TOPIC_COUNT")); err == nil && x > 0 {
		t = x
	}
	topics := []string{}
	for i := 0; i < t; i++ {
		topics = append(topics, fmt.Sprintf("benchmark.no_%d", i))
	}
	command.PersistentFlags().StringSliceP("topics", "t", topics, "name of the topic")

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Println("--- error ---")
				log.Println(err.Error())
			}
		}
	}()

	if err := command.Execute(); err != nil {
		panic(err)
	}
}
