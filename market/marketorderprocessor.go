package market

import "time"

type marketOrderProcessor struct {
	now func() time.Time
}

func (m *marketOrderProcessor) TryFillBid(bid Bid, ms MarketStorage) {
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
		if off.Amount <= bid.Amount {
			bid.Amount -= off.Amount
			ms.NewTransaction(
				Transaction{
					BidID:   bid.BidID,
					OfferID: off.ID,
					Price:   off.Price,
					Amount:  off.Amount,
					Date:    m.now(),
				},
			)
			off.Amount = 0
		} else {
			off.Amount -= bid.Amount
			ms.NewTransaction(
				Transaction{
					BidID:   bid.BidID,
					OfferID: off.ID,
					Price:   off.Price,
					Amount:  bid.Amount,
					Date:    m.now(),
				},
			)
			bid.Amount = 0
		}
		ms.UpdateOffer(off)
		ms.SetLastPrice(off.Symbol, off.Price)
	}
}

func (m *marketOrderProcessor) GetAskingPrice(o Offer, ms MarketStorage) int64 {
	if o.Price > 0 {
		return o.Price
	}
	return ms.LastPrice(o.Symbol)
}
