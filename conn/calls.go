package conn

import (
	"encoding/json"
	"fmt"

	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/requests"
)

// GetAccount gives the account information of the client.
//
// https://developers.cryptomkt.com/es/#cuenta
func (client *Client) GetAccount() (*Account, error) {
	resp, err := client.get("account", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var aResp AccountResponse
	json.Unmarshal(resp, &aResp)
	if aResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", aResp.Message)
	}
	return &aResp.Data, nil
}

// GetBalance returns the actual balance of the wallets of the client in CryptoMarket
//
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) GetBalance() (*[]Balance, error) {
	resp, err := client.get("balance", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var bResp BalancesResponse
	json.Unmarshal(resp, &bResp)
	if bResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", bResp.Message)
	}
	return &bResp.Data, nil
}

// GetWallets is an alias for Balance, returns the actual balance of wallets of the client in CryptoMarket
//
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) GetWallets() (*[]Balance, error) {
	return client.GetBalance()
}

// GetTransactions returns the movements of the wallets of the client.
//
// List of accepted Arguments:
//   - required: Currency
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#obtener-movimientos
func (client *Client) GetTransactions(arguments ...args.Argument) (*[]Transaction, error) {
	resp, err := client.getReq("transactions", "GetTransaction", []string{"currency"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var tResp TransactionsResponse
	json.Unmarshal(resp, &tResp)
	if tResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", tResp.Message)
	}
	return &tResp.Data, nil
}

// GetActiveOrders returns the list of active orders of the client
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-activas
func (client *Client) GetActiveOrders(arguments ...args.Argument) (*OrderList, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetActiveOrders: %s", err)
	}
	resp, err := client.post("orders/active", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var oListResp OrderListResp
	json.Unmarshal(resp, &oListResp)
	if oListResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", oListResp.Message)
	}
	orderList := OrderList{
		pagination: oListResp.Pagination,
		client:     client,
		Data:       oListResp.Data,
		caller:     "active_orders",
		market:     req.GetArguments()["market"],
	}
	return &orderList, nil
}

// GetExecutedOrders return a list of the executed orders of the client
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-ejecutadas
func (client *Client) GetExecutedOrders(arguments ...args.Argument) (*OrderList, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetExecutedOrders: %s", err)
	}
	resp, err := client.post("orders/executed", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var oListResp OrderListResp
	json.Unmarshal(resp, &oListResp)
	if oListResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", oListResp.Message)
	}

	orderList := OrderList{
		pagination: oListResp.Pagination,
		client:     client,
		Data:       oListResp.Data,
		caller:     "executed_orders",
		market:     req.GetArguments()["market"],
	}
	return &orderList, nil
}

// GetOrderStatus gives the status of an order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden
func (client *Client) GetOrderStatus(arguments ...args.Argument) (*Order, error) {
	resp, err := client.getReq("orders/status", "GetOrderStatus", []string{"id"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var oResp OrderResponse
	json.Unmarshal(resp, &oResp)
	if oResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", oResp.Message)
	}
	oResp.Data.client = client
	return &oResp.Data, nil
}

// GetInstant emulates an order in the current state of the Instant Exchange of CryptoMarket
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#obtener-cantidad
func (client *Client) GetInstant(arguments ...args.Argument) (*Quantity, error) {
	resp, err := client.getReq("orders/instant/get", "GetInstant", []string{"market", "type", "amount"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", iResp.Message)
	}
	return &iResp.Data, nil
}

// CreateOrder signal to create an order of buy or sell in CryptoMarket
//
// List of accepted Arguments:
//   - required: Amount, Market, Price, Type
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden
func (client *Client) CreateOrder(arguments ...args.Argument) (*Order, error) {
	resp, err := client.postReq("orders/create", "CreateOrder", []string{"amount", "market", "price", "type"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var oResp OrderResponse
	json.Unmarshal(resp, &oResp)
	if oResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", oResp.Message)
	}
	oResp.Data.client = client
	return &oResp.Data, nil
}

// CancelOrder signal to cancel an order in CryptoMarket
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#cancelar-una-orden
func (client *Client) CancelOrder(arguments ...args.Argument) (*Order, error) {
	resp, err := client.postReq("orders/cancel", "CancelOrder", []string{"id"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var oResp OrderResponse
	json.Unmarshal(resp, &oResp)
	if oResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", oResp.Message)
	}
	oResp.Data.client = client
	return &oResp.Data, nil
}

// CreateInstant makes an order in the Instant Exchange of CryptoMarket
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden-2
func (client *Client) CreateInstant(arguments ...args.Argument) error {
	resp, err := client.postReq("orders/instant/create", "CreateInstant", []string{"market", "type", "amount"}, arguments...)
	if err != nil {
		return fmt.Errorf("error making the request: %s", err)
	}
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from the server side: %s", iResp.Message)
	}
	return nil
}

// RequestDeposit notifies a deposit to a wallet of local currency
//
// List of accepted Arguments:
//   - required: Amount, BankAccount
//   - required only for México, Brasil and European Union: Voucher
//   - required only for México: Date, TrackingCode
// https://developers.cryptomkt.com/es/#notificar-deposito
func (client *Client) RequestDeposit(arguments ...args.Argument) error {
	resp, err := client.postReq("request/deposit", "RequestDeposit", []string{"amount", "bank_account"}, arguments...)
	if err != nil {
		return fmt.Errorf("error making the request: %s", err)
	}
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from the server side: %s", iResp.Message)
	}
	return nil
}

// RequestWithdrawal notifies a withdrawal from a bank account of the client
//
// List of accepted Arguments:
//   - required: Amount, BankAccount
//   - optional: none
// https://developers.cryptomkt.com/es/#notificar-retiro
func (client *Client) RequestWithdrawal(arguments ...args.Argument) error {
	resp, err := client.postReq("request/withdrawal", "RequestWithdrawal", []string{"amount", "bank_account"}, arguments...)
	if err != nil {
		return fmt.Errorf("error making the request: %s", err)
	}
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from the server side: %s", iResp.Message)
	}
	return nil
}

// Transfer move crypto between wallets
//
// List of accepted Arguments:
//   - required: Address, Amount, Currency
//   - optional: Memo
// https://developers.cryptomkt.com/es/#transferir
func (client *Client) Transfer(arguments ...args.Argument) error {
	resp, err := client.postReq("transfer", "Transfer", []string{"address", "amount", "currency"}, arguments...)
	if err != nil {
		return fmt.Errorf("error making the request: %s", err)
	}
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from the server side: %s", iResp.Message)
	}
	return nil

}

// NewOrder enables a payment order, and gives a QR and urls
//
// List of accepted Arguments:
//   - required: ToReceive, ToReceiveCurrency, PaymentReceiver
//   - optional: ExternalId, CallbackUrl, ErrorUrl, SuccessUrl, RefundEmail, Language
// https://developers.cryptomkt.com/es/#crear-orden-de-pago
func (client *Client) NewOrder(arguments ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/new_order", "NewOrder", []string{"to_receive", "to_receive_currency", "payment_receiver"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var poResp PaymentResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
}

// CreateWallet creates a wallet to pay a payment order
//
// List of accepted Arguments:
//   - required: Id, Token, Wallet
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-billetera-de-orden-de-pago
func (client *Client) CreateWallet(arguments ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/create_wallet", "CreateWallet", []string{"id", "token", "wallet"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var poResp PaymentResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
}

// PaymentOrders returns all the generated payment orders
//
// List of accepted Arguments:
//   - required: StartDate, EndDate
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#listado-de-ordenes-de-pago
func (client *Client) PaymentOrders(arguments ...args.Argument) (*PaymentOrderList, error) {
	req, err := makeReq([]string{"start_date", "end_date"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in PaymentOrders: %s", err)
	}
	resp, err := client.post("payment/orders", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var poResp PaymentOrdersResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", poResp.Message)
	}
	argMap := req.GetArguments()
	paymentOrderList := PaymentOrderList{
		startDate:  argMap["start_date"],
		endDate:    argMap["end_date"],
		client:     client,
		pagination: poResp.Pagination,
		Data:       poResp.Data,
	}
	return &paymentOrderList, nil
}

// GetPaymentStatus gives the status of a payment order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden-de-pago
func (client *Client) GetPaymentStatus(arguments ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/status", "PaymentStatus", []string{"id"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var poResp PaymentResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
}

// Public Endpoints:

//

// GetMarkets returns a pointer to a MarketStruct with the field "data" given by the api. The data given is
// an array of strings indicating the markets in cryptomkt. This function returns two values.
// The first is a reference to the struct created and the second is a error message. It returns (nil, error)
// when an error is raised.
// This method does not accept any arguments.
// https://developers.cryptomkt.com/es/mercado
func (client *Client) GetMarkets() ([]string, error) {
	resp, err := client.get("market", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}

	var mResp MarketListResponse
	json.Unmarshal(resp, &mResp)
	if mResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", mResp.Message)
	}
	return mResp.Data, nil
}

// GetTicker returns a pointer to a Ticker struct with the data given by the api and an error message. It returns (nil,error)
//when an error is raised and (*Ticker, nil) when the operation is successful. The data fields are: High, Low, Ask, Bid,
//LastPrice, Volume, Market and Timestamp
//
// List of accepted Arguments:
//   - required: none
//   - optional: Market
// https://developers.cryptomkt.com/es/#ticker
func (client *Client) GetTicker(arguments ...args.Argument) (*[]Ticker, error) {
	resp, err := client.getReq("ticker", "GetTicker", []string{}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var tResp TickerResponse
	json.Unmarshal(resp, &tResp)
	if tResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", tResp.Message)
	}
	return &tResp.Data, nil
}

// GetBook returns a pointer to a Book struct with the data given by
// the api and an error message. It returns (nil, error) when an error
// is raised and (*Book, nil) when the operation is successful.
// The data fields are: price, amount and timestamp.
//
// List of accepted Arguments:
//   - required: Market, Type
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes
func (client *Client) GetBook(arguments ...args.Argument) (*Book, error) {
	req, err := makeReq([]string{"market", "type"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetBook: %s", err)
	}
	resp, err := client.getPublic("book", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var bResp BookResponse
	json.Unmarshal(resp, &bResp)
	if bResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", bResp.Message)
	}
	book := Book{
		args:       req.GetArguments(),
		pagination: bResp.Pagination,
		client:     client,
		Data:       bResp.Data,
	}
	return &book, nil
}

// GetTrades returns a pointer to a Trades struct with the data given
// by the api and an error message. It returns (nil, error) when an error
// is raised and (*Trades, nil) when the operation is successful.
// The data fields are market_taker, price, amount, tid, timestamp and market.
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Start, End, Page, Limit
// https://developers.cryptomkt.com/es/#trades
func (client *Client) GetTrades(arguments ...args.Argument) (*Trades, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTradesPage: %s", err)
	}
	resp, err := client.getPublic("trades", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var tResp TradesResponse
	json.Unmarshal(resp, &tResp)
	if tResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", tResp.Message)
	}
	trades := Trades{
		args:       req.GetArguments(),
		pagination: tResp.Pagination,
		client:     client,
		Data:       tResp.Data,
	}
	return &trades, nil
}

// GetPrices return a pointer to a Prices struct with the data given by
// the api and an error message. It returns (nil,error) when an error
// is raised and (*Prices, nil) when the operation is successful.
//
// List of accepted Arguments:
//   - required: Market, Timeframe
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#precios
func (client *Client) GetPrices(arguments ...args.Argument) (*Prices, error) {
	req, err := makeReq([]string{"market", "timeframe"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetPrices: %s", err)
	}
	resp, err := client.getPublic("prices", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var pResp PricesResponse
	json.Unmarshal(resp, &pResp)
	if pResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", pResp.Message)
	}
	prices := Prices{
		args:       req.GetArguments(),
		pagination: pResp.Pagination,
		client:     client,
		Data:       pResp.Data,
	}
	return &prices, nil
}