package market

func MakeMarket(s MarketStorage) *Market {
	return &Market{
		storage: s,
	}
}

type Market struct {
	storage MarketStorage
}

func (m *Market) Offer(o Offer) {
	m.storage.AddOffer(o)
}

func (m *Market) Bid(b Bid) BidID {
	return m.storage.AddBid(b)
}
