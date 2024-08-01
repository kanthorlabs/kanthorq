package pub

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/command"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "pub",
		Short: "pub is a command to publish event to KanthorQ",
		Run: func(cmd *cobra.Command, args []string) {

			uri := command.GetString(cmd.Flags(), "connection-string")
			stream := command.GetString(cmd.Flags(), "stream")

			publisher, err := kanthorq.NewPublisher(uri, &kanthorq.PublisherOptions{
				StreamName: stream,
			})
			if err != nil {
				panic(err)
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			publisher.Start(ctx)
			// ignore stopping error
			defer publisher.Stop(ctx)

			topic := command.GetString(cmd.Flags(), "topic")
			body := GetBody(cmd.Flags())
			metadata := GetMetadata(cmd.Flags())

			count := command.GetInt(cmd.Flags(), "count")
			for i := 0; i < count; i++ {
				event := kanthorq.NewEvent(topic, body)
				event.Metadata.Merge(metadata)
				event.Metadata["index"] = i

				if err := publisher.Send(ctx, event); err != nil {
					log.Println(err)
					continue
				}

				log.Println("sent:", event.Id)
			}
		},
	}

	command.Flags().IntP("count", "c", 1, "number of events to publish")
	command.Flags().String("connection-string", os.Getenv("KANTHORQ_POSTGRES_URI"), "connection string of storage (PostgreSQL)")
	command.Flags().StringP("stream", "s", os.Getenv("KANTHORQ_STREAM"), "a stream name to publish event to")
	command.Flags().StringP("topic", "t", os.Getenv("KANTHORQ_TOPIC"), "a topic name of published event")
	command.Flags().StringP("body", "b", os.Getenv("KANTHORQ_BODY"), "a body of published event")
	command.Flags().StringArray("metadata", []string{}, "a metadata of published event")

	return command
}
