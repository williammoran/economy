package economy

import (
	"testing"
	"time"
)

func TestFillBidFilledByExactOffer(t *testing.T) {
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 10}
	offerID := storage.AddOffer(o)
	o.ID = offerID
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	bid = fillBid(storage, time.Time{}, bid, o, 7)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	o = storage.offers["m"][offerID]
	if o.Amount != 0 {
		t.Fatalf("%+v", o)
	}
	if storage.LastPrice("m") != 7 {
		t.Fatalf("%d != 7", storage.LastPrice("m"))
	}
}

func TestFillBidFilledByLargerOffer(t *testing.T) {
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 20}
	offerID := storage.AddOffer(o)
	o.ID = offerID
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	bid = fillBid(storage, time.Time{}, bid, o, 7)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	o = storage.offers["m"][offerID]
	if o.Amount != 10 {
		t.Fatalf("%+v", o)
	}
	if storage.LastPrice("m") != 7 {
		t.Fatalf("%d != 7", storage.LastPrice("m"))
	}
}

func TestFillBidPartiallyFilled(t *testing.T) {
	storage := makeMemoryMarketStorage()
	o := Offer{Symbol: "m", Amount: 5}
	offerID := storage.AddOffer(o)
	o.ID = offerID
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	bid = fillBid(storage, time.Time{}, bid, o, 7)
	if bid.Amount != 5 {
		t.Fatalf("%+v", bid)
	}
	o = storage.offers["m"][offerID]
	if o.Amount != 0 {
		t.Fatalf("%+v", o)
	}
	if storage.LastPrice("m") != 7 {
		t.Fatalf("%d != 7", storage.LastPrice("m"))
	}
}
