package economy

import (
	"reflect"
	"testing"
)

const (
	account = 99
	sym     = "S"
)

func TestAddOffer(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Account: account, Symbol: sym}
	offer.ID = ms.AddOffer(offer)
	r := ms.offers[sym][offer.ID]
	if !reflect.DeepEqual(offer, r) {
		t.Fatalf("%+v != %+v", r, offer)
	}
}

func TestBestOfferBasic(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Account: account, Symbol: sym, Amount: 10}
	offer.ID = ms.AddOffer(offer)
	r, found := ms.BestOffer(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(offer, r) {
		t.Fatalf("%+v != %+v", r, offer)
	}
}

func TestBestOfferIgnoresEmpty(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Account: account, Symbol: sym, Amount: 0}
	offer.ID = ms.AddOffer(offer)
	r, found := ms.BestOffer(sym)
	if found {
		t.Fatalf("Incorrectly found %+v", r)
	}
}

func TestBestOfferSelectsCorrectly(t *testing.T) {
	ms := MakeMemoryStorage()
	offer0 := Offer{Account: account, Symbol: sym, Amount: 10, Price: 5, OfferType: OrderTypeLimit}
	offer0.ID = ms.AddOffer(offer0)
	offer1 := Offer{Account: account, Symbol: sym, Amount: 10, Price: 2, OfferType: OrderTypeLimit}
	offer1.ID = ms.AddOffer(offer1)
	r, found := ms.BestOffer(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(offer1, r) {
		t.Fatalf("%+v != %+v", r, offer1)
	}
}

func TestBestOfferSelectsNonEmpty(t *testing.T) {
	ms := MakeMemoryStorage()
	offer0 := Offer{Account: account, Symbol: sym, Amount: 10, Price: 5, OfferType: OrderTypeLimit}
	offer0.ID = ms.AddOffer(offer0)
	offer1 := Offer{Account: account, Symbol: sym, Amount: 0, Price: 2, OfferType: OrderTypeLimit}
	offer1.ID = ms.AddOffer(offer1)
	r, found := ms.BestOffer(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(offer0, r) {
		t.Fatalf("%+v != %+v", r, offer0)
	}
}

func TestBestBidBasic(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Account: account, Symbol: sym, Amount: 10}
	bid.ID = ms.AddBid(bid)
	r, found := ms.BestBid(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(bid, r) {
		t.Fatalf("%+v != %+v", r, bid)
	}
}
