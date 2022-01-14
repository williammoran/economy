package market

import "testing"

func TestOffer(t *testing.T) {
	m := MakeMarket(makeMemoryMarketStorage())
	o := Offer{Symbol: "m"}
	m.Offer(o)
}

func TestBid(t *testing.T) {

}
