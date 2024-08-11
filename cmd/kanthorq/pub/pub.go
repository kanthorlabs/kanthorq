package pub

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/command"
	"github.com/kanthorlabs/kanthorq/pkg/faker"
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

			duration := command.GetInt(cmd.Flags(), "duration")
			if duration > 0 {
				timeout := time.After(time.Millisecond * time.Duration(duration))
				ticker := time.Tick(time.Millisecond * time.Duration(faker.F.IntBetween(100, 1000)))

				for {
					select {
					case <-ctx.Done():
						return nil
					case <-timeout:
						return nil
					case <-ticker:
						events := GetEvents(cmd.Flags())
						if err := publisher.Send(ctx, events...); err != nil {
							return err
						}
					}
				}
			}

			events := GetEvents(cmd.Flags())
			return publisher.Send(ctx, events...)
		},
	}

	command.Flags().IntP("count", "c", 1, "number of events to publish")
	command.Flags().Int("duration", 0, "millisecond duration of publishing events")
	command.Flags().String("connection-string", os.Getenv("KANTHORQ_POSTGRES_URI"), "connection string of storage (PostgreSQL)")
	command.Flags().StringP("stream", "s", os.Getenv("KANTHORQ_STREAM"), "a stream name to publish event to")
	command.Flags().StringP("subject", "t", os.Getenv("KANTHORQ_SUBJECT"), "a subject name of published event")
	command.Flags().StringP("body", "b", os.Getenv("KANTHORQ_BODY"), "a body of published event")
	command.Flags().StringArray("metadata", []string{}, "a metadata of published event")

	return command
}
