package base

import (
	"fairy"
	"fairy/util"
)

type Config struct {
	AttrMap
	CfgReconnectOpen     bool // 尝试重连次数，-1一直重连,0表示不需要断线重连
	CfgReconnectInterval int  // 单位秒
	CfgReaderBufferSize  int  // 读缓冲器大小
}

func (self *Config) SetDefaultConfig() {
	self.CfgReconnectOpen = true
	self.CfgReconnectInterval = 10
	self.CfgReaderBufferSize = 1024
}

func (self *Config) SetConfig(key *fairy.AttrKey, val interface{}) {
	switch key {
	case fairy.CfgReconnectOpen:
		if ret, err := util.ConvBool(val); err == nil {
			self.CfgReconnectOpen = ret
		}
	case fairy.CfgReconnectInterval:
		if ret, err := util.ConvInt(val); err == nil {
			self.CfgReconnectInterval = ret
		}
	case fairy.CfgReaderBufferSize:
		if ret, err := util.ConvInt(val); err == nil {
			self.CfgReaderBufferSize = ret
		}
	default:
		self.SetAttr(key, val)
	}
}

func (self *Config) GetConfig(key *fairy.AttrKey) interface{} {
	switch key {
	case fairy.CfgReconnectOpen:
		return self.CfgReconnectOpen
	case fairy.CfgReconnectInterval:
		return self.CfgReconnectInterval
	case fairy.CfgReaderBufferSize:
		return self.CfgReaderBufferSize
	default:
		return self.GetAttr(key)
	}
}
