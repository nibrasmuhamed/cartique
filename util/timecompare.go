package util

import "time"

func CompareTime(t time.Time) int64 {
	now := time.Now()
	return int64(t.Sub(now).Hours() / 24)
}
