package pub

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/command"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "pub",
		Short: "publish events to KanthorQ",
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := command.GetString(cmd.Flags(), "connection-string")
			stream := command.GetString(cmd.Flags(), "stream")

			publisher, err := kanthorq.NewPublisher(uri, &kanthorq.PublisherOptions{
				StreamName: stream,
			})
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			if err := publisher.Start(ctx); err != nil {
				return err
			}
			// ignore stopping error
			defer publisher.Stop(ctx)

			subject := command.GetString(cmd.Flags(), "subject")
			body := GetBody(cmd.Flags())
			metadata := GetMetadata(cmd.Flags())

			count := command.GetInt(cmd.Flags(), "count")
			for i := 0; i < count; i++ {
				event := kanthorq.NewEvent(subject, body)
				event.Metadata.Merge(metadata)
				event.Metadata["index"] = i

				if err := publisher.Send(ctx, event); err != nil {
					log.Println(err)
					continue
				}

				ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
				fmt.Printf("%s | %s | %s\n", event.Id, event.Subject, ts)
			}
			return nil
		},
	}

	command.Flags().IntP("count", "c", 1, "number of events to publish")
	command.Flags().String("connection-string", os.Getenv("KANTHORQ_POSTGRES_URI"), "connection string of storage (PostgreSQL)")
	command.Flags().StringP("stream", "s", os.Getenv("KANTHORQ_STREAM"), "a stream name to publish event to")
	command.Flags().StringP("subject", "t", os.Getenv("KANTHORQ_SUBJECT"), "a subject name of published event")
	command.Flags().StringP("body", "b", os.Getenv("KANTHORQ_BODY"), "a body of published event")
	command.Flags().StringArray("metadata", []string{}, "a metadata of published event")

	return command
}
