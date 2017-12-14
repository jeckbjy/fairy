package base

import (
	"fairy"
)

type BaseTransport struct {
	filters fairy.FilterChain
}

func (self *BaseTransport) New() {

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
