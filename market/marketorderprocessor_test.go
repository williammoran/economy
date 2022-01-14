package market

import (
	"testing"
	"time"
)

func TestTryFillFilledByExactOffer(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 10, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(bid, storage)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
}
