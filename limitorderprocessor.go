package economy

import "time"

type limitOrderProcessor struct {
	now func() time.Time
}

func (m *limitOrderProcessor) TryFillBid(
	ms MarketStorage,
	accounts Accounts,
	opl map[OrderType]orderProcessor,
	bid Bid,
) {
	for {
		if bid.Amount < 1 {
			ms.UpdateBid(bid)
			return
		}
		off, found := ms.BestOffer(bid.Symbol)
		if !found {
			ms.UpdateBid(bid)
			return
		}
		askPrice := opl[off.OfferType].GetAskingPrice(ms, off)
		if askPrice > bid.Price {
			ms.UpdateBid(bid)
			return
		}
		marketPrice := ms.LastPrice(bid.Symbol)
		var price int64
		if marketPrice >= askPrice {
			if marketPrice <= bid.Price {
				price = marketPrice
			} else {
				price = bid.Price
			}
		} else {
			price = askPrice
		}
		var filled bool
		bid, filled = fillBid(ms, accounts, m.now(), bid, off, price)
		if !filled {
			ms.UpdateBid(bid)
			return
		}
	}
}

func (m *limitOrderProcessor) TrySell(
	ms MarketStorage,
	accounts Accounts,
	opl map[OrderType]orderProcessor,
	offer Offer,
) {
	for {
		if offer.Amount < 1 {
			ms.UpdateOffer(offer)
			return
		}
		bid, found := ms.BestBid(offer.Symbol)
		if !found {
			ms.UpdateOffer(offer)
			return
		}
		price := opl[bid.BidType].GetBidPrice(ms, bid)
		if price <= offer.Price {
			bid, _ = fillBid(ms, accounts, m.now(), bid, offer, price)
			ms.UpdateBid(bid)
		} else {
			ms.UpdateOffer(offer)
			return
		}
	}
}

func (m *limitOrderProcessor) GetAskingPrice(ms MarketStorage, o Offer) int64 {
	return o.Price
}

func (m *limitOrderProcessor) GetBidPrice(ms MarketStorage, b Bid) int64 {
	return b.Price
}
