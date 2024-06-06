package publisher

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
			connection := cmd.Flags().Lookup("connection").Value.String()

			size, err := cmd.Flags().GetInt64("size")
			if err != nil {
				return err
			}
			streams, err := cmd.Flags().GetStringSlice("streams")
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

			p := pool.New().WithErrors().WithContext(ctx)
			for i := range streams {
				stream := streams[i]

				for j := range topics {
					topic := topics[j]
					counterKey := fmt.Sprintf("%s -> %s", stream, topic)

					p.Go(func(ctx context.Context) error {
						pub := publisher.New(&publisher.Config{
							ConnectionUri: connection,
							StreamName:    stream,
						})

						if err := pub.Start(ctx); err != nil {
							return err
						}
						defer func() {
							cancelling, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
							if err := pub.Stop(cancelling); err != nil {
								log.Println(err)
							}
						}()

						for {
							if ctx.Err() != nil {
								return nil
							}

							events := testify.GenStreamEvents(ctx, topic, size)
							if err := pub.Send(ctx, events); err != nil {
								continue
							}

							total, ok := counter.Load(counterKey)
							if ok {
								total = total.(int) + len(events)
							} else {
								total = len(events)
							}
							counter.Store(counterKey, total)
							log.Printf("[%s] %d", counterKey, total)
						}
					})
				}
			}

			if err := p.Wait(); err != nil {
				log.Println(err)
			}

			return nil
		},
	}

	command.Flags().Int64("size", 100, "number of messages to send per batch")
	return command
}

func publish(topic []string) {}
