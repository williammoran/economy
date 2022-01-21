package main

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/williammoran/economy"
)

func makeMemoryMarketStorage() *memoryMarketStorage {
	return &memoryMarketStorage{
		offers:    make(map[economy.Symbol]map[uuid.UUID]economy.Offer),
		bids:      make(map[economy.BidID]economy.Bid),
		lastPrice: make(map[economy.Symbol]int64),
	}
}

type memoryMarketStorage struct {
	mutex        sync.Mutex
	offers       map[economy.Symbol]map[uuid.UUID]economy.Offer
	bids         map[economy.BidID]economy.Bid
	transactions []economy.Transaction
	lastBidID    economy.BidID
	lastPrice    map[economy.Symbol]int64
}

func (s *memoryMarketStorage) Lock() {
	s.mutex.Lock()
}

func (s *memoryMarketStorage) Unlock() {
	s.mutex.Unlock()
}

func (s *memoryMarketStorage) AddOffer(o economy.Offer) uuid.UUID {
	o.ID = uuid.New()
	offers := s.offers[o.Symbol]
	if offers == nil {
		offers = make(map[uuid.UUID]economy.Offer)
	}
	offers[o.ID] = o
	s.offers[o.Symbol] = offers
	return o.ID
}

func (s *memoryMarketStorage) BestOffer(sym economy.Symbol) (economy.Offer, bool) {
	l := s.offers[sym]
	if len(l) == 0 {
		return economy.Offer{}, false
	}
	o := economy.Offer{Price: 2 ^ 60}
	for _, offer := range l {
		if offer.Price < o.Price {
			o = offer
		}
	}
	return o, true
}

func (s *memoryMarketStorage) UpdateOffer(o economy.Offer) {
	l := s.offers[o.Symbol]
	if o.Amount > 0 {
		l[o.ID] = o
	} else {
		delete(l, o.ID)
	}
	s.offers[o.Symbol] = l
}

func (s *memoryMarketStorage) AddBid(b economy.Bid) economy.BidID {
	s.lastBidID++
	b.BidID = s.lastBidID
	s.bids[s.lastBidID] = b
	return s.lastBidID
}

func (s *memoryMarketStorage) UpdateBid(b economy.Bid) {
	s.bids[b.BidID] = b
}

func (s *memoryMarketStorage) GetBid(id economy.BidID) economy.Bid {
	bid, found := s.bids[s.lastBidID]
	if !found {
		log.Panicf("Bid %d not found", id)
	}
	return bid
}

func (s *memoryMarketStorage) NewTransaction(t economy.Transaction) {
	t.ID = uuid.New()
	s.transactions = append(s.transactions, t)
}

func (s *memoryMarketStorage) LastPrice(symbol economy.Symbol) int64 {
	p, found := s.lastPrice[symbol]
	if found {
		return p
	}
	return 1
}

func (s *memoryMarketStorage) SetLastPrice(
	symbol economy.Symbol, price int64,
) {
	s.lastPrice[symbol] = price
}

func (s *memoryMarketStorage) AllSymbols() []economy.Symbol {
	l := make(map[economy.Symbol]bool)
	for s := range s.lastPrice {
		l[s] = true
	}
	for s := range s.offers {
		l[s] = true
	}
	var rv []economy.Symbol
	for s := range l {
		rv = append(rv, s)
	}
	return rv
}
