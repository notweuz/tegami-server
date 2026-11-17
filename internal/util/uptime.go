package util

import "time"

var startUpTime = time.Now()

func GetUptime() time.Duration {
	return time.Since(startUpTime)
}