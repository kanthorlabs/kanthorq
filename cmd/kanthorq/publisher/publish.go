package publisher

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kanthorlabs/kanthorq/publisher"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
)

func Publish() *cobra.Command {
	command := &cobra.Command{
		Use:   "publish",
		Short: "publish messages to a message stream",
		RunE: func(cmd *cobra.Command, args []string) error {
			pub := publisher.New(&publisher.Config{
				ConnectionUri: cmd.Flags().Lookup("connection").Value.String(),
				StreamName:    cmd.Flags().Lookup("stream").Value.String(),
			})

			ctx := cmd.Context()
			if err := pub.Start(ctx); err != nil {
				return err
			}
			defer pub.Stop(ctx)

			count, err := cmd.Flags().GetInt64("count")
			if err != nil {
				return err
			}
			topics, err := cmd.Flags().GetStringSlice("topics")
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			var counter = sync.Map{}

			p := pool.New().WithContext(ctx).WithMaxGoroutines(len(topics))
			for i := range topics {
				topic := topics[i]
				p.Go(func(ctx context.Context) error {
					select {
					case <-ctx.Done():
						return nil
					default:
						for {
							if ctx.Err() != nil {
								return nil
							}

							events := testify.GenStreamEvents(ctx, topic, count)
							if err := pub.Send(ctx, events); err != nil {
								log.Println(err)
								continue
							}

							total, ok := counter.Load(topic)
							if ok {
								total = total.(int) + len(events)
							} else {
								total = len(events)
							}
							counter.Store(topic, total)
							log.Printf("[%s] %d", topic, total)
						}
					}
				})
			}

			return p.Wait()
		},
	}

	command.Flags().Int64("count", 100, "number of messages to send")
	return command
}

func publish(topic []string) {}
