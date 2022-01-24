package economy

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

func MakeMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		offers:    make(map[string]map[uuid.UUID]Offer),
		bids:      make(map[uuid.UUID]Bid),
		lastPrice: make(map[string]int64),
	}
}

const eof = "EOF"

type MemoryStorage struct {
	mutex        sync.Mutex
	offers       map[string]map[uuid.UUID]Offer
	bids         map[uuid.UUID]Bid
	transactions []Transaction
	lastPrice    map[string]int64
}

func (s *MemoryStorage) Lock() {
	s.mutex.Lock()
}

func (s *MemoryStorage) Unlock() {
	s.mutex.Unlock()
}

func (s *MemoryStorage) Marshal(w io.Writer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	offers := make(map[uuid.UUID]Offer)
	for _, oList := range s.offers {
		for _, o := range oList {
			offers[o.ID] = o
		}
	}
	saveOffers(w, offers)
	writeEOF(w)
	savePrices(w, s.lastPrice)
	writeEOF(w)
	saveBids(w, s.bids)
	writeEOF(w)
	saveTransactions(w, s.transactions)
	writeEOF(w)
}

func writeEOF(w io.Writer) {
	fmt.Fprint(w, eof+"\n")
}

func (s *MemoryStorage) UnMarshal(r io.Reader) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	offers := loadOffers(r)
	s.offers = make(map[string]map[uuid.UUID]Offer)
	for _, o := range offers {
		oList := s.offers[o.Symbol]
		if oList == nil {
			oList = make(map[uuid.UUID]Offer)
		}
		oList[o.ID] = o
		s.offers[o.Symbol] = oList
	}
	s.bids = loadBids(r)
	s.lastPrice = loadPrices(r)
	s.transactions = loadTransactions(r)
}

func (s *MemoryStorage) AddOffer(o Offer) uuid.UUID {
	o.ID = uuid.New()
	offers := s.offers[o.Symbol]
	if offers == nil {
		offers = make(map[uuid.UUID]Offer)
	}
	offers[o.ID] = o
	s.offers[o.Symbol] = offers
	return o.ID
}

func (s *MemoryStorage) BestOffer(sym string) (Offer, bool) {
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

func (s *MemoryStorage) BestBid(sym string) (Bid, bool) {
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

func (s *MemoryStorage) UpdateOffer(o Offer) {
	l := s.offers[o.Symbol]
	if l == nil {
		l = make(map[uuid.UUID]Offer)
	}
	l[o.ID] = o
	s.offers[o.Symbol] = l
}

func (s *MemoryStorage) AddBid(b Bid) uuid.UUID {
	b.ID = uuid.New()
	s.bids[b.ID] = b
	return b.ID
}

func (s *MemoryStorage) UpdateBid(b Bid) {
	s.bids[b.ID] = b
}

func (s *MemoryStorage) GetBid(id uuid.UUID) Bid {
	bid, found := s.bids[id]
	if !found {
		log.Panicf("Bid %d not found", id)
	}
	return bid
}

func (s *MemoryStorage) GetOffer(id uuid.UUID) Offer {
	for _, l := range s.offers {
		bid, found := l[id]
		if found {
			return bid
		}
	}
	panic(fmt.Sprintf("Bid %d not found", id))
}

func (s *MemoryStorage) NewTransaction(t Transaction) {
	t.ID = uuid.New()
	s.transactions = append(s.transactions, t)
}

func (s *MemoryStorage) LastPrice(symbol string) int64 {
	p, found := s.lastPrice[symbol]
	if found {
		return p
	}
	return 1
}

func (s *MemoryStorage) SetLastPrice(
	symbol string, price int64,
) {
	s.lastPrice[symbol] = price
}

func (s *MemoryStorage) AllSymbols() []string {
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

func loadOffers(r io.Reader) map[uuid.UUID]Offer {
	offers := make(map[uuid.UUID]Offer)
	reader := csv.NewReader(r)
	for {
		record, err := reader.Read()
		if err != nil {
			panic(err.Error())
		}
		if len(record) == 1 {
			if record[0] == eof {
				return offers
			}
			log.Panicf("Invalid record %+v", record)
		}
		offer := Offer{
			ID:        uuid.MustParse(record[0]),
			OfferType: OrderType(mustParseByte(record[1])),
			Account:   mustParseInt64(record[2]),
			Symbol:    record[3],
			Price:     mustParseInt64(record[4]),
			Amount:    mustParseInt64(record[5]),
		}
		offers[offer.ID] = offer
	}
}

func saveOffers(w io.Writer, offers map[uuid.UUID]Offer) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	for _, offer := range offers {
		var r []string
		r = append(r, offer.ID.String())
		r = append(r, fmt.Sprintf("%d", offer.OfferType))
		r = append(r, fmt.Sprintf("%d", offer.Account))
		r = append(r, offer.Symbol)
		r = append(r, fmt.Sprintf("%d", offer.Price))
		r = append(r, fmt.Sprintf("%d", offer.Amount))
		writer.Write(r)
	}
}

func loadPrices(r io.Reader) map[string]int64 {
	prices := make(map[string]int64)
	reader := csv.NewReader(r)
	for {
		record, err := reader.Read()
		if err != nil {
			panic(err.Error())
		}
		if len(record) == 1 {
			if record[0] == eof {
				return prices
			}
			log.Panicf("Invalid record %+v", record)
		}
		symbol := record[0]
		prices[symbol] = mustParseInt64(record[1])
	}
}

func savePrices(w io.Writer, prices map[string]int64) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	for symbol, price := range prices {
		var r []string
		r = append(r, symbol)
		r = append(r, fmt.Sprintf("%d", price))
		writer.Write(r)
	}
}

func loadBids(r io.Reader) map[uuid.UUID]Bid {
	reader := csv.NewReader(r)
	bids := make(map[uuid.UUID]Bid)
	for {
		record, err := reader.Read()
		if err != nil {
			panic(err.Error())
		}
		if len(record) == 1 {
			if record[0] == eof {
				return bids
			}
			log.Panicf("Unrecognized line: %+v", record)
		}
		bid := Bid{
			ID:      uuid.MustParse(record[0]),
			BidType: OrderType(mustParseByte(record[1])),
			Account: mustParseInt64(record[2]),
			Symbol:  record[3],
			Price:   mustParseInt64(record[4]),
			Amount:  mustParseInt64(record[5]),
			NSF:     mustParseBool(record[6]),
		}
		bids[bid.ID] = bid
	}
}

func saveBids(w io.Writer, bids map[uuid.UUID]Bid) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	for _, bid := range bids {
		var r []string
		r = append(r, bid.ID.String())
		r = append(r, fmt.Sprintf("%d", bid.BidType))
		r = append(r, fmt.Sprintf("%d", bid.Account))
		r = append(r, bid.Symbol)
		r = append(r, fmt.Sprintf("%d", bid.Price))
		r = append(r, fmt.Sprintf("%d", bid.Amount))
		r = append(r, fmt.Sprintf("%t", bid.NSF))
		writer.Write(r)
	}
}

func loadTransactions(r io.Reader) []Transaction {
	var txs []Transaction
	reader := csv.NewReader(r)
	for {
		record, err := reader.Read()
		if err != nil {
			panic(err.Error())
		}
		if len(record) == 1 {
			if record[0] == eof {
				return txs
			}
			log.Panicf("Invalid record %+v", record)
		}
		tx := Transaction{
			ID:      uuid.MustParse(record[0]),
			BidID:   uuid.MustParse(record[1]),
			OfferID: uuid.MustParse(record[2]),
			Price:   mustParseInt64(record[3]),
			Amount:  mustParseInt64(record[4]),
			Date:    time.UnixMilli(mustParseInt64(record[5])),
		}
		txs = append(txs, tx)
	}
}

func saveTransactions(w io.Writer, txs []Transaction) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	for _, tx := range txs {
		var r []string
		r = append(r, tx.ID.String())
		r = append(r, tx.BidID.String())
		r = append(r, tx.OfferID.String())
		r = append(r, fmt.Sprintf("%d", tx.Price))
		r = append(r, fmt.Sprintf("%d", tx.Amount))
		r = append(r, fmt.Sprintf("%d", tx.Date.UnixMilli()))
		writer.Write(r)
	}
}

func mustParseByte(v string) byte {
	r, err := strconv.ParseInt(v, 10, 8)
	if err != nil {
		panic(err.Error())
	}
	return byte(r)
}

func mustParseInt64(v string) int64 {
	r, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return r
}

func mustParseBool(v string) bool {
	r, err := strconv.ParseBool(v)
	if err != nil {
		panic(err.Error())
	}
	return r
}
