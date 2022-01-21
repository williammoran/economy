package economy

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Offer struct {
	ID        uuid.UUID
	OfferType OrderType
	Account   int64
	Symbol    string
	Price     int64
	Amount    int64
}

type OrderType byte

const (
	BidMarket        OrderType = 0
	OfferMarket      OrderType = 0
	BidLimit         OrderType = 1
	OfferLimit       OrderType = 1
	BidStatusPending           = 0
	BidStatusFilled            = 1
)

type Bid struct {
	ID      int64
	BidType OrderType
	Account int64
	Symbol  string
	Price   int64
	Amount  int64
	Status  byte
}

type Transaction struct {
	ID      uuid.UUID
	BidID   int64
	OfferID uuid.UUID
	Price   int64
	Amount  int64
	Date    time.Time
}

type MarketStorage interface {
	Lock()
	Unlock()
	// AddOffer returns the UUID of the created offer
	AddOffer(Offer) uuid.UUID
	// BestOffer returns the offer with the best price
	// for the specified string, or false if there are
	// no offers
	BestOffer(string) (Offer, bool)
	UpdateOffer(Offer)
	AddBid(Bid) int64
	UpdateBid(Bid)
	GetBid(int64) Bid
	NewTransaction(Transaction)
	LastPrice(string) int64
	SetLastPrice(string, int64)
	// Return all the known strings
	AllSymbols() []string
}

func makeMemoryMarketStorage() *memoryMarketStorage {
	return &memoryMarketStorage{
		offers:    make(map[string]map[uuid.UUID]Offer),
		bids:      make(map[int64]Bid),
		lastPrice: make(map[string]int64),
	}
}

type memoryMarketStorage struct {
	mutex        sync.Mutex
	offers       map[string]map[uuid.UUID]Offer
	bids         map[int64]Bid
	transactions []Transaction
	lastint64    int64
	lastPrice    map[string]int64
}

func (s *memoryMarketStorage) Lock() {
	s.mutex.Lock()
}

func (s *memoryMarketStorage) Unlock() {
	s.mutex.Unlock()
}

func (s *memoryMarketStorage) AddOffer(o Offer) uuid.UUID {
	o.ID = uuid.New()
	offers := s.offers[o.Symbol]
	if offers == nil {
		offers = make(map[uuid.UUID]Offer)
	}
	offers[o.ID] = o
	s.offers[o.Symbol] = offers
	return o.ID
}

func (s *memoryMarketStorage) BestOffer(sym string) (Offer, bool) {
	l := s.offers[sym]
	if len(l) == 0 {
		return Offer{}, false
	}
	o := Offer{Price: 2 ^ 60}
	for _, offer := range l {
		if offer.Price < o.Price {
			o = offer
		}
	}
	return o, true
}

func (s *memoryMarketStorage) UpdateOffer(o Offer) {
	l := s.offers[o.Symbol]
	if o.Amount > 0 {
		l[o.ID] = o
	} else {
		delete(l, o.ID)
	}
	s.offers[o.Symbol] = l
}

func (s *memoryMarketStorage) AddBid(b Bid) int64 {
	s.lastint64++
	b.ID = s.lastint64
	s.bids[s.lastint64] = b
	return s.lastint64
}

func (s *memoryMarketStorage) UpdateBid(b Bid) {
	s.bids[b.ID] = b
}

func (s *memoryMarketStorage) GetBid(id int64) Bid {
	bid, found := s.bids[s.lastint64]
	if !found {
		log.Panicf("Bid %d not found", id)
	}
	return bid
}

func (s *memoryMarketStorage) NewTransaction(t Transaction) {
	t.ID = uuid.New()
	s.transactions = append(s.transactions, t)
}

func (s *memoryMarketStorage) LastPrice(string string) int64 {
	p, found := s.lastPrice[string]
	if found {
		return p
	}
	return 1
}

func (s *memoryMarketStorage) SetLastPrice(
	string string, price int64,
) {
	s.lastPrice[string] = price
}

func (s *memoryMarketStorage) AllSymbols() []string {
	l := make(map[string]bool)
	for s := range s.lastPrice {
		l[s] = true
	}
	for s := range s.offers {
		l[s] = true
	}
	var rv []string
	for s := range l {
		rv = append(rv, s)
	}
	return rv
}
