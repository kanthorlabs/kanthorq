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
		fmt.Printf("PRINTER: %s | %s | %s\n", event.Id, event.Subject, ts)
		return nil
	}
}

func RandomErrorHandler(mod int64) Handler {
	return func(ctx context.Context, event *entities.Event) error {
		ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("RANDOM_ERROR: %s | %s | %s\n", event.Id, event.Subject, ts)

		modulus := event.CreatedAt % mod
		if modulus == 0 {
			return fmt.Errorf("random error because %d %% %d = 0", event.CreatedAt, mod)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * time.Duration(modulus)):
			return nil
		}
	}
}