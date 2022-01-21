package economy

import (
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
	OrderTypeMarket  OrderType = 0
	OrderTypeLimit   OrderType = 1
	BidStatusPending           = 0
	BidStatusFilled            = 1
)

type Bid struct {
	ID      uuid.UUID
	BidType OrderType
	Account int64
	Symbol  string
	Price   int64
	Amount  int64
	Status  byte
}

type Transaction struct {
	ID      uuid.UUID
	BidID   uuid.UUID
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
	AddBid(Bid) uuid.UUID
	UpdateBid(Bid)
	GetBid(uuid.UUID) Bid
	NewTransaction(Transaction)
	LastPrice(string) int64
	SetLastPrice(string, int64)
	// Return all the known strings
	AllSymbols() []string
}

type OrderProcessor interface {
	TryFillBid(MarketStorage, map[OrderType]OrderProcessor, Bid)
	GetAskingPrice(MarketStorage, Offer) int64
}

func MakeMarket(s MarketStorage) *Market {
	return &Market{
		storage: s,
		orderProcessors: map[OrderType]OrderProcessor{
			OrderTypeMarket: &marketOrderProcessor{now: time.Now},
			OrderTypeLimit:  &limitOrderProcessor{now: time.Now},
		},
	}
}

type Market struct {
	storage         MarketStorage
	orderProcessors map[OrderType]OrderProcessor
}

func (m *Market) Offer(o Offer) {
	m.storage.Lock()
	defer m.storage.Unlock()
	m.storage.AddOffer(o)
}

func (m *Market) Bid(b Bid) uuid.UUID {
	m.storage.Lock()
	defer m.storage.Unlock()
	b.ID = m.storage.AddBid(b)
	m.orderProcessors[b.BidType].TryFillBid(m.storage, m.orderProcessors, b)
	return b.ID
}

func (m *Market) GetBid(id uuid.UUID) Bid {
	m.storage.Lock()
	defer m.storage.Unlock()
	return m.storage.GetBid(id)
}

func (m *Market) AllSymbols() []string {
	m.storage.Lock()
	defer m.storage.Unlock()
	return m.storage.AllSymbols()
}

func (m *Market) LastPrice(s string) int64 {
	m.storage.Lock()
	defer m.storage.Unlock()
	return m.storage.LastPrice(s)
}
