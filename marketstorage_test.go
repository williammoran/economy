package economy

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const (
	sym = "S"
)

func TestAddOffer(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Symbol: sym}
	offer.ID = ms.AddOffer(offer)
	if offer.ID == uuid.Nil {
		t.Fatalf("UUID not generated")
	}
	r := ms.offers[sym][offer.ID]
	if !reflect.DeepEqual(offer, r) {
		t.Fatalf("%+v != %+v", r, offer)
	}
}

func TestBestOfferBasic(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Symbol: sym, Amount: 10}
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
	offer := Offer{Symbol: sym, Amount: 0}
	offer.ID = ms.AddOffer(offer)
	r, found := ms.BestOffer(sym)
	if found {
		t.Fatalf("Incorrectly found %+v", r)
	}
}

func TestBestOfferSelectsCorrectly(t *testing.T) {
	ms := MakeMemoryStorage()
	offer0 := Offer{Symbol: sym, Amount: 10, Price: 5, OfferType: OrderTypeLimit}
	offer0.ID = ms.AddOffer(offer0)
	offer1 := Offer{Symbol: sym, Amount: 10, Price: 2, OfferType: OrderTypeLimit}
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
	offer0 := Offer{Symbol: sym, Amount: 10, Price: 5, OfferType: OrderTypeLimit}
	offer0.ID = ms.AddOffer(offer0)
	offer1 := Offer{Symbol: sym, Amount: 0, Price: 2, OfferType: OrderTypeLimit}
	offer1.ID = ms.AddOffer(offer1)
	r, found := ms.BestOffer(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(offer0, r) {
		t.Fatalf("%+v != %+v", r, offer0)
	}
}

func TestBestBidMarketBasic(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Symbol: sym, Amount: 10}
	bid.ID = ms.AddBid(bid)
	r, found := ms.BestBid(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(bid, r) {
		t.Fatalf("%+v != %+v", r, bid)
	}
}

func TestBestBidLimitBasic(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Symbol: sym, Amount: 10, BidType: OrderTypeLimit, Price: 5}
	bid.ID = ms.AddBid(bid)
	r, found := ms.BestBid(sym)
	if !found {
		t.Fatal("Not found")
	}
	if !reflect.DeepEqual(bid, r) {
		t.Fatalf("%+v != %+v", r, bid)
	}
}

func TestBestBidIgnoresEmpty(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Symbol: sym, Amount: 0}
	bid.ID = ms.AddBid(bid)
	r, found := ms.BestBid(sym)
	if found {
		t.Fatalf("Should not have been found: %+v", r)
	}
}

func TestBestBidSelectsCorrectly(t *testing.T) {
	ms := MakeMemoryStorage()
	bid0 := Bid{Symbol: sym, Amount: 10, Price: 5, BidType: OrderTypeLimit}
	bid0.ID = ms.AddBid(bid0)
	bid1 := Bid{Symbol: sym, Amount: 10, Price: 7, BidType: OrderTypeLimit}
	bid1.ID = ms.AddBid(bid1)
	r, found := ms.BestBid(sym)
	if !found {
		t.Fatal("Not Found")
	}
	if !reflect.DeepEqual(bid1, r) {
		t.Fatalf("Incorrectly found %+v", r)
	}
}

func TestBestBidSelectsNonempty(t *testing.T) {
	ms := MakeMemoryStorage()
	bid0 := Bid{Symbol: sym, Amount: 10, Price: 5, BidType: OrderTypeLimit}
	bid0.ID = ms.AddBid(bid0)
	bid1 := Bid{Symbol: sym, Amount: 0, Price: 7, BidType: OrderTypeLimit}
	bid1.ID = ms.AddBid(bid1)
	r, found := ms.BestBid(sym)
	if !found {
		t.Fatal("Not Found")
	}
	if !reflect.DeepEqual(bid0, r) {
		t.Fatalf("Incorrectly found %+v", r)
	}
}

func TestUpdateOffer(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Symbol: sym, ID: uuid.New()}
	ms.offers[sym] = make(map[uuid.UUID]Offer)
	ms.offers[sym][offer.ID] = offer
	offer.Amount = 10
	ms.UpdateOffer(offer)
	r := ms.offers[sym][offer.ID]
	if r != offer {
		t.Fatalf("%+v != %+v", r, offer)
	}
}

func TestAddBid(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Symbol: sym}
	bid.ID = ms.AddBid(bid)
	if bid.ID == uuid.Nil {
		t.Fatalf("UUID not generated")
	}
	r := ms.bids[bid.ID]
	if !reflect.DeepEqual(bid, r) {
		t.Fatalf("%+v != %+v", r, bid)
	}
}

func TestUpdateBid(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Symbol: sym, ID: uuid.New()}
	ms.bids[bid.ID] = bid
	bid.Amount = 10
	ms.UpdateBid(bid)
	r := ms.bids[bid.ID]
	if r != bid {
		t.Fatalf("%+v != %+v", r, bid)
	}
}

func TestGetBid(t *testing.T) {
	ms := MakeMemoryStorage()
	bid := Bid{Symbol: sym, ID: uuid.New()}
	ms.bids[bid.ID] = bid
	r := ms.GetBid(bid.ID)
	if r != bid {
		t.Fatalf("%+v != %+v", r, bid)
	}
}

func TestGetOffer(t *testing.T) {
	ms := MakeMemoryStorage()
	offer := Offer{Symbol: sym, ID: uuid.New()}
	ms.offers[sym] = make(map[uuid.UUID]Offer)
	ms.offers[sym][offer.ID] = offer
	r := ms.GetOffer(offer.ID)
	if r != offer {
		t.Fatalf("%+v != %+v", r, offer)
	}
}

func TestNewTransaction(t *testing.T) {
	ms := MakeMemoryStorage()
	tx := Transaction{BidID: uuid.New()}
	ms.NewTransaction(tx)
	if len(ms.transactions) != 1 {
		t.Fatalf("Not 1: %+v", ms.transactions)
	}
	r := ms.transactions[0]
	if r.ID == uuid.Nil {
		t.Fatal("ID not generated")
	}
	if r.BidID != tx.BidID {
		t.Fatalf("Data incorred: %+v", r)
	}
	tx.BidID = uuid.New()
	ms.NewTransaction(tx)
	if len(ms.transactions) != 2 {
		t.Fatalf("Not 2: %+v", ms.transactions)
	}
	r = ms.transactions[1]
	if r.ID == uuid.Nil {
		t.Fatal("ID not generated")
	}
	if r.BidID != tx.BidID {
		t.Fatalf("Data incorrect: %+v", r)
	}
}

func TestSetLastPrice(t *testing.T) {
	ms := MakeMemoryStorage()
	ms.SetLastPrice(sym, 55)
	if ms.lastPrice[sym] != 55 {
		t.Fatalf("Should be 55: %+v", ms.lastPrice)
	}
	ms.SetLastPrice("J", 42)
	if ms.lastPrice["J"] != 42 {
		t.Fatalf("Should be 42: %+v", ms.lastPrice)
	}
	ms.SetLastPrice(sym, 35)
	if ms.lastPrice[sym] != 35 {
		t.Fatalf("Should be 35: %+v", ms.lastPrice)
	}
}

func TestGetLastPrice(t *testing.T) {
	ms := MakeMemoryStorage()
	ms.lastPrice[sym] = 42
	if ms.LastPrice(sym) != 42 {
		t.Fatalf("Should be 42: %+v", ms.lastPrice)
	}
}

func TestAllSymbols(t *testing.T) {
	const symX = "X"
	ms := MakeMemoryStorage()
	ms.lastPrice[symX] = 42
	offer := Offer{Symbol: sym, ID: uuid.New()}
	ms.offers[sym] = make(map[uuid.UUID]Offer)
	ms.offers[sym][offer.ID] = offer
	syms := ms.AllSymbols()
	if len(syms) != 2 {
		t.Fatalf("Wrong: %+v", syms)
	}
	expected := map[string]bool{symX: true, sym: true}
	for _, s := range syms {
		if !expected[s] {
			t.Logf("Did not expect '%s'", s)
			t.Fail()
		}
		delete(expected, s)
	}
	if len(expected) != 0 {
		t.Fatalf("Mising symbols: %+v", expected)
	}
}

func TestMarshalEmpty(t *testing.T) {
	ms := MakeMemoryStorage()
	buffer := bytes.Buffer{}
	ms.Marshal(&buffer)
	t.Log(buffer.String())
	msr := MakeMemoryStorage()
	reader := bytes.NewReader(buffer.Bytes())
	msr.UnMarshal(reader)
	if !reflect.DeepEqual(ms, msr) {
		t.FailNow()
	}
}

func TestMarshalData(t *testing.T) {
	ms := MakeMemoryStorage()
	ms.AddBid(Bid{Symbol: sym, Amount: 10})
	ms.AddBid(Bid{Symbol: "G", Amount: 11, Account: 2})
	ms.AddOffer(Offer{Symbol: "Z", Amount: 14})
	ms.AddOffer(Offer{Symbol: "Y", Amount: 8, Account: 4, OfferType: OrderTypeLimit, Price: 42})
	ms.NewTransaction(Transaction{Price: 24})
	ms.NewTransaction(Transaction{Price: 424})
	ms.SetLastPrice("Q", 233)
	ms.SetLastPrice("X", 322)
	buffer := bytes.Buffer{}
	ms.Marshal(&buffer)
	t.Log("\n" + buffer.String())
	msr := MakeMemoryStorage()
	reader := bytes.NewReader(buffer.Bytes())
	msr.UnMarshal(reader)
	if !reflect.DeepEqual(ms.bids, msr.bids) {
		t.Fatalf("bids != %+v", msr.bids)
	}
	if !reflect.DeepEqual(ms.offers, msr.offers) {
		t.Fatalf("offers != %+v", msr.offers)
	}
	if !reflect.DeepEqual(ms.lastPrice, msr.lastPrice) {
		t.Fatalf("prices != %+v", msr.lastPrice)
	}
	if !reflect.DeepEqual(ms.transactions, msr.transactions) {
		t.Fatalf("transactions:\n%+v\n%+v", msr.transactions, ms.transactions)
	}
}
