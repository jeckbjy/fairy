package util

import (
	"time"

	"github.com/jeckbjy/fairy/util/sonyflake"
)

var sonyflakeId *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{}
	st.StartTime = time.Now()
	sonyflakeId = sonyflake.NewSonyflake(st)
}

// NextID生成一个唯一ID
func NextID() (uint64, error) {
	return sonyflakeId.NextID()
}
