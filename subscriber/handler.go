package subscriber

import (
	"context"
	"fmt"
	"time"
)

func PrinterHandler() Handler {
	return func(ctx context.Context, msg *Message) error {
		ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("PRINTER: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)
		return nil
	}
}

func RandomErrorHandler(mod int64) Handler {
	return func(ctx context.Context, msg *Message) error {
		ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
		fmt.Printf("RANDOM_ERROR: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)

		modulus := msg.Event.CreatedAt % mod
		if modulus == 0 {
			return fmt.Errorf("random error because %d %% %d = 0", msg.Event.CreatedAt, mod)
		}

		if modulus == 1 {
			panic(fmt.Sprintf("random error because %d %% %d = 1", msg.Event.CreatedAt, mod))
		}

		// simulate execution time
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * time.Duration(modulus)):
			return nil
		}
	}
}
