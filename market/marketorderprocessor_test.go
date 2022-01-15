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
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
	if storage.LastPrice("m") != 20 {
		t.Fatalf("%d != 20", storage.LastPrice("m"))
	}
}

func TestTryFillFilledByLargerOffer(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 20, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
	for _, o := range storage.offers["m"] {
		if o.Amount != 10 {
			t.Fatalf("%+v", o)
		}
	}
	if storage.LastPrice("m") != 20 {
		t.Fatalf("%d != 20", storage.LastPrice("m"))
	}
}

func TestTryFillPartiallyFilled(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 5, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusPending {
		t.Fatalf("%+v", bid)
	}
	if bid.Amount != 5 {
		t.Fatalf("%+v", bid)
	}
	if storage.LastPrice("m") != 20 {
		t.Fatalf("%d != 20", storage.LastPrice("m"))
	}
}

func TestTryFillFilledBy2Offers(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 8, Price: 20}
	storage.AddOffer(o)
	o = Offer{Symbol: "m", Amount: 8, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
	if len(storage.offers["m"]) > 1 {
		t.Fatalf("%+v", storage.offers["m"])
	}
	var total int64
	for _, o := range storage.offers["m"] {
		total += o.Amount
	}
	if total != 6 {
		t.Fatalf("%d remaining", total)
	}
}
