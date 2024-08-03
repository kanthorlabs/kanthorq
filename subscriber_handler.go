package kanthorq

import (
	"context"
	"fmt"
	"time"
)

func SubscriberHandlerPrinter() SubscriberHandler {
	return func(ctx context.Context, event *Event) error {
		ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("%s | %s | %s\n", event.Id, event.Subject, ts)
		return nil
	}
}
