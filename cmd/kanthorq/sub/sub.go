package sub

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/pkg/command"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "sub",
		Short: "subscribe task from KanthorQ",
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := command.GetString(cmd.Flags(), "connection-string")
			stream := command.GetString(cmd.Flags(), "stream")
			subject := command.GetString(cmd.Flags(), "subject")
			consumer := command.GetString(cmd.Flags(), "consumer")

			subscriber, err := kanthorq.NewSubscriber(uri, &kanthorq.SubscriberOptions{
				StreamName:         stream,
				ConsumerName:       consumer,
				ConsumerSubject:    subject,
				ConsumerAttemptMax: 1,
				HandleTimeout:      5000,
			})
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			if err := subscriber.Start(ctx); err != nil {
				return err
			}
			// ignore stopping error
			defer subscriber.Stop(ctx)

			if err := subscriber.Receive(ctx, kanthorq.SubscriberHandlerPrinter()); !errors.Is(err, context.Canceled) {
				return err
			}

			return nil
		},
	}

	command.Flags().String("connection-string", os.Getenv("KANTHORQ_POSTGRES_URI"), "connection string of storage (PostgreSQL)")
	command.Flags().StringP("stream", "s", os.Getenv("KANTHORQ_STREAM"), "a stream name we want to subscribe events from")
	command.Flags().StringP("subject", "t", os.Getenv("KANTHORQ_SUBJECT"), "a subject name we want to subscribe")
	command.Flags().StringP("consumer", "c", os.Getenv("KANTHORQ_CONSUMER"), "a consumer name")

	return command
}
