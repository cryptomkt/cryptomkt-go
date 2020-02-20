package conn

import (
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/args"
)

type PricesResponse struct {
	Status     string
	Message    string
	Pagination Pagination
	Data       DataPrices
}

type Prices struct {
	args       map[string]string
	pagination Pagination
	client     *Client
	Data       DataPrices
}

type DataPrices struct {
	Ask []Candle
	Bid []Candle
}

type Candle struct {
	CandleId   int    `json:"candle_id"`
	OpenPrice  string `json:"open_price"`
	HightPrice string `json:"hight_price"`
	ClosePrice string `json:"close_price"`
	LowPrice   string `json:"low_price"`
	VolumeSum  string `json:"volume_sum"`
	CandleDate string `json:"candle_date"`
	TickCount  string `json:"tick_count"`
}

// GetPrevious lets you go to the previous page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (p *Prices) GetPrevious() (*Prices, error) {
	if p.pagination.Previous == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	return p.client.GetPrices(
		args.Market(p.args["market"]),
		args.Type(p.args["timeframe"]),
		args.Page(p.pagination.Previous.(int)),
		args.Limit(p.pagination.Limit))
}

// GetNext lets you go to the next page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (p *Prices) GetNext() (*Prices, error) {
	if p.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	return p.client.GetPrices(
		args.Market(p.args["market"]),
		args.Type(p.args["timeframe"]),
		args.Page(p.pagination.Next.(int)),
		args.Limit(p.pagination.Limit))
}

// GetPage returns the page you have
func (p *Prices) GetPage() int {
	return p.pagination.Page
}

// GetLimit returns the limit you have provided, but if you have not, it provides the default
func (p *Prices) GetLimit() int {
	return p.pagination.Limit
}
