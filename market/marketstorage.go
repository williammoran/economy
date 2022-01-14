package market

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type AccountID int64

type Symbol string

type Offer struct {
	Account AccountID
	Symbol  Symbol
	Amount  int64
}

type BidID int64

type BidType byte

const (
	BidMarket BidType = 1
	BidLimit  BidType = 2
)

type Bid struct {
	BidID   BidID
	BidType BidType
	Account AccountID
	Symbol  Symbol
	Volume  int64
	Amount  int64
}

type Transaction struct {
	ID    uuid.UUID
	BidID BidID
	Date  time.Time
}

type MarketStorage interface {
	AddOffer(Offer)
	AddBid(Bid) BidID
	GetTransactions(BidID) []Transaction
}

func makeMemoryMarketStorage() *memoryMarketStorage {
	return &memoryMarketStorage{
		offers: make(map[Symbol][]Offer),
		bids:   make(map[BidID]Bid),
	}
}

type memoryMarketStorage struct {
	mutex     sync.Mutex
	offers    map[Symbol][]Offer
	bids      map[BidID]Bid
	lastBidID BidID
}

func (s *memoryMarketStorage) AddOffer(o Offer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	offers := s.offers[o.Symbol]
	offers = append(offers, o)
	s.offers[o.Symbol] = offers
}

func (s *memoryMarketStorage) AddBid(b Bid) BidID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastBidID++
	b.BidID = s.lastBidID
	s.bids[s.lastBidID] = b
	return s.lastBidID
}

func (s *memoryMarketStorage) GetTransactions(id BidID) []Transaction {
	panic("incomplete")
}
