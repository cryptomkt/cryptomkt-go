package conn

type MarketStruct struct {
	Data []string
}

type Ticker struct {
	Data []map[string]string
}

type Order struct {
	Data []map[string]string
}
type Trades struct {
	Data []map[string]string
}
type Prices struct {
	Data map[string][]Field
}

type Field struct {
	Candle_id   int
	Open_price  string
	Hight_price string
	Close_price string
	Low_price   string
	Volume_sum  string
	Candle_date string
	Tick_count  string
}
