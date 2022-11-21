package models

import (
	"fmt"

	"github.com/cryptomarket/cryptomarket-go/args"
)

type Currency struct {
	FullName          string    `json:"full_name"`
	PayinEnabled      bool      `json:"payin_enabled"`
	PayoutEnabled     bool      `json:"payout_enabled"`
	TransferEnabled   bool      `json:"transfer_enabled"`
	PrecisionTransfer string    `json:"precision_transfer"`
	Networks          []Network `json:"networks"`
}

type Network struct {
	Network            string `json:"network"`
	Protocol           string `json:"protocol"`
	Default            bool   `json:"default"`
	PayinEnabled       bool   `json:"payin_enabled"`
	PayoutEnabled      bool   `json:"payout_enabled"`
	PrecisionPayout    string `json:"presicion_payout"`
	PayoutFee          string `str:"payout_fee"`
	PayoutIsPaymentID  bool   `json:"payout_is_payment_id"`
	PayinPaymentID     bool   `json:"payin_payment_id"`
	PayinConfirmations int    `json:"payin_confirmation"`
	AddressRegex       string `json:"address_confirmation"`
	PaymentIDRegex     string `json:"payment_id_regex"`
	LowProcessingTime  string `json:"low_processing_time"`
	HighProcessingTime string `json:"high_processing_time"`
	AvgProcessingTime  string `json:"avg_processing_time"`
}

// Balance is the amount of currency a user have
type Balance struct {
	Currency       string `json:"currency"`
	Available      string `json:"available"`
	Reserved       string `json:"reserved"`
	ReservedMargin string `json:"reserved_margin"`
}

type SubAccountBalances struct {
	Wallet []Balance `json:"wallet"`
	Spot   []Balance `json:"spot"`
}

type Price struct {
	Currency  string `json:"currency"`
	Price     string `json:"price"`
	Timestamp string `json:"timestamp"`
}

type PriceHistory struct {
	Currency string         `json:"currency"`
	History  []HistoryPoint `json:"history"`
}

type HistoryPoint struct {
	Timestamp string `json:"timestamp"`
	Open      string `json:"open"`
	Close     string `json:"close"`
	Min       string `json:"min"`
	Max       string `json:"max"`
}

// PublicTrade is the available information from public trades
type PublicTrade struct {
	ID        int64  `json:"id"`
	Price     string `json:"price"`
	Quantity  string `json:"qty"`
	Side      string `json:"side"`
	Timestamp string `json:"timestamp"`
}

// BookLevel agregates orders by price in a symbol
type BookLevel struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

// OrderBook is the current state of a symbol
type OrderBookJson struct {
	Ask       [][]string `json:"ask"`
	Bid       [][]string `json:"bid"`
	Timestamp string     `json:"timestamp"`
}

// OrderBook is the current state of a symbol
type OrderBook struct {
	Ask       []BookLevel `json:"ask"`
	Bid       []BookLevel `json:"bid"`
	Timestamp string      `json:"timestamp"`
}

// TradingCommission is the asociated cost to trade in the exchange
type TradingCommission struct {
	Symbol   string `json:"symbol"`
	TakeRate string `json:"take_rate"`
	MakeRate string `json:"make_rate"`
}

// Symbol is a market made of two currencies being exchanged
type Symbol struct {
	Type               string                `json:"type"`
	BaseCurrency       string                `json:"base_currency"`
	QuoteCurrency      string                `json:"quote_currency"`
	Status             args.SymbolStatusType `json:"status"`
	QuantityIncrement  string                `json:"quantity_increment"`
	TickSize           string                `json:"tick_size"`
	TakeRate           string                `json:"take_rate"`
	MakeRate           string                `json:"make_rate"`
	FeeCurrency        string                `json:"fee_currency"`
	MarginTrading      bool                  `json:"margin_trading"`
	MaxInitialLeverage string                `json:"max_initial_leverage"`
}

// Order is the abstraction of an order in a symbol in the exchange
type Order struct {
	ID                    int64                `json:"id"`
	ClientOrderID         string               `json:"client_order_id"`
	Symbol                string               `json:"symbol"`
	Side                  string               `json:"side"`
	Status                args.OrderStatusType `json:"status"`
	Type                  args.OrderType       `json:"type"`
	TimeInForce           args.TimeInForceType `json:"time_in_force"`
	Quantity              string               `json:"quantity"`
	QuantityCumulative    string               `json:"quantity_cumulative"`
	Price                 string               `json:"price"`
	StopPrice             string               `json:"stop_price"`
	ExpireTime            string               `json:"expire_time"`
	PostOnly              bool                 `json:"post_only"`
	OriginalClientOrderID string               `json:"original_client_order_id"`
	CreatedAt             string               `json:"created_at"`
	UpdatedAt             string               `json:"updated_at"`
	Trades                []TradeOfOrder       `json:"trades"`
	Contingency           args.ContingencyType `json:"contingency_type"`
	OrderListID           string               `json:"order_list_id"`
}

// TradeOfOrder is the trade information of trades of an order
type TradeOfOrder struct {
	ID        int64  `json:"id"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Fee       string `json:"fee"`
	Taker     bool   `json:"taker"`
	Timestamp string `json:"timestamp"`
}

// Trade is a movement of currency where the user takes part
type Trade struct {
	ID            int64         `json:"id"`
	OrderID       int64         `json:"order_id"`
	ClientOrderID string        `json:"client_order_id"`
	Symbol        string        `json:"symbol"`
	Side          args.SideType `json:"side"`
	Quantity      string        `json:"quantity"`
	Price         string        `json:"price"`
	Fee           string        `json:"fee"`
	Timestamp     string        `json:"timestamp"`
	Taker         bool          `json:"taker"`
}

// Transaction is a movement of currency,
// not in the market, but related on the exchange
type Transaction struct {
	ID        int64                       `json:"id,result"`
	Status    args.TransactionStatusType  `json:"status"`
	Type      args.TransactionTypeType    `json:"type"`
	SubType   args.TransactionSubTypeType `json:"subtype"`
	CreatedAt string                      `json:"created_at"`
	UpdatedAt string                      `json:"updated_at"`
	Native    NativeTransaction           `json:"native"`
	Meta      MetaTransaction             `json:"meta"`
}

type NativeTransaction struct {
	ID            string   `json:"tx_id"`
	Index         int64    `json:"index"`
	Currency      string   `json:"currency"`
	Amount        string   `json:"amount"`
	Fee           string   `json:"fee"`
	Address       string   `json:"address"`
	PaymentID     string   `json:"payment_id"`
	Hash          string   `json:"hash"`
	OffchainID    string   `json:"offchain_id"`
	Confirmations int64    `json:"confirmations"`
	PublicComment string   `json:"public_comment"`
	ErrorCode     string   `json:"error_code"`
	Senders       []string `json:"senders"`
}

type MetaTransaction struct {
	FiatToCrypto map[string]interface{} `json:"fiat_to_crypto"`

	ID                string                         `json:"id"`
	ProviderName      string                         `json:"provider_name"`
	OrderType         string                         `json:"order_type"`
	SourceCurrency    string                         `json:"source_currency"`
	TargetCurrency    string                         `json:"target_currency"`
	WalletAddress     string                         `json:"wallet_address"`
	TransactionHash   string                         `json:"tx_hash"`
	TargetAmount      string                         `json:"target_amount"`
	SourceAmount      string                         `json:"source_amount"`
	Status            args.MetaTransactionStatusType `json:"status"`
	CreatedAt         string                         `json:"created_at"`
	UpdatedAt         string                         `json:"updated_at"`
	DeletedAt         string                         `json:"deleted_at"`
	PaymentMethodType string                         `json:"payment_method_type"`
}

// CryptoAddress is an crypto address
type CryptoAddress struct {
	Currency  string `json:"currency"`
	Address   string `json:"address"`
	PaymentID string `json:"payment_id"`
	PublicKey string `json:"publicKey"`
}

// PayoutCryptoAddress is for external crypto addresses
type PayoutCryptoAddress struct {
	Address   string `json:"address"`
	PaymentID string `json:"payment_id"`
}

// Candle is an OHLC representation of the market
// This version uses Max instead of High nad Min instead of Low
type Candle struct {
	Timestamp   string `json:"timestamp"`
	Open        string `json:"open"`
	Close       string `json:"close"`
	High        string `json:"max"`
	Low         string `json:"min"`
	Volume      string `json:"volume"`
	VolumeQuote string `json:"volume_quote"`
}

type Ticker struct {
	Timestamp   string `json:"timestamp"`
	Open        string `json:"open"`
	Close       string `json:"last,close"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Volume      string `json:"volume"`
	VolumeQuote string `json:"volume_quote"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
}

// Error is an error from the exchange
type APIError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (err APIError) String() string {
	return fmt.Sprintf("CryptomarketAPIError: (code=%d) %s. %s", err.Code, err.Message, err.Description)
}

// ErrorMetadata is the data asociated with an error
// from the exchange
type ErrorMetadata struct {
	Timestamp string    `json:"timestamp"`
	Path      string    `json:"path"`
	APIError  *APIError `json:"error"`
	RequestID string    `json:"request_id"`
	Status    int       `json:"status"`
}

// Report is used for websocket trading reports.
type Report struct {
	ID                    int64                `json:"id"`
	ClientOrderID         string               `json:"client_order_id"`
	Symbol                string               `json:"symbol"`
	Side                  args.SideType        `json:"side"`
	Status                args.OrderStatusType `json:"status"`
	OrderType             args.OrderType       `json:"type"`
	TimeInForce           args.TimeInForceType `json:"time_in_force"`
	Quantity              string               `json:"quantity"`
	Price                 string               `json:"price"`
	QuantityCumulative    string               `json:"quantity_cumulative"`
	PostOnly              bool                 `json:"post_only"`
	OrderListID           string               `json:"order_list_id"`
	CreatedAt             string               `json:"created_at"`
	UpdatedAt             string               `json:"updated_at"`
	StopPrice             string               `json:"stopPrice"`
	ExpireTime            string               `json:"expire_time"`
	OriginalClientOrderID string               `json:"original_client_order_id"`
	TradeID               int64                `json:"trade_id"`
	TradeQuantity         string               `json:"trade_quantity"`
	TradePrice            string               `json:"trade_price"`
	TradeFee              string               `json:"trade_fee"`
	TradeTaker            bool                 `json:"trade_taker"`
	ReportType            args.ReportType      `json:"report_type"`
}

type AmountLock struct {
	ID                int64  `json:"id"`
	Currency          string `json:"currency"`
	Amount            string `json:"amount"`
	DateEnd           string `json:"date_end"`
	Description       string `json:"description"`
	Cancelled         bool   `json:"cancelled"`
	CancelledAt       string `json:"cancelled_at"`
	CancelDescription string `json:"cancel_description"`
	CreatedAt         string `json:"created_at"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type ResultResponse struct {
	ID string `json:"result"`
}

type ResultListResponse struct {
	IDs []string `json:"result"`
}

type BooleanResponse struct {
	Result bool `json:"result"`
}

type FeeResponse struct {
	Fee string `json:"fee"`
}

type SubAccount struct {
	ID     string `json:"sub_account_id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type ACLSettings struct {
	SubAccountID                    string `json:"sub_account_id"`
	DepositAddressGenerationEnabled bool   `json:"deposit_address_generation_enabled"`
	WithdrawEnabled                 bool   `json:"withdraw_enabled"`
	Description                     string `json:"description"`
	CreatedAt                       string `json:"created_at"`
	UpdatedAt                       string `json:"updated_at"`
}
