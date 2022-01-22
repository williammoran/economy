package economy

type mockOrderProcessor struct {
	fulfill     int64
	askingPrice int64
}

func (m *mockOrderProcessor) TryFillBid(
	ms MarketStorage,
	accounts Accounts,
	opl map[OrderType]orderProcessor,
	bid Bid,
) {
	bid.Amount -= m.fulfill
	if bid.Amount < 0 {
		bid.Amount = 0
	}
	ms.UpdateBid(bid)
}

func (m *mockOrderProcessor) TrySell(
	ms MarketStorage,
	accounts Accounts,
	opl map[OrderType]orderProcessor,
	offer Offer,
) {
	offer.Amount -= m.fulfill
	if offer.Amount < 0 {
		offer.Amount = 0
	}
	ms.UpdateOffer(offer)
}

func (m *mockOrderProcessor) GetAskingPrice(ms MarketStorage, o Offer) int64 {
	if m.askingPrice > 0 {
		return m.askingPrice
	}
	return ms.LastPrice(o.Symbol)
}

func (m *mockOrderProcessor) GetBidPrice(ms MarketStorage, b Bid) int64 {
	return ms.LastPrice(b.Symbol)
}
