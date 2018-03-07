package base

import "math"

type Config struct {
	AttrMap
	CfgReconnectCount    int // 尝试重连次数，-1一直重连,0表示不需要断线重连
	CfgReconnectInterval int // 单位秒
	CfgReaderBufferSize  int // 读缓冲器大小
}

func (self *Config) SetDefaultConfig() {
	self.CfgReconnectCount = math.MaxInt32
	self.CfgReconnectInterval = 10
	self.CfgReaderBufferSize = 1024
}

func (self *Config) IsNeedReconnect() bool {
	return self.CfgReconnectInterval >= 0
}
