package util

import "time"

func Now() int64 {
	return time.Now().UnixNano() * int64(time.Microsecond)
}
