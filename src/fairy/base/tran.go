package base

import (
	"fairy"
)

type Tran struct {
	Config
	filters fairy.FilterChain
}

func (self *Tran) SetFilterChain(filters fairy.FilterChain) {
	self.filters = filters
}

func (self *Tran) GetFilterChain() fairy.FilterChain {
	return self.filters
}

func (self *Tran) AddFilters(filters ...fairy.Filter) {
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

func (self *Tran) Start() {
}
