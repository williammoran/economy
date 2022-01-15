package market

type mockOrderProcessor struct {
	fulfill     int64
	askingPrice int64
}

func (m *mockOrderProcessor) TryFillBid(bid Bid, ms MarketStorage) {
	bid.Amount -= m.fulfill
	if bid.Amount < 0 {
		bid.Amount = 0
	}
	ms.UpdateBid(bid)
}

func (m *mockOrderProcessor) GetAskingPrice(o Offer, ms MarketStorage) int64 {
	if m.askingPrice > 0 {
		return m.askingPrice
	}
	return ms.LastPrice(o.Symbol)
}
