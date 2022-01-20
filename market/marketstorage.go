package market

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AccountID int64

type Symbol string

type Offer struct {
	ID        uuid.UUID
	OfferType OrderType
	Account   AccountID
	Symbol    Symbol
	Price     int64
	Amount    int64
}

type BidID int64

type OrderType byte

const (
	BidMarket        OrderType = 0
	OfferMarket      OrderType = 0
	BidLimit         OrderType = 1
	BidStatusPending           = 0
	BidStatusFilled            = 1
)

type Bid struct {
	BidID   BidID
	BidType OrderType
	Account AccountID
	Symbol  Symbol
	Price   int64
	Amount  int64
	Status  byte
}

type Transaction struct {
	ID      uuid.UUID
	BidID   BidID
	OfferID uuid.UUID
	Price   int64
	Amount  int64
	Date    time.Time
}

type MarketStorage interface {
	Lock()
	Unlock()
	AddOffer(Offer)
	BestOffer(Symbol) (Offer, bool)
	UpdateOffer(Offer)
	AddBid(Bid) BidID
	UpdateBid(Bid)
	GetBid(BidID) Bid
	NewTransaction(Transaction)
	LastPrice(Symbol) int64
	SetLastPrice(Symbol, int64)
}

func makeMemoryMarketStorage() *memoryMarketStorage {
	return &memoryMarketStorage{
		offers:    make(map[Symbol]map[uuid.UUID]Offer),
		bids:      make(map[BidID]Bid),
		lastPrice: make(map[Symbol]int64),
	}
}

type memoryMarketStorage struct {
	mutex        sync.Mutex
	offers       map[Symbol]map[uuid.UUID]Offer
	bids         map[BidID]Bid
	transactions []Transaction
	lastBidID    BidID
	lastPrice    map[Symbol]int64
}

func (s *memoryMarketStorage) Lock() {
	s.mutex.Lock()
}

func (s *memoryMarketStorage) Unlock() {
	s.mutex.Unlock()
}

func (s *memoryMarketStorage) AddOffer(o Offer) {
	o.ID = uuid.New()
	offers := s.offers[o.Symbol]
	if offers == nil {
		offers = make(map[uuid.UUID]Offer)
	}
	offers[o.ID] = o
	s.offers[o.Symbol] = offers
}

func (s *memoryMarketStorage) BestOffer(sym Symbol) (Offer, bool) {
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

func (s *memoryMarketStorage) AddBid(b Bid) BidID {
	s.lastBidID++
	b.BidID = s.lastBidID
	s.bids[s.lastBidID] = b
	return s.lastBidID
}

func (s *memoryMarketStorage) UpdateBid(b Bid) {
	s.bids[b.BidID] = b
}

func (s *memoryMarketStorage) GetBid(id BidID) Bid {
	bid, found := s.bids[s.lastBidID]
	if !found {
		log.Panicf("Bid %d not found", id)
	}
	return bid
}

func (s *memoryMarketStorage) NewTransaction(t Transaction) {
	t.ID = uuid.New()
	s.transactions = append(s.transactions, t)
}

func (s *memoryMarketStorage) LastPrice(symbol Symbol) int64 {
	p, found := s.lastPrice[symbol]
	if found {
		return p
	}
	return 1
}

func (s *memoryMarketStorage) SetLastPrice(
	symbol Symbol, price int64,
) {
	s.lastPrice[symbol] = price
}
