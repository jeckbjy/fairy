package base

import (
	"fairy"
	"fairy/util"
	"math"
)

type TransportEx interface {
	fairy.Transport
	ConnectBy()
}

type Transport struct {
	Config
	filters fairy.FilterChain
}

func (self *Transport) SetConfig(key *fairy.AttrKey, val interface{}) {
	switch key {
	case fairy.CfgReconnectCount:
		if ret, err := util.ConvInt(val); err == nil {
			self.CfgReconnectCount = ret
			if self.CfgReconnectCount < 0 {
				self.CfgReconnectCount = math.MaxInt32
			}
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

func (self *Transport) GetConfig(key *fairy.AttrKey) interface{} {
	switch key {
	case fairy.CfgReconnectCount:
		return self.CfgReconnectCount
	case fairy.CfgReconnectInterval:
		return self.CfgReconnectInterval
	case fairy.CfgReaderBufferSize:
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
