package conn

type MarketStruct struct {
	Data []string
}

type DataTicker struct {
	High      string
	Volume    string
	Low       string
	Ask       string
	Timestamp string
	Bid       string
	LastPrice string `json:"last_price"`
	Market    string
}

type TemporalTicker struct {
	Status string
	Data   []DataTicker
}

type Ticker struct {
	Data []DataTicker
}

type BookData struct {
	Price     string
	Amount    string
	Timestamp string
}

type TemporalBook struct {
	Status     string
	Pagination Pagination
	Data       []BookData
}

type Book struct {
	args       map[string]string
	pagination Pagination
	client     *Client
	Data       []BookData
}

<<<<<<< HEAD
type TradesData struct {
	MarketTaker string `json:"market_taker"`
	Price       string
	Amount      string
	Tid         string
	Timestamp   string
	Market      string
=======
type Book struct {
	Data []map[string]string
>>>>>>> 77cc5aa8f6445fba169e90a878e8bc0b3c0734e4
}
type TemporalTrades struct {
	Status     string
	Pagination Pagination
	Data       []TradesData
}

type Trades struct {
	args       map[string]string
	pagination Pagination
	client     *Client
	Data       []TradesData
}

type TemporalPrices struct {
	Status     string
	Pagination Pagination
	Data       DataPrices1
}

type Prices struct {
	args       map[string]string
	pagination Pagination
	client     *Client
	Data       DataPrices1
}

type DataPrices1 struct {
	Ask []DataPrices2
	Bid []DataPrices2
}

type DataPrices2 struct {
	CandleId   int    `json:"candle_id"`
	OpenPrice  string `json:"open_price"`
	HightPrice string `json:"hight_price"`
	ClosePrice string `json:"close_price"`
	LowPrice   string `json:"low_price"`
	VolumeSum  string `json:"volume_sum"`
	CandleDate string `json:"candle_date"`
	TickCount  string `json:"tick_count"`
}
type Pagination struct {
	Previous interface{}
	Limit    int
	Page     int
	Next     interface{}
}
