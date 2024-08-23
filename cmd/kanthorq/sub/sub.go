package sub

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/xcmd"
	"github.com/kanthorlabs/kanthorq/puller"
	"github.com/kanthorlabs/kanthorq/subscriber"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "sub",
		Short: "subscribe task from KanthorQ",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			options := &subscriber.Options{
				Connection:            xcmd.GetString(cmd.Flags(), "connection"),
				StreamName:            xcmd.GetString(cmd.Flags(), "stream"),
				ConsumerName:          xcmd.GetString(cmd.Flags(), "consumer"),
				ConsumerSubjectFilter: xcmd.GetStringSlice(cmd.Flags(), "subjects"),
				ConsumerAttemptMax:    1,
				Puller: &puller.PullerIn{
					Size:        100,
					WaitingTime: 5000,
				},
			}

			err := kanthorq.Sub(ctx, options, subscriber.PrinterHandler())
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}

			return nil
		},
	}

	command.Flags().String("connection", os.Getenv("KANTHORQ_POSTGRES_URI"), "connection string of storage (PostgreSQL)")
	command.Flags().StringP("stream", "s", os.Getenv("KANTHORQ_STREAM"), "a stream name we want to subscribe events from")
	command.Flags().StringP("consumer", "c", os.Getenv("KANTHORQ_CONSUMER"), "a consumer name")
	command.Flags().StringSliceP("subjects", "t", []string{os.Getenv("KANTHORQ_SUBJECT")}, "a subject name we want to subscribe")

	return command
}
