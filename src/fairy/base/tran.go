package base

import (
	"fairy"
	"fairy/util"
)

type Config struct {
	AttrMap
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

type Transport struct {
	Config
	filters fairy.FilterChain
}

func (self *Transport) NewBase() {
	self.SetDefaultConfig()
}

func (self *Transport) SetConfig(key *fairy.AttrKey, val interface{}) {
	switch key {
	case fairy.KeyReconnectInterval:
		if ret, err := util.ConvInt(val); err == nil {
			self.CfgReconnectInterval = ret
		}
	case fairy.KeyReaderBufferSize:
		if ret, err := util.ConvInt(val); err == nil {
			self.CfgReaderBufferSize = ret
		}
	default:
		self.SetAttr(key, val)
	}
}

func (self *Transport) GetConfig(key *fairy.AttrKey) interface{} {
	switch key {
	case fairy.KeyReconnectInterval:
		return self.CfgReconnectInterval
	case fairy.KeyReaderBufferSize:
		return self.CfgReaderBufferSize
	default:
		return self.GetAttr(key)
	}
}

func (self *Transport) SetFilterChain(filters fairy.FilterChain) {
	self.filters = filters
}

func (self *Transport) GetFilterChain() fairy.FilterChain {
	return self.filters
}

func (self *Transport) AddFilters(filters ...fairy.Filter) {
	if self.filters == nil {
		self.filters = NewFilterChain()
	}

	// auto add filter
	if self.filters.Len() == 0 {
		if _, ok := filters[0].(*TransferFilter); !ok {
			self.filters.AddLast(NewTransferFilter())
		}
	}

	for _, filter := range filters {
		self.filters.AddLast(filter)
	}
}

func (self *Transport) Start() {
}
