package economy

import (
	"testing"
	"time"
)

func TestNoTransactionWhenLimitPriceTooHigh(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{OfferType: OfferLimit, Symbol: "m", Amount: 10, Price: 20}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidLimit, Price: 10}
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusPending {
		t.Fatalf("%+v", bid)
	}
	if bid.Amount != 10 {
		t.Fatalf("%d != 10", bid.Amount)
	}
}

func TestSatisfiesAtMarketPriceInBetween(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{OfferType: OfferLimit, Symbol: "m", Amount: 10, Price: 10}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidLimit, Price: 20}
	storage.SetLastPrice("m", 15)
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
	tx := storage.transactions[0]
	if tx.Price != 15 {
		t.Fatalf("%+v", tx)
	}
}

func TestSatisfiesAtBidWhenMarketHigh(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{OfferType: OfferLimit, Symbol: "m", Amount: 10, Price: 10}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidLimit, Price: 20}
	storage.SetLastPrice("m", 25)
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
	tx := storage.transactions[0]
	if tx.Price != 20 {
		t.Fatalf("%+v", tx)
	}
}

func TestSatisfiesAtOfferWhenMarketLow(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := makeMemoryMarketStorage()
	o := Offer{OfferType: OfferLimit, Symbol: "m", Amount: 10, Price: 10}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: BidLimit, Price: 20}
	storage.SetLastPrice("m", 5)
	id := storage.AddBid(bid)
	bid.BidID = id
	mop.TryFillBid(storage, map[OrderType]OrderProcessor{OfferLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Status != BidStatusFilled {
		t.Fatalf("%+v", bid)
	}
	tx := storage.transactions[0]
	if tx.Price != 10 {
		t.Fatalf("%+v", tx)
	}
}
