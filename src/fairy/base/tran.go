package base

import (
	"fairy"
)

type TransportEx interface {
	fairy.Transport
	ConnectBy()
}

type Transport struct {
	Config
	filters fairy.FilterChain
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
