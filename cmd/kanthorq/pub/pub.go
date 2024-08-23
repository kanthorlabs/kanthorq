package pub

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/entities"
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

			ch := make(chan *entities.Event, 1)
			options := &publisher.Options{
				Connection: xcmd.GetString(cmd.Flags(), "connection"),
				StreamName: xcmd.GetString(cmd.Flags(), "stream"),
			}
			go func() {
				if err := kanthorq.Pub(ctx, options, ch); err != nil {
					panic(err)
				}
			}()

			duration := xcmd.GetInt(cmd.Flags(), "duration")
			if duration > 0 {
				timeout := time.After(time.Millisecond * time.Duration(duration))
				ticker := time.Tick(time.Millisecond * time.Duration(xfaker.F.IntBetween(100, 1000)))

				for {
					select {
					case <-ctx.Done():
						return nil
					case <-timeout:
						return nil
					case <-ticker:
						events := GetEvents(cmd.Flags())
						for i := 0; i < len(events); i++ {
							ch <- events[i]
						}
					}
				}
			}

			events := GetEvents(cmd.Flags())
			for i := 0; i < len(events); i++ {
				ch <- events[i]
			}
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
