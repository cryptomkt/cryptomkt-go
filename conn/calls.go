package conn

import (
	"encoding/json"
	"fmt"

	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/requests"
)

// GetAccount gives the information of the cryptoMarket account.
// Returns the data in an Account struct
//
// https://developers.cryptomkt.com/#cuenta
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
// Returns the a slice of Balance structs
//
// https://developers.cryptomkt.com/#obtener-balance
func (client *Client) GetBalance() ([]Balance, error) {
	resp, err := client.get("balance", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var bResp BalancesResponse
	json.Unmarshal(resp, &bResp)
	if bResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", bResp.Message)
	}
	return bResp.Data, nil
}

// GetWallets is an alias for Balance, returns the actual balance of wallets of the client in CryptoMarket
//
// https://developers.cryptomkt.com/#obtener-balance
func (client *Client) GetWallets() ([]Balance, error) {
	return client.GetBalance()
}

// GetTransactions returns the movements of the wallets of the client for a given currency.
// Returns a TransactionList struct, where all the transactions are in the Data field
// in a slice of Transaction. TransactionList supports Next() and Previous to get the
// corresponding pages.
//
// List of accepted Arguments:
//   - required: Currency (string)
//   - optional: Page (int), Limit (int)
// https://developers.cryptomkt.com/#obtener-movimientos
func (client *Client) GetTransactions(arguments ...args.Argument) (*TransactionList, error) {
	required := []string{"currency"}
	req, err := makeReq(required, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTransactions: %s", err)
	}
	resp, err := client.post("transactions", req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var tResp TransactionsResponse
	json.Unmarshal(resp, &tResp)
	if tResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", tResp.Message)
	}
	tList := TransactionList{
		currency:   req.GetArguments()["currency"],
		client:     client,
		pagination: tResp.Pagination,
		Data:       tResp.Data,
	}
	return &tList, nil
}

// GetActiveOrders returns the list of active orders of the client in a given market.
// Retunrs an OrderList struct, where all the orders are in the Data field, in a slice of Order.
// OrderLists supports Next() and Previous() to get the corresponding pages.
// OrderLists also supports Close() and Refresh(), to close or refresh all the orders of the list.
//
// List of accepted Arguments:
//   - required: Market (string)
//   - optional: Page (int), Limit (int)
// https://developers.cryptomkt.com/#ordenes-activas
func (client *Client) GetActiveOrders(arguments ...args.Argument) (*OrderList, error) {
	required := []string{"market"}
	req, err := makeReq(required, arguments...)
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

// GetExecutedOrders return a list of the executed orders of the client in a given market.
// Retunrs an OrderList struct, where all the orders are in the Data field, in a slice of Order.
// OrderLists supports Next() and Previous() to get the corresponding pages.
// OrderLists also supports Close() and Refresh(), to close or refresh all the orders of the list.
//
// List of accepted Arguments:
//   - required: Market (string)
//   - optional: Page (int), Limit (int)
// https://developers.cryptomkt.com/#ordenes-ejecutadas
func (client *Client) GetExecutedOrders(arguments ...args.Argument) (*OrderList, error) {
	required := []string{"market"}
	req, err := makeReq(required, arguments...)
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

// GetOrderStatus gives the status of an order given its id.
// Returns an Order struct that supports Close() and Refresh() to cancel or update
// the order respectively.
//
// List of accepted Arguments:
//   - required: Id (string)
//   - optional: none
// https://developers.cryptomkt.com/#estado-de-orden
func (client *Client) GetOrderStatus(arguments ...args.Argument) (*Order, error) {
	required := []string{"id"}
	resp, err := client.getReq("orders/status", "GetOrderStatus", required, arguments...)
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
// Returns a Quantity struct holding the data.
//
// List of accepted Arguments:
//   - required: Market (string), Type (string), Amount (string)
//   - optional: none
// https://developers.cryptomkt.com/#obtener-cantidad
func (client *Client) GetInstant(arguments ...args.Argument) (*Quantity, error) {
	required := []string{"market", "type", "amount"}
	resp, err := client.getReq("orders/instant/get", "GetInstant", required, arguments...)
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

// CreateOrder creates an order to buy or sell in a market of CryptoMarket
// Returns an Order struct that supports Close() and Refresh() to cancel or update
// the order respectively.
//
// List of accepted Arguments:
//   - required: Amount (string), Market (string), Price (string), Type (string)
//   - optional: none
// https://developers.cryptomkt.com/#crear-orden
func (client *Client) CreateOrder(arguments ...args.Argument) (*Order, error) {
	required := []string{"amount", "market", "price", "type"}
	resp, err := client.postReq("orders/create", "CreateOrder", required, arguments...)
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

// CancelOrder cancel an order in CryptoMarket, given its id
// Returns an Order struct that supports Close() and Refresh() to cancel or update
// the order respectively.
//
// List of accepted Arguments:
//   - required: Id (string)
//   - optional: none
// https://developers.cryptomkt.com/#cancelar-una-orden
func (client *Client) CancelOrder(arguments ...args.Argument) (*Order, error) {
	required := []string{"id"}
	resp, err := client.postReq("orders/cancel", "CancelOrder", required, arguments...)
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
// Returns an error if something goes wrong
//
// List of accepted Arguments:
//   - required: Market (string), Type (string), Amount (string)
//   - optional: none
// https://developers.cryptomkt.com/#crear-orden-2
func (client *Client) CreateInstant(arguments ...args.Argument) error {
	required := []string{"market", "type", "amount"}
	resp, err := client.postReq("orders/instant/create", "CreateInstant", required, arguments...)
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

// RequestDeposit makes a deposit to a wallet of local currency
// Returns an error if something goes wrong
//
// List of accepted Arguments:
//   - required: Amount (string), BankAccount (string)
//   - required only for México, Brasil and European Union: Voucher (file)
//   - required only for México: Date (string dd/mm/yyyy), TrackingCode (string)
// https://developers.cryptomkt.com/#notificar-deposito
func (client *Client) RequestDeposit(arguments ...args.Argument) error {
	required := []string{"amount", "bank_account"}
	resp, err := client.postReq("request/deposit", "RequestDeposit", required, arguments...)
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

// RequestWithdrawal makes a withdrawal from a bank account of the client
// Returns an error if something goes wrong
//
// List of accepted Arguments:
//   - required: Amount (string), BankAccount (string)
//   - optional: none
// https://developers.cryptomkt.com/#notificar-retiro
func (client *Client) RequestWithdrawal(arguments ...args.Argument) error {
	required := []string{"amount", "bank_account"}
	resp, err := client.postReq("request/withdrawal", "RequestWithdrawal", required, arguments...)
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

// Transfer moves crypto between wallets
// Returns an error if something goes wrong
//
// List of accepted Arguments:
//   - required: Address (string), Amount (string) , Currency (string)
//   - optional: Memo (string)
// https://developers.cryptomkt.com/#transferir
func (client *Client) Transfer(arguments ...args.Argument) error {
	required := []string{"address", "amount", "currency"}
	resp, err := client.postReq("transfer", "Transfer", required, arguments...)
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

// NewOrder creates a payment order, and gives a QR and urls to make the payment.
// Returns a PaymentOrder struct that supports Refresh(), to update its information.
//
// List of accepted Arguments:
//   - required: ToReceive (double), ToReceiveCurrency (string), PaymentReceiver (string)
//   - optional: ExternalId (string), CallbackUrl (string) , ErrorUrl (string), SuccessUrl (string), RefundEmail (string) , Language (string)
// https://developers.cryptomkt.com/#crear-orden-de-pago
func (client *Client) NewOrder(arguments ...args.Argument) (*PaymentOrder, error) {
	required := []string{"to_receive", "to_receive_currency", "payment_receiver"}
	resp, err := client.postReq("payment/new_order", "NewOrder", required, arguments...)
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

// CreateWallet creates a wallet to pay a payment order.
// Returns a PaymentOrder struct that supports Refresh(), to update its information.
//
// List of accepted Arguments:
//   - required: Id (string), Token (string), Wallet (string)
//   - optional: none
// https://developers.cryptomkt.com/#crear-billetera-de-orden-de-pago
func (client *Client) CreateWallet(arguments ...args.Argument) (*PaymentOrder, error) {
	required := []string{"id", "token", "wallet"}
	resp, err := client.postReq("payment/create_wallet", "CreateWallet", required, arguments...)
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

// GetPaymentOrders returns the generated payment orders.
// Returns a PaymentOrder struct that supports Refresh(), to update its information.
//
// List of accepted Arguments:
//   - required: StartDate (string dd/mm/aaaa), EndDate (string dd/mm/yyyy)
//   - optional: Page (int), Limit (int)
// https://developers.cryptomkt.com/#listado-de-ordenes-de-pago
func (client *Client) GetPaymentOrders(arguments ...args.Argument) (*PaymentOrderList, error) {
	required := []string{"start_date", "end_date"}
	req, err := makeReq(required, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetPaymentOrders: %s", err)
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

// GetPaymentStatus gives the status of a payment order.
// Returns a PaymentOrder struct that supports Refresh(), to update its information.
//
// List of accepted Arguments:
//   - required: Id (string)
//   - optional: none
// https://developers.cryptomkt.com/#estado-de-orden-de-pago
func (client *Client) GetPaymentStatus(arguments ...args.Argument) (*PaymentOrder, error) {
	required := []string{"id"}
	resp, err := client.postReq("payment/status", "GetPaymentStatus", required, arguments...)
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

// GetMarkets returns the available markets in cryptomarket as a string slice
//
// https://developers.cryptomkt.com/mercado
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

// GetTicker returns a list of Tickers.
// It returns (nil,error) when an error is raised and ([]Tiker, nil)
// when the operation is successful. The data fields are: High, Low, Ask, Bid,
// LastPrice, Volume, Market and Timestamp.
//
// List of accepted Arguments:
//   - required: none
//   - optional: Market (string)
// https://developers.cryptomkt.com/#ticker
func (client *Client) GetTicker(arguments ...args.Argument) ([]Ticker, error) {
	resp, err := client.getReq("ticker", "GetTicker", []string{}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %s", err)
	}
	var tResp TickerResponse
	json.Unmarshal(resp, &tResp)
	if tResp.Status == "error" {
		return nil, fmt.Errorf("error from the server side: %s", tResp.Message)
	}
	return tResp.Data, nil
}

// GetBook returns a pointer to a Book struct with the data given by
// the api and an error message. It returns (nil, error) when an error
// is raised and (*Book, nil) when the operation is successful.
// The data fields are: Price, Amount and Timestamp. To access these fields,
// you can call them by *Book.Data[indexYouWant].FieldYouWant
//
// List of accepted Arguments:
//   - required: Market (string), Type (string)
//   - optional: Page (int), Limit (int)
// https://developers.cryptomkt.com/#ordenes
func (client *Client) GetBook(arguments ...args.Argument) (*Book, error) {
	required := []string{"market", "type"}
	req, err := makeReq(required, arguments...)
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
// The data fields are MarketTaker, Price, Amount, Tid, Timestamp and Market.
// You can access them by *Trades.Data[indexYouWant].FieldYouWant
//
// List of accepted Arguments:
//   - required: Market (string)
//   - optional: Start (string YYYY-MM-DD), End (YYYY-MM-DD), Page (int), Limit (int)
// https://developers.cryptomkt.com/#trades
func (client *Client) GetTrades(arguments ...args.Argument) (*Trades, error) {
	required := []string{"market"}
	req, err := makeReq(required, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTrades: %s", err)
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
// The data fields are classified in two categories, Ask and Bid.
// The fields are CandleId, OpenPrice, HightPrice, ClosePrice, LowPrice, VolumeSum
// CandleDate and TickCount. To access the data you can call this way:
// *Prices.Data.Ask[indexYouWant].FieldYouWant or *Prices.Data.Bid[indexYouWant].FieldYouWant
//
// List of accepted Arguments:
//   - required: Market (string), Timeframe (string)
//   - optional: Page (int), Limit (int)
// https://developers.cryptomkt.com/#precios
func (client *Client) GetPrices(arguments ...args.Argument) (*Prices, error) {
	required := []string{"market", "timeframe"}
	req, err := makeReq(required, arguments...)
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
