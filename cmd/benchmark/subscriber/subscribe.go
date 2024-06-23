package subscriber

import (
	"context"
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
			mode, err := cmd.Flags().GetString("mode")
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
					sub := instance(mode, conf)

					if err := sub.Start(ctx); err != nil {
						return err
					}

					subscribers = append(subscribers, sub)

					c := fmt.Sprintf("%s | %s -> %s -> %s", mode, conf.StreamName, conf.Topic, conf.ConsumerName)
					fmt.Printf("[%s] subscribing...\n", c)

					go sub.Consume(
						ctx,
						func(subctx context.Context, event *entities.StreamEvent) entities.JobState {
							time.Sleep(time.Millisecond * time.Duration(testify.Fake.IntBetween(100, 1000)))

							state := testify.Fake.IntBetween(int(entities.StateDiscarded)-1, int(entities.StateRetryable)+1)

							if state == int(entities.StateCancelled) {
								return entities.StateCancelled
							}
							if state == int(entities.StateRetryable) {
								return entities.StateRetryable
							}

							datac <- fmt.Sprintf("[%s] %s", c, event.EventId)
							return entities.StateCompleted
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
	command.Flags().String("mode", "available", "name subscriber mode/type you want to use")

	return command
}

func instance(mode string, conf *subscriber.Config) subscriber.Subscriber {
	if mode == "retryable" {
		return subscriber.NewRetryable(conf)
	}
	if mode == "stuck" {
		return subscriber.NewStuck(conf)
	}
	return subscriber.NewAvailable(conf)
}
