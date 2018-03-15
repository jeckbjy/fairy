package test

import (
	"fairy/container/zset"
	"fairy/util"
	"fmt"
)

func TestZSet() {
	zs := zset.New(true)
	max := 50000

	beg := util.Now()
	for i := int64(1); i <= int64(max); i++ {
		zs.Insert(util.ConvStr(i), i)
	}
	end := util.Now()

	fmt.Printf("init use time ms:%+v, len=%v\n", end-beg, zs.Len())

	beg = util.Now()
	zs.Insert("1000000", 1000000)
	end = util.Now()
	fmt.Printf("one  use time ms:%+v\n", end-beg)

	beg = util.Now()
	zs.Scan(0, 0, func(rank uint64, el *zset.Element) {
		if rank != uint64(el.Score) {
			fmt.Printf("bad rank:%+v", rank)
		}
	})
	end = util.Now()
	fmt.Printf("scan use time ms:%+v\n", end-beg)
}
