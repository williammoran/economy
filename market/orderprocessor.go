package market

import "time"

func fillBid(ms MarketStorage, ts time.Time, bid Bid, off Offer, price int64) Bid {
	if off.Amount <= bid.Amount {
		bid.Amount -= off.Amount
		ms.NewTransaction(
			Transaction{
				BidID:   bid.BidID,
				OfferID: off.ID,
				Price:   price,
				Amount:  off.Amount,
				Date:    ts,
			},
		)
		off.Amount = 0
	} else {
		off.Amount -= bid.Amount
		ms.NewTransaction(
			Transaction{
				BidID:   bid.BidID,
				OfferID: off.ID,
				Price:   price,
				Amount:  bid.Amount,
				Date:    ts,
			},
		)
		bid.Amount = 0
	}
	ms.UpdateOffer(off)
	ms.SetLastPrice(off.Symbol, price)
	return bid
}
