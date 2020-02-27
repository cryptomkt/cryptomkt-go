package conn

import (
	"bytes"
	"strconv"
)

type Pagination struct {
	Previous interface{}
	Next     interface{}
	Limit    int
	Page     int
}

func (pagination *Pagination) String() string {
	var b bytes.Buffer
	b.WriteString("Pagination{")
	b.WriteString("previous:")
	if pagination.Previous != nil {
		b.WriteString(strconv.FormatFloat(pagination.Previous.(float64), 'f', -1, 64))
	} else {
		b.WriteString("nil")
	}
	b.WriteString(" next:")
	if pagination.Next != nil {
		b.WriteString(strconv.FormatFloat(pagination.Next.(float64), 'f', -1, 64))
	} else {
		b.WriteString("nil")
	}
	b.WriteString(" limit:")
	b.WriteString(strconv.Itoa(pagination.Limit))
	b.WriteString(" page:")
	b.WriteString(strconv.Itoa(pagination.Page))
	b.WriteString("}")
	return b.String()
}

type MarketListResponse struct {
	Status  string
	Message string
	Data    []string
}

type TickerResponse struct {
	Status  string
	Message string
	Data    []Ticker
}

type Ticker struct {
	High      string
	Volume    string
	Low       string
	Ask       string
	Timestamp string
	Bid       string
	LastPrice string `json:"last_price"`
	Market    string
}

type BalancesResponse struct {
	Status  string
	Message string
	Data    []Balance
}

type Balance struct {
	Wallet    string
	Available string
	Balance   string
}

type AccountResponse struct {
	Status  string
	Message string
	Data    Account
}

type Account struct {
	Name         string
	Email        string
	Rate         Rate
	BankAccounts []BankAccount `json:"bank_accounts"`
}

type Rate struct {
	MarketMaker string `json:"market_maker"`
	MarketTaker string `json:"market_taker"`
}

type BankAccount struct {
	Id          int
	Bank        string
	Description string
	Country     string
	Number      string
}

type InstantResponse struct {
	Status  string
	Message string
	Data    Instant
}

type Instant struct {
	Obtained float64
	Required float64
}

func (instant *Instant) String() string {
	return "instant{obtained:" + 
	strconv.FormatFloat(instant.Obtained, 'f', -1, 64) +
	" required:" +
	strconv.FormatFloat(instant.Required, 'f', -1, 64) + 
	"}"
}
