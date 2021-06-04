package models

// SideType is the side of the order or trade
type SideType string

// OrderStatus is the status of an order in the exchange
type OrderStatus string

// OrderType is the type of an order
type OrderType string

// TimeInForceType is the TimeInForce of an order
type TimeInForceType string

// TransactionStatus is the status of transaction
type TransactionStatus string

// TransactionType is the type of a transaction
type TransactionType string

// TransactionSubType is the sub type of a transaction
type TransactionSubType string

// ReportType shows the type of report
type ReportType string

// sides types of order and trades
const (
	SideTypeSell SideType = "sell"
	SideTypeBuy  SideType = "buy"

	OrderStatusNew             OrderStatus = "new"
	OrderStatusSuspended       OrderStatus = "suspended"
	OrderStatusPartiallyFilled OrderStatus = "partiallyFilled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCanceled        OrderStatus = "canceled"
	OrderStatusExpired         OrderStatus = "expired"

	OrderTypeLimit      OrderType = "limit"
	OrderTypeMarket     OrderType = "market"
	OrderTypeStopLimit  OrderType = "stopLimit"
	OrderTypeStopMarket OrderType = "stopMarket"

	TimeInForceTypeGTC TimeInForceType = "GTC"
	TimeInForceTypeIOC TimeInForceType = "IOC"
	TimeInForceTypeFOK TimeInForceType = "FOK"
	TimeInForceTypeDAY TimeInForceType = "DAY"
	TimeInForceTypeGTD TimeInForceType = "GTD"

	TransactionStatusCreated TransactionStatus = "created"
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusFailed  TransactionStatus = "failed"
	TransactionStatusSuccess TransactionStatus = "success"

	TransactionTypePayout         TransactionType = "payout"
	TransactionTypePayin          TransactionType = "payin"
	TransactionTypeDeposit        TransactionType = "deposit"
	TransactionTypeWithdraw       TransactionType = "withdraw"
	TransactionTypeBankToExchange TransactionType = "bankToExchange"
	TransactionTypeExchangeToBank TransactionType = "exchangeToBank"

	TransactionSubTypeOffchain TransactionSubType = "offchain"
	TransactionSubTypeSwap     TransactionSubType = "swap"

	ReportTypeStatus    ReportType = "status"
	ReportTypeNew       ReportType = "new"
	ReportTypeCanceled  ReportType = "canceled"
	ReportTypeExpired   ReportType = "expired"
	ReportTypeSuspended ReportType = "suspended"
	ReportTypeTrade     ReportType = "trade"
	ReportTypeReplaced  ReportType = "replaced"
)

// Currency is the abstraction for a digital currency
type Currency struct {
	ID                  string `json:"id"`
	FullName            string `json:"fullName"`
	Crypto              bool   `json:"crypto"`
	PayinEnabled        bool   `json:"payinEnabled"`
	PayinPaymentID      bool   `json:"payinPaymentId"`
	PayinConfirmations  int    `json:"payinConfirmations"`
	PayoutEnabled       bool   `json:"payoutEnable"`
	PayoutFee           string `json:"payoutFee"`
	PayoutIsPaymentID   bool   `json:"payoutIsPaymentId"`
	Delisted            bool   `json:"delisted"`
	TransferEnabled     bool   `json:"transferEnabled"`
	PayoutMinimalAmount string `json:"payoutMinimalAmount"`
	PrecisionPayout     int    `json:"precisionPayout"`
	PrecisionTransfer   int    `json:"precisionTransfer"`
	LowProcessingTime   string `json:"lowProcessingTime"`
	HighProcessingTime  string `json:"highProcessingTime"`
	AvgProcessingTime   string `json:"avgProcessingTime"`
}

// Balance is the amount of currency a user have
type Balance struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Reserved  string `json:"reserved"`
}

// Ticker is a snapshot of a symbol
type Ticker struct {
	Symbol       string `json:"symbol"`
	Ask          string `json:"ask"`
	Bid          string `json:"bid"`
	Last         string `json:"last"`
	Low          string `json:"low"`
	High         string `json:"high"`
	Open         string `json:"open"`
	Volume       string `json:"volume"`
	VolumeQuote string `json:"volumeQuote"`
	Timestamp    string `json:"timestamp"`
}

// PublicTrade is the available information from public trades
type PublicTrade struct {
	ID        int64    `json:"id"`
	Price     string   `json:"price"`
	Quantity  string   `json:"quantity"`
	Side      SideType `json:"side"`
	Timestamp string   `json:"timestamp"`
}

// BookLevel agregates orders by price in a symbol
type BookLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// OrderBook is the current state of a symbol
type OrderBook struct {
	Symbol          string      `json:"symbol"`
	Ask             []BookLevel `json:"ask"`
	Bid             []BookLevel `json:"bid"`
	Timestamp       string      `json:"timestamp"`
	BatchingTime    string      `json:"batchingTime"`
	AskAveragePrice string      `json:"askAveragePrice"`
	BidAveragePrice string      `json:"bidAveragePrice"`
}

// TradingFee is the asociated cost to trade in the exchange
type TradingFee struct {
	TakeLiquidityRate    string `json:"takeLiquidityRate"`
	ProvideLiquidityRate string `json:"provideLiquidityRate"`
}

// Symbol is a market made of two currencies being exchanged
type Symbol struct {
	ID                   string `json:"id"`
	BaseCurrency         string `json:"baseCurrency"`
	QuoteCurrency        string `json:"quoteCurrency"`
	QuantityIncrement    string `json:"quantityIncrement"`
	TickSize             string `json:"tickSize"`
	TakeLiquidityRate    string `json:"takeLiquidityRate"`
	ProvideLiquidityRate string `json:"provideLiquidityRate"`
	FeeCurrency          string `json:"feeCurrency"`
}

// Order is the abstraction of an order in a symbol in the exchange
type Order struct {
	ID            int64           `json:"id"`
	ClientOrderID string          `json:"clientOrderId"`
	Symbol        string          `json:"symbol"`
	Side          SideType        `json:"side"`
	Status        OrderStatus     `json:"status"`
	Type          OrderType       `json:"type"`
	TimeInForce   TimeInForceType `json:"timeInForce"`
	ExpireTime    string          `json:"expireTime"`
	Quantity      string          `json:"quantity"`
	Price         string          `json:"price"`
	StopPrice     string          `json:"stopPrice"`
	AvgPrice      string          `json:"avgPrice"`
	PostOnly      bool            `json:"postOnly"`
	CumQuantity   string          `json:"cumQuantity"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	PositionID    string          `json:"positionId"`
	TradesReport  []TradeReport   `json:"tradesReport"`
}

// TradeReport is the trade information of trades of an order
type TradeReport struct {
	ID        int64  `json:"id"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Fee       string `json:"fee"`
	Timestamp string `json:"timestamp"`
}

// Trade is a movement of currency where the user takes part
type Trade struct {
	ID            int64    `json:"id"`
	ClientOrderID string   `json:"clientOrderId"`
	OrderID       int64    `json:"orderId"`
	Symbol        string   `json:"symbol"`
	Side          SideType `json:"side"`
	Quantity      string   `json:"quantity"`
	Fee           string   `json:"fee"`
	Price         string   `json:"price"`
	PositionID    string   `json:"positionId"`
	Pnl           string   `json:"pnl"`
	Timestamp     string   `json:"timestamp"`
	Liquidation   string   `json:"liquidation"`
}

// Transaction is a movement of currency,
// not in the market, but related on the exchange
type Transaction struct {
	ID            string             `json:"id,result"`
	Index         int64              `json:"index"`
	Currency      string             `json:"currency"`
	Amount        string             `json:"amount"`
	Fee           string             `json:"fee"`
	Address       string             `json:"address"`
	PaymentID     string             `json:"paymentId"`
	Hash          string             `json:"hash"`
	Status        TransactionStatus  `json:"status"`
	Type          TransactionType    `json:"type"`
	SubType       TransactionSubType `json:"subType"`
	OffchainID    string             `json:"offchainId"`
	Confirmations int64              `json:"confirmations"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
	PublicComment string             `json:"publicComment"`
	ErrorCode     string             `json:"errorCode"`
}

// CryptoAddress is an crypto address
type CryptoAddress struct {
	Address   string `json:"address"`
	PaymentID string `json:"paymentId"`
	PublicKey string `json:"publicKey"`
}

// PayoutCryptoAddress is for external crypto addresses
type PayoutCryptoAddress struct {
	Address   string `json:"address"`
	PaymentID string `json:"paymentId"`
}

// Candle is an OHLC representation of the market
// This version uses Max instead of High nad Min instead of Low
type Candle struct {
	Timestamp   string `json:"timestamp"`
	Open        string `json:"open"`
	Close       string `json:"close"`
	Min         string `json:"min"`
	Max         string `json:"max"`
	Volume      string `json:"volume"`
	VolumeQuote string `json:"volumeQuote"`
}

// Error is an error from the exchange
type Error struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// ErrorMetadata is the data asociated with an error
// from the exchange
type ErrorMetadata struct {
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Error     *Error `json:"error"`
	RequestID string `json:"requestId"`
	Status    int    `json:"status"`
}

// Report is used for websocket trading reports.
type Report struct {
	ID                           int64           `json:"id"`
	ClientOrderID                string          `json:"clientOrderId"`
	Symbol                       string          `json:"symbol"`
	Side                         SideType        `json:"side"`
	Status                       OrderStatus     `json:"status"`
	Type                         OrderType       `json:"type"`
	TimeInForce                  TimeInForceType `json:"timeInForce"`
	ExpireTime                   string          `json:"expireTime"`
	Quantity                     string          `json:"quantity"`
	Price                        string          `json:"price"`
	StopPrice                    string          `json:"stopPrice"`
	PostOnly                     bool            `json:"postOnly"`
	CumQuantity                  string          `json:"cumQuantity"`
	CreatedAt                    string          `json:"createdAt"`
	UpdatedAt                    string          `json:"updatedAt"`
	ReportType                   ReportType      `json:"reportType"`
	TradeID                      int64           `json:"tradeId"`
	TradeQuantity                string          `json:"tradeQuantity"`
	TradePrice                   string          `json:"tradePrice"`
	TradeFee                     string          `json:"tradeFee"`
	OriginalRequestClientOrderID string          `json:"originalRequestClientOrderId"`
}
