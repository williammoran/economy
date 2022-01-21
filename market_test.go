package economy

import (
	"testing"
)

func TestOfferAddedToStorage(t *testing.T) {
	storage := makeMemoryMarketStorage()
	m := MakeMarket(storage, makeMockAccounts())
	o := Offer{Symbol: "m"}
	m.Offer(o)
	if len(storage.offers["m"]) != 1 {
		t.Fatalf("%+v", storage)
	}
}

func TestBidAddedToStorage(t *testing.T) {
	storage := makeMemoryMarketStorage()
	m := MakeMarket(storage, makeMockAccounts())
	m.orderProcessors = map[OrderType]orderProcessor{
		OrderTypeMarket: &mockOrderProcessor{},
	}
	b := Bid{}
	id := m.Bid(b)
	if _, found := storage.bids[id]; !found {
		t.Fatalf("%+v", storage)
	}
}
