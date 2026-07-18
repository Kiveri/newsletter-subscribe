package timer

import "time"

type Timer struct{}

func NewTimer() *Timer {
	return &Timer{}
}

func (t Timer) NowUTC() time.Time {
	return time.Now().UTC()
}

func (t Timer) NowUTCString(format string) string {
	return time.Now().UTC().Format(format)
}
