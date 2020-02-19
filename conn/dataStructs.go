package conn

type MarketStruct struct {
	Data []string
}

type Ticker struct {
	Data []map[string]string
}

type Book struct {
	Data []map[string]string
}
type Trades struct {
	Data []map[string]string
}
type PricesResponse struct {
	status  string
	Message string
	data    []Prices
}

type Prices struct {
	CandleId   int    `json:"candle_id"`
	OpenPrice  string `json:"open_price"`
	HightPrice string `json:"hight_price"`
	ClosePrice string `json:"close_price"`
	LowPrice   string `json:"low_price"`
	VolumeSum  string `json:"volume_sum"`
	CandleDate string `json:"candle_date"`
	TickCount  string `json:"tick_count"`
}

type Balance struct {
	Wallet    string
	Available string
	Balance   string
}

type BalancesResponse struct {
	Status  string
	Message string
	Data    []Balance
}

type PaymentOrder struct {
	Id                string
	ExternalId        string `json:"external_id"`
	Status            int
	ToReceive         string `json:"to_receive"`
	ToReceiveCurrency string `json:"to_receive_currency"`
	ExpectedAmount    string `json:"expected_amount"`
	ExpectedCurrency  string `json:"expected_currency"`
	DepositAddress    string `json:"deposit_address"`
	RefundEmail       string `json:"refund_email"`
	Qr                string
	Obs               string
	CallbackUrl       string `json:"callback_url"`
	ErrorUrl          string `json:"error_url"`
	SuccessUrl        string `json:"success_url"`
	PaymentUrl        string `json:"payment_url"`
	Remaining         int    `json:"remanining"`
	Language          string
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ServerAt          string `json:"server_at"`
}

type PaymentResponse struct {
	Status  string
	Message string
	Data    PaymentOrder
}

type PaymentOrdersResponse struct {
	client * Client
	Status     string
	Message    string
	Pagination Pagination
	Data       []PaymentOrder
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

type AccountResponse struct {
	Status  string
	Message string
	Data    Account
}

type Quantity struct {
	Obtained string
	Required string
}

type InstantResponse struct {
	Status  string
	Message string
	Data    Quantity
}

type Transaction struct {
	Id         string
	Type       int
	Amount     string
	FeePercent string `json:"fee_percent"`
	FeeAmount  string `json:"fee_amount"`
	Balance    string
	Date       string
	Hash       string
	Address    string
	Memo       string
}

type TransactionsResponse struct {
	client *Client
	Status     string
	Message    string
	Pagination Pagination
	Data       []Transaction
}

type Pagination struct {
	PreviousHolder interface{} `json:"previous"`
	NextHolder     interface{} `json:"next"`
	Previous       int
	Next           int
	Limit          int
	Page           int
}

type Amount struct {
	Original  string
	Remaining string
	Executed  string
}

type OrderListResp struct {
	Client  *Client
	Status     string
	Message    string
	Pagination Pagination
	Warnings   string
	Data       []Order
}

type OrderResponse struct {
	Status  string
	Message string
	Data    Order
}

type Order struct {
	client         *Client
	Id                string
	Status            string
	Type              string
	Price             string
	Amount            Amount
	ExecutionPrice    string `json:"execution_price"`
	AvgExecutionPrice int    `json:"avg_execution_price"`
	Market            string
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ExecutedAt        string `json:"executed_at"`
}
