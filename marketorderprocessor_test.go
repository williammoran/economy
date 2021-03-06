package economy

import (
	"testing"
	"time"
)

func TestTryFillFilledByExactOffer(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{Symbol: "m", Amount: 10, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	storage.SetLastPrice("m", 7)
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	if storage.LastPrice("m") != 7 {
		t.Fatalf("%d != 7", storage.LastPrice("m"))
	}
}

func TestTryFillFilledByLargerOffer(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{Symbol: "m", Amount: 20, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	storage.SetLastPrice("m", 7)
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	for _, o := range storage.offers["m"] {
		if o.Amount != 10 {
			t.Fatalf("%+v", o)
		}
	}
	if storage.LastPrice("m") != 7 {
		t.Fatalf("%d != 7", storage.LastPrice("m"))
	}
}

func TestTryFillPartiallyFilled(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{Symbol: "m", Amount: 5, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	storage.SetLastPrice("m", 7)
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 5 {
		t.Fatalf("%+v", bid)
	}
	if bid.Amount != 5 {
		t.Fatalf("%+v", bid)
	}
	if storage.LastPrice("m") != 7 {
		t.Fatalf("%d != 7", storage.LastPrice("m"))
	}
}

func TestTryFillFilledBy2Offers(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	o := Offer{Symbol: "m", Amount: 8, Price: 20}
	storage.AddOffer(o)
	o = Offer{Symbol: "m", Amount: 8, Price: 20}
	storage.AddOffer(o)
	t.Logf("%+v", storage.offers)
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	mop.TryFillBid(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, bid)
	bid = storage.GetBid(id)
	if bid.Amount != 0 {
		t.Fatalf("%+v", bid)
	}
	if len(storage.offers["m"]) != 2 {
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

func TestTrySellFillsExactMatch(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	bid.ID = storage.AddBid(bid)
	offer := Offer{Symbol: "m", Amount: 10, OfferType: OrderTypeMarket}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, offer)
	bid = storage.GetBid(bid.ID)
	if bid.IsActive() {
		t.Fatal("Bid still active")
	}
	offer = storage.GetOffer(offer.ID)
	if offer.IsActive() {
		t.Fatal("Offer still active")
	}
}

func TestTrySellNoBidsCompletes(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	offer := Offer{Symbol: "m", Amount: 10, OfferType: OrderTypeMarket}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, offer)
	offer = storage.GetOffer(offer.ID)
	if !offer.IsActive() {
		t.Fatal("Offer not active")
	}
	if offer.Amount != 10 {
		t.Fatalf("%+v", offer)
	}
}

func TestTrySell2BidsSell(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	bid0 := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	bid0.ID = storage.AddBid(bid0)
	bid1 := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	bid1.ID = storage.AddBid(bid1)
	offer := Offer{Symbol: "m", Amount: 20, OfferType: OrderTypeMarket}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, offer)
	bid0 = storage.GetBid(bid0.ID)
	if bid0.IsActive() {
		t.Fatal("Bid0 still active")
	}
	bid1 = storage.GetBid(bid1.ID)
	if bid1.IsActive() {
		t.Fatal("Bid1 still active")
	}
	offer = storage.GetOffer(offer.ID)
	if offer.IsActive() {
		t.Fatal("Offer still active")
	}
}

func TestTrySellPartialBidCompletesAndDecrimentsOffer(t *testing.T) {
	mop := marketOrderProcessor{now: func() time.Time { return time.Time{} }}
	storage := MakeMemoryStorage()
	bid := Bid{Symbol: "m", Amount: 5, BidType: OrderTypeMarket}
	bid.ID = storage.AddBid(bid)
	offer := Offer{Symbol: "m", Amount: 10, OfferType: OrderTypeMarket}
	offer.ID = storage.AddOffer(offer)
	mop.TrySell(storage, makeMockAccounts(), map[OrderType]orderProcessor{OrderTypeMarket: &mop}, offer)
	bid = storage.GetBid(bid.ID)
	if bid.IsActive() {
		t.Fatal("Bid still active")
	}
	offer = storage.GetOffer(offer.ID)
	if !offer.IsActive() {
		t.Fatal("Offer inactive")
	}
	if offer.Amount != 5 {
		t.Fatalf("Wrong remaining amount: %+v", offer)
	}
}
