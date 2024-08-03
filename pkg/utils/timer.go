package utils

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done(): //context cancelled
	case <-time.After(duration): //timeout
	}
}
