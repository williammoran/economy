package economy

import "time"

func fillBid(
	ms MarketStorage,
	accounts Accounts,
	ts time.Time,
	bid Bid,
	off Offer,
	price int64,
) (Bid, bool) {
	var amount int64
	if off.Amount <= bid.Amount {
		amount = off.Amount
	} else {
		amount = bid.Amount
	}
	totalPrice := amount * price
	if !accounts.DebitIfPossible(bid.Account, totalPrice) {
		return bid, false
	}
	accounts.Credit(off.Account, totalPrice)
	if off.Amount <= bid.Amount {
		bid.Amount -= off.Amount
		off.Amount = 0
	} else {
		off.Amount -= bid.Amount
		bid.Amount = 0
	}
	ms.NewTransaction(
		Transaction{
			BidID:   bid.ID,
			OfferID: off.ID,
			Price:   price,
			Amount:  amount,
			Date:    ts,
		},
	)
	ms.UpdateOffer(off)
	ms.SetLastPrice(off.Symbol, price)
	return bid, true
}
