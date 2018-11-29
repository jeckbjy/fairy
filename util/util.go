package util

import "time"

// Now return millisecond timestamp
func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
