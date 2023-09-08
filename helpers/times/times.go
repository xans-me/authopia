package times

import (
	"time"
)

// Now is a func to get current time string
func Now(hourOffset int, format string) string {
	return time.Now().UTC().Add(time.Hour * time.Duration(hourOffset)).Format(format)
}
