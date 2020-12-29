package utils

import "time"

func MillisecondsFrom(t time.Time) int64 {
	return time.Now().Sub(t).Milliseconds()
}
