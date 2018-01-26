package util

import "time"

func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func FromSec(t int64) int64 {
	return t * 1000
}

func GetTimeByMsec(timestamp int64) time.Time {
	sec := timestamp / 1000
	nsec := (timestamp - sec*1000) * 1000
	return time.Unix(sec, nsec)
}

func HourAMPM(hour int) int {
	if hour < 1 {
		return 12
	} else if hour > 12 {
		return hour - 12
	} else {
		return hour
	}
}

func IsAM(hour int) bool {
	return hour < 12
}

func IsPM(hour int) bool {
	return hour >= 12
}
