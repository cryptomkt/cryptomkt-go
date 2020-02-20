package conn

import (
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/args"
)

type TradesResponse struct {
	Status     string
	Message    string
	Pagination Pagination
	Data       []TradeData
}

type Trades struct {
	args       map[string]string
	pagination Pagination
	client     *Client
	Data       []TradeData
}

type TradeData struct {
	MarketTaker string `json:"market_taker"`
	Price       string
	Amount      string
	Tid         string
	Timestamp   string
	Market      string
}

// GetPrevious lets you go to the previous page if it exists, returns (*Trades, nil) if
// it is successfull and (nil, error) otherwise
func (t *Trades) GetPrevious() (*Trades, error) {
	if t.pagination.Previous == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	var newArgs []args.Argument = make([]args.Argument, 0, 5)
	// there is always a market and a pagination
	newArgs = append(newArgs, args.Market(t.args["market"]))
	newArgs = append(newArgs, args.Page(t.pagination.Previous.(int)))
	newArgs = append(newArgs, args.Limit(t.pagination.Limit))
	if v, ok := t.args["start"]; ok {
		newArgs = append(newArgs, args.Start(v))
	}
	if v, ok := t.args["end"]; ok {
		newArgs = append(newArgs, args.Start(v))
	}
	return t.client.GetTrades(newArgs...)
}

// GetNext lets you go to the next page if it exists, returns (*Trades, nil) if
// it is successfull and (nil, error) otherwise
func (t *Trades) GetNext() (*Trades, error) {
	if t.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	var newArgs []args.Argument = make([]args.Argument, 0, 5)
	// there is always a market and a pagination
	newArgs = append(newArgs, args.Market(t.args["market"]))
	newArgs = append(newArgs, args.Page(t.pagination.Next.(int)))
	newArgs = append(newArgs, args.Limit(t.pagination.Limit))
	if v, ok := t.args["start"]; ok {
		newArgs = append(newArgs, args.Start(v))
	}
	if v, ok := t.args["end"]; ok {
		newArgs = append(newArgs, args.Start(v))
	}
	return t.client.GetTrades(newArgs...)
}

// GetPage returns the actual page of the request.
func (t *Trades) GetPage() int {
	return t.pagination.Page
}

// GetLimit returns the number of trades per page.
func (t *Trades) GetLimit() int {
	return t.pagination.Limit
}
