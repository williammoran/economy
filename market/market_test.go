package market

import "testing"

func TestOffer(t *testing.T) {
	storage := makeMemoryMarketStorage()
	m := MakeMarket(storage)
	o := Offer{Symbol: "m"}
	m.Offer(o)
	if len(storage.offers["m"]) != 1 {
		t.Fatalf("%+v", storage)
	}
}

func TestBid(t *testing.T) {

}
