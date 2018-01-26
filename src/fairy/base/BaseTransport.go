package base

import (
	"fairy"
	"fairy/util"
)

type Config struct {
	BaseAttrMap
	CfgReconnectInterval int // 单位秒
	CfgReaderBufferSize  int
}

func (self *Config) SetDefaultConfig() {
	self.CfgReconnectInterval = 10
	self.CfgReaderBufferSize = 1024
}

func (self *Config) IsNeedReconnect() bool {
	return self.CfgReconnectInterval >= 0
}

type BaseTransport struct {
	Config
	filters fairy.FilterChain
}

func (self *BaseTransport) NewBase() {
	self.SetDefaultConfig()
}

func (self *BaseTransport) SetConfig(key *fairy.AttrKey, val string) {
	switch key {
	case fairy.KeyReconnectInterval:
		util.SafeParseInt(&self.CfgReconnectInterval, val)
	case fairy.KeyReaderBufferSize:
		util.SafeParseInt(&self.CfgReaderBufferSize, val)
	default:
		self.SetAttr(key, val)
	}
}

func (self *BaseTransport) GetConfig(key *fairy.AttrKey) interface{} {
	switch key {
	case fairy.KeyReconnectInterval:
		return self.CfgReconnectInterval
	case fairy.KeyReaderBufferSize:
		return self.CfgReaderBufferSize
	default:
		return self.GetAttr(key)
	}
}

func (self *BaseTransport) SetFilterChain(filters fairy.FilterChain) {
	self.filters = filters
}

func (self *BaseTransport) GetFilterChain() fairy.FilterChain {
	return self.filters
}

func (self *BaseTransport) AddFilters(filters ...fairy.Filter) {
	if self.filters == nil {
		self.filters = NewFilterChain()
	}

	for _, filter := range filters {
		self.filters.AddLast(filter)
	}
}
