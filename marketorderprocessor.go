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
			return
		}
		off, found := ms.BestOffer(bid.Symbol)
		if !found {
			return
		}
		price := opl[off.OfferType].GetAskingPrice(ms, off)
		var filled bool
		bid, _, filled = fillBid(ms, accounts, m.now(), bid, off, price)
		if !filled {
			return
		}
	}
}

func (m *marketOrderProcessor) TrySell(
	ms MarketStorage,
	accounts Accounts,
	opl map[OrderType]orderProcessor,
	offer Offer,
) {
	for {
		if offer.Amount < 1 {
			return
		}
		bid, found := ms.BestBid(offer.Symbol)
		if !found {
			return
		}
		price := opl[bid.BidType].GetBidPrice(ms, bid)
		_, offer, _ = fillBid(ms, accounts, m.now(), bid, offer, price)
	}
}

func (m *marketOrderProcessor) GetAskingPrice(ms MarketStorage, o Offer) int64 {
	return ms.LastPrice(o.Symbol)
}

func (m *marketOrderProcessor) GetBidPrice(ms MarketStorage, b Bid) int64 {
	return ms.LastPrice(b.Symbol)
}
