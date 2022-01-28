// This package contains a CLI client for experimenting
// with the library's capabilities. Build and start the
// program and type "help<enter>" for more information.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/williammoran/economy"
)

const help = `
Use this CLI to experiment with the market library.

This program persists market data to CSV files in the
current directory. Thus your market activity persists
across multiple executions, unless you delete the
CSV files.

Commands are as follows:
help - this message
account $id $account_balance - create or update an account
 with the specified id and balance of funds
accounts - list all known accounts
offer $account $symbol $volume - Offer for $account to
 sell $volume of $symbol at the current market price
offer $account $symbol $volume limit $price - Offer for
 $account to sell $volume of $symbol at or above $price
bid $account $symbol $volume - Put in an order for $account
 to purchase $volume of $symbol at current market price
bid $account $symbol $volume limit $price - Put in an
 order for $account to purchase $volume of $symbol at any
 price at or below $price
market - List current prices of all known symbols
`

func main() {
	accounts := makeAccounts()
	storage := economy.MakeMemoryStorage()
	market := economy.MakeMarket(time.Now, storage, accounts)
	reader := bufio.NewReader(os.Stdin)
	for {
		tokens := nextCommand(reader)
		if len(tokens) == 0 {
			continue
		}
		command := strings.ToLower(tokens[0])
		switch command {
		case "help":
			fmt.Print(help)
		case "account":
			setAccount(tokens[1:], accounts)
		case "accounts":
			showAccounts(accounts)
		case "bid":
			load(storage)
			bid(tokens[1:], market)
			save(storage)
		case "offer":
			load(storage)
			offer(tokens[1:], market)
			save(storage)
		case "market":
			load(storage)
			showMarket(market)
		default:
			fmt.Printf("Unrecognized command '%s'\n", command)
		}
	}
}

func nextCommand(r io.ByteReader) []string {
	var buf []byte
	for {
		c, err := r.ReadByte()
		if err != nil {
			break
		}
		if c == '\n' {
			break
		}
		if c != '\r' {
			buf = append(buf, c)
		}
	}
	c := strings.TrimSpace(string(buf))
	space := regexp.MustCompile(`\s+`)
	c = space.ReplaceAllString(c, " ")
	return strings.Split(c, " ")
}

func setAccount(c []string, accounts *accounts) {
	id, ok := parseAccount(c[0])
	if !ok {
		return
	}
	amount, ok := parseAmount(c[1])
	if !ok {
		return
	}
	accounts.accounts[id] = amount
	fmt.Printf("Account %d now has %d\n", id, amount)
}

func showAccounts(accounts *accounts) {
	fmt.Println("AccountID   Balance")
	for id, balance := range accounts.accounts {
		fmt.Printf("%9d %9d\n", id, balance)
	}
}

// bid $account $symbol $volume
// bid $account $symbol $volume limit $price
func bid(c []string, market *economy.Market) {
	bidType, account, symbol, volume, price, ok := parseBidOrOffer(c)
	if !ok {
		return
	}
	bid := economy.Bid{
		BidType: bidType,
		Account: account,
		Symbol:  symbol,
		Price:   price,
		Amount:  volume,
	}
	market.Bid(bid)
	fmt.Printf("Made bid %+v\n", bid)
}

// offer $account $symbol $volume
// offer $account $symbol $volume limit $price
func offer(c []string, market *economy.Market) {
	offerType, account, symbol, volume, price, ok := parseBidOrOffer(c)
	if !ok {
		return
	}
	offer := economy.Offer{
		OfferType: offerType,
		Account:   account,
		Symbol:    symbol,
		Price:     price,
		Amount:    volume,
	}
	market.Offer(offer)
	fmt.Printf("Made offer %+v\n", offer)
}

// $_ $account $symbol $volume
// $_ $account $symbol $volume limit $price
func parseBidOrOffer(c []string) (economy.OrderType, int64, string, int64, int64, bool) {
	var ok bool
	var orderType economy.OrderType
	var price int64
	if len(c) == 3 {
		orderType = economy.OrderTypeMarket
	} else {
		if len(c) == 5 {
			orderType = economy.OrderTypeLimit
			price, ok = parsePrice(c[4])
			if !ok {
				return 0, 0, "", 0, 0, false
			}
		} else {
			fmt.Printf("Invalid offer %+v", c)
			return 0, 0, "", 0, 0, false
		}
	}
	symbol := c[1]
	account, ok := parseAccount(c[0])
	if !ok {
		return 0, 0, "", 0, 0, false
	}
	amount, ok := parseAmount(c[2])
	if !ok {
		return 0, 0, "", 0, 0, false
	}
	return orderType, account, symbol, amount, price, true
}

func showMarket(market *economy.Market) {
	symbols := market.AllSymbols()
	fmt.Println("Symbol Last Price")
	for _, symbol := range symbols {
		price := market.LastPrice(symbol)
		fmt.Printf("%6s %10d\n", symbol, price)
	}
}

func parseAccount(in string) (int64, bool) {
	return parseInt64(in, "Account ID must be an int64")
}

func parseAmount(in string) (int64, bool) {
	return parseInt64(in, "Amount must be an int64")
}

func parsePrice(in string) (int64, bool) {
	return parseInt64(in, "Price must be an int64")
}

func parseInt64(in, msg string) (int64, bool) {
	id, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(msg)
		return 0, false
	}
	return id, true
}

const filename = "economy.data"

func load(s *economy.MemoryStorage) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	s.UnMarshal(f)
}

func save(s *economy.MemoryStorage) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	s.Marshal(f)
}
