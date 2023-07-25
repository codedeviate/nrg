package helper

import "time"

func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func TimeToUnix(time time.Time) int64 {
	return time.Unix()
}
