package economy

import (
	"testing"
	"time"
)

func TestNoTransactionWhenLimitPriceTooHigh(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{OfferType: OrderTypeLimit, Symbol: "m", Amount: 10, Price: 20}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 10}
	id := storage.AddBid(bid)
	bid.ID = id
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount == 0 {
		t.Fatalf("%+v", bid)
	}
	if bid.Amount != 10 {
		t.Fatalf("%d != 10", bid.Amount)
	}
}

func TestSatisfiesAtMarketPriceInBetween(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{OfferType: OrderTypeLimit, Symbol: "m", Amount: 10, Price: 10}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 20}
	storage.SetLastPrice("m", 15)
	id := storage.AddBid(bid)
	bid.ID = id
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	tx := storage.transactions[0]
	if tx.Price != 15 {
		t.Fatalf("%+v", tx)
	}
}

func TestSatisfiesAtBidWhenMarketHigh(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{OfferType: OrderTypeLimit, Symbol: "m", Amount: 10, Price: 10}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 20}
	storage.SetLastPrice("m", 25)
	id := storage.AddBid(bid)
	bid.ID = id
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	tx := storage.transactions[0]
	if tx.Price != 20 {
		t.Fatalf("%+v", tx)
	}
}

func TestSatisfiesAtOfferWhenMarketLow(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{OfferType: OrderTypeLimit, Symbol: "m", Amount: 10, Price: 10}
	storage.AddOffer(o)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 20}
	storage.SetLastPrice("m", 5)
	id := storage.AddBid(bid)
	bid.ID = id
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	tx := storage.transactions[0]
	if tx.Price != 10 {
		t.Fatalf("%+v", tx)
	}
}

func TestLimitTrySellFillsExactMatch(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 5}
	bid.ID = storage.AddBid(bid)
	offer := Offer{Symbol: "m", Amount: 10, OfferType: OrderTypeLimit, Price: 5}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, offer)
	bid = storage.GetBid(bid.ID)
	if bid.IsActive() {
		t.Fatalf("Bid still active: %+v", bid)
	}
	offer = storage.GetOffer(offer.ID)
	if offer.IsActive() {
		t.Fatalf("Offer still active: %+v", offer)
	}
}

func TestLimitTrySellNoBidsCompletes(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	offer := Offer{Symbol: "m", Amount: 10, OfferType: OrderTypeLimit}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, offer)
	offer = storage.GetOffer(offer.ID)
	if !offer.IsActive() {
		t.Fatal("Offer not active")
	}
	if offer.Amount != 10 {
		t.Fatalf("%+v", offer)
	}
}

func TestLimitTrySell2BidsSell(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	bid0 := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 1}
	bid0.ID = storage.AddBid(bid0)
	bid1 := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeLimit, Price: 1}
	bid1.ID = storage.AddBid(bid1)
	offer := Offer{Symbol: "m", Amount: 20, OfferType: OrderTypeLimit, Price: 1}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, offer)
	bid0 = storage.GetBid(bid0.ID)
	if bid0.IsActive() {
		t.Fatalf("Bid0 still active: %+v", bid0)
	}
	bid1 = storage.GetBid(bid1.ID)
	if bid1.IsActive() {
		t.Fatal("Bid1 still active")
	}
	offer = storage.GetOffer(offer.ID)
	if offer.IsActive() {
		t.Fatalf("Offer still active: %+v", offer)
	}
}

func TestLimitTrySellPartialBidCompletesAndDecrimentsOffer(t *testing.T) {
	mop := limitOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	bid := Bid{Symbol: "m", Amount: 5, BidType: OrderTypeLimit, Price: 1}
	bid.ID = storage.AddBid(bid)
	offer := Offer{Symbol: "m", Amount: 10, OfferType: OrderTypeLimit, Price: 1}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeLimit: &mop}, offer)
	bid = storage.GetBid(bid.ID)
	if bid.IsActive() {
		t.Fatalf("Bid still active: %+v", bid)
	}
	offer = storage.GetOffer(offer.ID)
	if !offer.IsActive() {
		t.Fatal("Offer inactive")
	}
	if offer.Amount != 5 {
		t.Fatalf("Wrong remaining amount: %+v", offer)
	}
}
