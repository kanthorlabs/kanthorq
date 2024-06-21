package subscriber

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/subscriber"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/spf13/cobra"
)

func Subscribe() *cobra.Command {
	command := &cobra.Command{
		Use:   "subscribe",
		Short: "subscribe messages of a topic from a consumer",
		RunE: func(cmd *cobra.Command, args []string) error {
			connection := cmd.Flags().Lookup("connection").Value.String()

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

			var datac = make(chan string)
			var subscribers []subscriber.Subscriber

			for i := range streams {
				stream := streams[i]

				for j := range topics {
					topic := topics[j]

					conf := &subscriber.Config{
						ConnectionUri: connection,
						StreamName:    stream,
						Topic:         topic,
						ConsumerName:  fmt.Sprintf("%s-%s-worker", stream, topic),
					}
					sub := subscriber.NewAvailable(conf)

					if err := sub.Start(ctx); err != nil {
						return err
					}

					subscribers = append(subscribers, sub)

					c := fmt.Sprintf("%s -> %s -> %s", conf.StreamName, conf.Topic, conf.ConsumerName)
					fmt.Printf("[%s] subscribing...\n", c)

					go sub.Consume(
						ctx,
						func(subctx context.Context, event *entities.StreamEvent) error {
							time.Sleep(time.Millisecond * time.Duration(testify.Fake.IntBetween(100, 1000)))

							if testify.Fake.IntBetween(0, 3) == 1 {
								return errors.New("aborting")
							}

							datac <- fmt.Sprintf("[%s] %s", c, event.EventId)
							return nil
						},
					)
				}
			}

			go func() {
				fmt.Println("listening for events...")
				for data := range datac {
					fmt.Println(data)
				}
			}()

			// wait for interrupt signal
			<-ctx.Done()

			fmt.Println("terminating...")
			for _, sub := range subscribers {
				if err := sub.Stop(ctx); err != nil {
					fmt.Println(err.Error())
				}
			}
			fmt.Println("terminated")

			return nil
		},
	}

	return command
}
