package economy

import "time"

type marketOrderProcessor struct {
	now func() time.Time
}

func (m *marketOrderProcessor) TryFillBid(
	ms MarketStorage,
	accounts Accounts,
	opl map[OrderType]orderProcessor,
	bid Bid,
) {
	for {
		if bid.Amount < 1 {
			bid.Status = BidStatusFilled
			ms.UpdateBid(bid)
			return
		}
		off, found := ms.BestOffer(bid.Symbol)
		if !found {
			ms.UpdateBid(bid)
			return
		}
		price := opl[off.OfferType].GetAskingPrice(ms, off)
		bid = fillBid(ms, m.now(), bid, off, price)
	}
}

func (m *marketOrderProcessor) GetAskingPrice(ms MarketStorage, o Offer) int64 {
	return ms.LastPrice(o.Symbol)
}
