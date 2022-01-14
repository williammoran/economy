package market

import "sync"

func MakeMarket(s MarketStorage) *Market {
	return &Market{
		storage: s,
	}
}

type Market struct {
	sync.Mutex
	storage MarketStorage
}

func (m *Market) Offer(o Offer) {
	m.storage.AddOffer(o)
}
