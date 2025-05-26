package utils

import (
	"context"
	"time"
)

// SleepContext sleeps for the provided duration but returns immediately if
// the provided is canceled before the sleep duration.
// In case we sleep for the provided duration this function returns true
// otherwise false.
func SleepContext(ctx context.Context, delay time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(delay):
		return true
	}
}
