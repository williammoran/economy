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
	OrderTypeMarket OrderType = 0
	OrderTypeLimit  OrderType = 1
)

type Bid struct {
	ID      uuid.UUID
	BidType OrderType
	Account int64
	Symbol  string
	Price   int64
	Amount  int64
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

// Accounts provides a method for code to inject a callback
// for crediteding or debiteding funds when transactions
// occur. Note that both functions must be transaction
// safe otherwise funds could go missing.
type Accounts interface {
	// Credit must add the specified funds to the
	// spedified account
	Credit(accountID, funds int64)
	// DebitIfPossible must debit the specified funds
	// or return false
	DebitIfPossible(accountID, funds int64) bool
}

type orderProcessor interface {
	TryFillBid(MarketStorage, Accounts, map[OrderType]orderProcessor, Bid)
	GetAskingPrice(MarketStorage, Offer) int64
}

func MakeMarket(s MarketStorage, a Accounts) *Market {
	return &Market{
		storage: s,
		orderProcessors: map[OrderType]orderProcessor{
			OrderTypeMarket: &marketOrderProcessor{now: time.Now},
			OrderTypeLimit:  &limitOrderProcessor{now: time.Now},
		},
	}
}

type Market struct {
	storage         MarketStorage
	accounts        Accounts
	orderProcessors map[OrderType]orderProcessor
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
	m.orderProcessors[b.BidType].TryFillBid(
		m.storage, m.accounts, m.orderProcessors, b,
	)
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
