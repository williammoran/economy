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
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	bid, _ = fillBid(storage, makeMockAccounts(), time.Time{}, bid, o, 7)
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
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	bid, _ = fillBid(storage, makeMockAccounts(), time.Time{}, bid, o, 7)
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
	bid := Bid{Symbol: "m", Amount: 10, BidType: OrderTypeMarket}
	id := storage.AddBid(bid)
	bid.ID = id
	bid, _ = fillBid(storage, makeMockAccounts(), time.Time{}, bid, o, 7)
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

func TestFillBidExchangesFunds(t *testing.T) {
	storage := makeMemoryMarketStorage()
	accounts := makeMockAccounts()
	o := Offer{Symbol: "m", Amount: 10, Account: 1}
	offerID := storage.AddOffer(o)
	o.ID = offerID
	bid := Bid{Symbol: "m", Amount: 10, Account: 2}
	id := storage.AddBid(bid)
	bid.ID = id
	storage.SetLastPrice("m", 10)
	bid, filled := fillBid(storage, accounts, time.Time{}, bid, o, 10)
	if !filled {
		t.Fatal("Bid not filled")
	}
	if accounts.accounts[1] != 100 {
		t.Fatalf("Not credited: %+v", accounts.accounts)
	}
	if accounts.accounts[2] != -100 {
		t.Fatalf("Not debited: %+v", accounts.accounts)
	}
}

func TestFillBidRejectsOnNoFunds(t *testing.T) {
	storage := makeMemoryMarketStorage()
	accounts := makeMockAccounts()
	accounts.rejects[2] = true
	o := Offer{Symbol: "m", Amount: 10, Account: 1}
	offerID := storage.AddOffer(o)
	o.ID = offerID
	bid := Bid{Symbol: "m", Amount: 10, Account: 2}
	id := storage.AddBid(bid)
	bid.ID = id
	storage.SetLastPrice("m", 10)
	bid, filled := fillBid(storage, accounts, time.Time{}, bid, o, 10)
	if filled {
		t.Fatal("Bid was filled")
	}
	if accounts.accounts[1] != 0 {
		t.Fatalf("Credited: %+v", accounts.accounts)
	}
	if accounts.accounts[2] != 0 {
		t.Fatalf("Debited: %+v", accounts.accounts)
	}
}
