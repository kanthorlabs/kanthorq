package main

import (
	"log"
	"os"

	_ "embed"

	"github.com/kanthorlabs/kanthorq/cmd/base"
	"github.com/kanthorlabs/kanthorq/cmd/kanthorq/publisher"
)

func main() {
	command := base.New()
	command.AddCommand(publisher.New())

	command.PersistentFlags().StringP("connection", "c", os.Getenv("KANTHORQ_POSTGRES_URI"), "name of the stream")
	command.PersistentFlags().StringP("stream", "s", "testing", "name of the stream")

	topics := []string{"testing.benchmark_1", "testing.benchmark_2", "testing.benchmark_3"}
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
