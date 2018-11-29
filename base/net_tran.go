package base

import (
	"github.com/jeckbjy/fairy"
)

type Tran struct {
	chain fairy.IFilterChain
}

func (tran *Tran) GetChain() fairy.IFilterChain {
	return tran.chain
}

func (tran *Tran) SetOptions(...fairy.Option) {
}

func (tran *Tran) AddFilters(filters ...fairy.IFilter) {
	if tran.chain == nil {
		tran.chain = NewChain()
	}

	// Auto Add TransferFilter
	if tran.chain.Len() == 0 {
		if _, ok := filters[0].(*TransferFilter); !ok {
			tran.chain.AddLast(NewTransferFilter())
		}
	}

	tran.chain.AddLast(filters...)
}

func (tran *Tran) Start() {
}

func (tran *Tran) Stop() {
}
