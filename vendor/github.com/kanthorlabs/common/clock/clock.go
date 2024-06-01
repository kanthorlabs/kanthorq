package clock

import "time"

func New() Clock {
	return &clock{}
}

type Clock interface {
	Now() time.Time
	UnixMilli(msec int64) time.Time
}

type clock struct {
}

// Now return current UTC time
func (t *clock) Now() time.Time {
	return time.Now().UTC()
}

// Now return current UTC time from milliseconds
func (t *clock) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec).UTC()
}
