package market

import "time"

type limitOrderProcessor struct {
	now func() time.Time
}

func (m *limitOrderProcessor) TryFillBid(ms MarketStorage, opl map[OrderType]OrderProcessor, bid Bid) {
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
		bid = fillBid(ms, m.now(), bid, off, price)
	}
}

func (m *limitOrderProcessor) GetAskingPrice(ms MarketStorage, o Offer) int64 {
	return o.Price
}
