package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/kanthorlabs/kanthorq/entities"
)

func PrinterHandler() Handler {
	return func(ctx context.Context, event *entities.Event) error {
		ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("%s | %s | %s\n", event.Id, event.Subject, ts)
		return nil
	}
}
