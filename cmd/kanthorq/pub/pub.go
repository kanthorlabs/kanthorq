package pub

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/xcmd"
	"github.com/kanthorlabs/kanthorq/pkg/xfaker"
	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "pub",
		Short: "publish events to KanthorQ",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			options := &publisher.Options{
				Connection: xcmd.GetString(cmd.Flags(), "connection"),
				StreamName: xcmd.GetString(cmd.Flags(), "stream"),
			}

			pub, cleanup := kanthorq.Pub(ctx, options)
			defer cleanup()

			var count int

			duration := xcmd.GetInt(cmd.Flags(), "duration")
			if duration > 0 {
				timeout := time.After(time.Millisecond * time.Duration(duration))
				ticker := time.Tick(time.Millisecond * time.Duration(xfaker.F.IntBetween(100, 1000)))
				for {
					select {
					case <-ctx.Done():
						fmt.Printf("published %d events\n", count)
						return nil
					case <-timeout:
						fmt.Printf("--------- %d events ---------\n", count)
						return nil
					case <-ticker:
						events := GetEvents(cmd.Flags())
						if err := pub.Send(ctx, events); err != nil {
							return err
						}
						count += len(events)
					}
				}
			}

			events := GetEvents(cmd.Flags())
			count += len(events)
			if err := pub.Send(ctx, events); err != nil {
				return err
			}

			fmt.Printf("--------- %d events ---------\n", count)
			return nil
		},
	}

	command.Flags().IntP("count", "c", 1, "number of events to publish")
	command.Flags().Int("duration", 0, "millisecond duration of publishing events")
	command.Flags().String("connection", os.Getenv("KANTHORQ_POSTGRES_URI"), "connection string of storage (PostgreSQL)")
	command.Flags().StringP("stream", "s", os.Getenv("KANTHORQ_STREAM"), "a stream name to publish event to")
	command.Flags().StringP("subject", "t", os.Getenv("KANTHORQ_SUBJECT"), "a subject name of published event")
	command.Flags().StringP("body", "b", os.Getenv("KANTHORQ_BODY"), "a body of published event")
	command.Flags().StringArray("metadata", []string{}, "a metadata of published event")

	return command
}
