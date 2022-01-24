package economy

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

func makeMemoryMarketStorage() *memoryMarketStorage {
	return &memoryMarketStorage{
		offers:    make(map[string]map[uuid.UUID]Offer),
		bids:      make(map[uuid.UUID]Bid),
		lastPrice: make(map[string]int64),
	}
}

type memoryMarketStorage struct {
	mutex        sync.Mutex
	offers       map[string]map[uuid.UUID]Offer
	bids         map[uuid.UUID]Bid
	transactions []Transaction
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
	marketPrice := s.LastPrice(sym)
	for _, offer := range l {
		if offer.IsActive() {
			switch offer.OfferType {
			case OrderTypeLimit:
				if offer.Price < o.Price {
					o = offer
				}
			case OrderTypeMarket:
				if marketPrice < o.Price {
					o = offer
				}
			default:
				log.Panicf("Unknown offer type %d", offer.OfferType)
			}
		}
	}
	if o.Amount > 0 {
		return o, true
	}
	return Offer{}, false
}

func (s *memoryMarketStorage) BestBid(sym string) (Bid, bool) {
	b := Bid{Price: 0}
	marketPrice := s.LastPrice(sym)
	for _, bid := range s.bids {
		if bid.IsActive() {
			switch bid.BidType {
			case OrderTypeLimit:
				if bid.Price > b.Price {
					b = bid
				}
			case OrderTypeMarket:
				if marketPrice > b.Price {
					b = bid
				}
			default:
				log.Panicf("Unknown bid type %d", bid.BidType)
			}
		}
	}
	if b.Amount > 0 {
		return b, true
	}
	return Bid{}, false
}

func (s *memoryMarketStorage) UpdateOffer(o Offer) {
	l := s.offers[o.Symbol]
	if l == nil {
		l = make(map[uuid.UUID]Offer)
	}
	l[o.ID] = o
	s.offers[o.Symbol] = l
}

func (s *memoryMarketStorage) AddBid(b Bid) uuid.UUID {
	b.ID = uuid.New()
	s.bids[b.ID] = b
	return b.ID
}

func (s *memoryMarketStorage) UpdateBid(b Bid) {
	s.bids[b.ID] = b
}

func (s *memoryMarketStorage) GetBid(id uuid.UUID) Bid {
	bid, found := s.bids[id]
	if !found {
		log.Panicf("Bid %s not found", id)
	}
	return bid
}

func (s *memoryMarketStorage) GetOffer(id uuid.UUID) Offer {
	for _, l := range s.offers {
		if offer, found := l[id]; found {
			return offer
		}
	}
	panic(fmt.Sprintf("Offer %s not found", id))
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
