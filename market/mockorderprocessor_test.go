package market

type mockOrderProcessor struct {
	fulfill int64
}

func (m *mockOrderProcessor) TryFillBid(bid Bid, ms MarketStorage) {
	bid.Amount -= m.fulfill
	if bid.Amount < 0 {
		bid.Amount = 0
	}
	ms.UpdateBid(bid)
}
