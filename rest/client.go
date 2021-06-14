package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/cryptomarket/cryptomarket-go/args"

	"github.com/cryptomarket/cryptomarket-go/models"
)

// http methods
const (
	methodGet    = "GET"
	methodPut    = "PUT"
	methodPost   = "POST"
	methodDelete = "DELETE"
)

// for authentication porpouses
const (
	publicCall  = true
	privateCall = false
)

const (
	transferTypeExchangeToBank = "exchangeToBank"
	transferTypeBankToExchange = "bankToExchange"
)

// Client handles all the comunication with the rest API
type Client struct {
	hclient httpclient
}

// NewClient creates a new rest client to communicate with the exchange.
// Requests to the exchange via this clients use the args package for aguments.
// All requests accepts contexts for cancelation.
func NewClient(apiKey, apiSecret string) (client *Client) {
	client = &Client{
		hclient: newHTTPClient(apiKey, apiSecret),
	}
	return
}

func (client *Client) publicGet(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodGet, publicCall, endpoint, params, model)
}

func (client *Client) privateGet(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodGet, privateCall, endpoint, params, model)
}

func (client *Client) post(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodPost, privateCall, endpoint, params, model)
}

func (client *Client) put(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodPut, privateCall, endpoint, params, model)
}

func (client *Client) delete(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodDelete, privateCall, endpoint, params, model)
}

func (client *Client) doRequest(ctx context.Context, method string, public bool, endpoint string, params map[string]interface{}, model interface{}) error {
	data, err := client.hclient.doRequest(ctx, method, endpoint, params, public)
	if err != nil {
		return err
	}
	return client.handleResponseData(data, model)
}

func (client *Client) handleResponseData(data []byte, model interface{}) error {
	errorResponse := models.ErrorMetadata{}
	json.Unmarshal(data, &errorResponse)
	serverError := errorResponse.Error
	if serverError != nil { // is a real error
		return fmt.Errorf("CryptomarketAPIError: (code=%v) %v. %v", serverError.Code, serverError.Message, serverError.Description)
	}
	err := json.Unmarshal(data, model)
	if err != nil {
		return errors.New("CryptomarketSDKError: Failed to parse response data: " + err.Error())
	}
	return nil
}

// GetCurrencies gets a list of all currencies or specified currencies.
//
// https://api.exchange.cryptomarket.com/#currencies
//
// Arguments:
//  Currencies([]string) // Optional. A list of currencies ids
func (client *Client) GetCurrencies(ctx context.Context, arguments ...args.Argument) (result []models.Currency, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointCurrency, params, &result)
	return
}

// GetCurrency gets the data of a currency.
//
// https://api.exchange.cryptomarket.com/#currencies
//
// Arguments:
//  Currency(string) // A currency id
func (client *Client) GetCurrency(ctx context.Context, arguments ...args.Argument) (result *models.Currency, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointCurrency+"/"+params["currency"].(string), nil, &result)
	return
}

// GetSymbols gets a list of all symbols or for specified symbols.
//
// A symbol is the combination of the base currency (first one) and quote currency (second one).
//
// https://api.exchange.cryptomarket.com/#symbols
//
// Arguments:
//  Symbols([]string) // Optional. A list of symbol ids
func (client *Client) GetSymbols(ctx context.Context, arguments ...args.Argument) (result []models.Symbol, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointSymbol, params, &result)
	return
}

// GetSymbol gets a symbol by its id.
//
// A symbol is the combination of the base currency (first one) and quote currency (second one).
//
// https://api.exchange.cryptomarket.com/#symbols
//
// Arguments:
//  Symbol(string) // A symbol id
func (client *Client) GetSymbol(ctx context.Context, arguments ...args.Argument) (result *models.Symbol, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointSymbol+"/"+params["symbol"].(string), nil, &result)
	return
}

// GetTickers gets tickers for all symbols or for specified symbols.
//
// https://api.exchange.cryptomarket.com/#tickers
//
// Arguments:
//  Symbols([]string) // Optional. A list of symbol ids
func (client *Client) GetTickers(ctx context.Context, arguments ...args.Argument) (result []models.Ticker, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointTicker, params, &result)
	return
}

// GetTicker gets the ticker of a symbol
//
// https://api.exchange.cryptomarket.com/#tickers
//
// Arguments:
//  Symbol(string) // A symbol id
func (client *Client) GetTicker(ctx context.Context, arguments ...args.Argument) (result *models.Ticker, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointTicker+"/"+params["symbol"].(string), nil, &result)
	return
}

// GetTrades gets trades for all symbols or for specified symbols.
//
// 'from' param and 'till' param must have the same format, both index of both timestamp.
//
// https://api.exchange.cryptomarket.com/#trades
//
// Arguments:
//  Symbols([]string) // Optional. A list of symbol ids
//  Sort(SortType)    // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  From(string)      // Optional. Initial value of the queried interval
//  Till(string)      // Optional. Last value of the queried interval
//  Limit(int)        // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)       // Optional. Default is 0. Max is 100000
func (client *Client) GetTrades(ctx context.Context, arguments ...args.Argument) (result map[string][]models.PublicTrade, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointTrade, params, &result)
	return
}

// GetTradesOfSymbol gets trades of a symbol.
//
// 'from' param and 'till' param must have the same format, both index of both timestamp.
//
// https://api.exchange.cryptomarket.com/#trades
//
// Arguments:
//  Symbol(string)     // A symbol id
//  Sort(SortType)     // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  SortBy(SortByType) // Optional. Defines the sorting type. SortByTimestamp or SortByID
//  From(string)       // Optional. Initial value of the queried interval
//  Till(string)       // Optional. Last value of the queried interval
//  Limit(int)         // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)        // Optional. Default is 0. Max is 100000
func (client *Client) GetTradesOfSymbol(ctx context.Context, arguments ...args.Argument) (result []models.PublicTrade, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointTrade+"/"+params["symbol"].(string), params, &result)
	return
}

// GetOrderbooks gets orderbooks for all symbols or for the specified symbols.
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level.
//
// https://api.exchange.cryptomarket.com/#order-book
//
// Arguments:
//  Symbols([]string) // Optional. A list of symbol ids
//  Limit(int)        // Optional. Limit of order book levels. Set to 0 to view full list of order book levels
func (client *Client) GetOrderbooks(ctx context.Context, arguments ...args.Argument) (result map[string]models.OrderBook, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointOrderbook, params, &result)
	return
}

// GetOrderbook gets an orderbook of a symbol.
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level.
//
// https://api.exchange.cryptomarket.com/#order-book
//
// Arguments:
//  Symbol(string) // A symbol id
//  Limit(int)     // Optional. Limit of order book levels. Set to 0 to view full list of order book levels
func (client *Client) GetOrderbook(ctx context.Context, arguments ...args.Argument) (result *models.OrderBook, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointOrderbook+"/"+params["symbol"].(string), params, &result)
	return
}

// MarketDepthSearch gets an orderbook of a symbol with market depth info.
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level.
//
// https://api.exchange.cryptomarket.com/#order-book
//
// Arguments:
//  Symbol(string) // The symbol id
//  Volume(string) // Desired volume for market depth search
func (client *Client) MarketDepthSearch(ctx context.Context, arguments ...args.Argument) (result *models.OrderBook, err error) {
	params, err := args.BuildParams(arguments, "symbol", "volume")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointOrderbook+"/"+params["symbol"].(string), params, &result)
	return
}

// GetCandles get candles for all symbols or for specified symbols.
//
// Candels are used for OHLC representation.
//
// https://api.exchange.cryptomarket.com/#candles
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
//  Period(PeriodType) // Optional. A valid tick interval. Default is PeriodType30Minutes
//  Sort(SortType)     // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  From(string)       // Optional. Initial value of the queried interval
//  Till(string)       // Optional. Last value of the queried interval
//  Limit(int)         // Optional. Candles per query. Defaul is 100. Max is 1000
//  Offset(int)        // Optional. Default is 0. Max is 100000
func (client *Client) GetCandles(ctx context.Context, arguments ...args.Argument) (result map[string][]models.Candle, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointCandle, params, &result)
	return
}

// GetCandlesOfSymbol get candles for a symbol.
//
// Candels are used for OHLC representation.
//
// https://api.exchange.cryptomarket.com/#candles
//
// Arguments:
//  Symbol(string)     // A symbol id
//  Period(PeriodType) // Optional. A valid tick interval. Default is PeriodType30Minutes
//  Sort(SortType)     // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  From(string)       // Optional. Initial value of the queried interval
//  Till(string)       // Optional. Last value of the queried interval
//  Limit(int)         // Optional. Candles per query. Defaul is 100. Max is 1000
//  Offset(int)        // Optional. Default is 0. Max is 100000
func (client *Client) GetCandlesOfSymbol(ctx context.Context, arguments ...args.Argument) (result []models.Candle, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointCandle+"/"+params["symbol"].(string), params, &result)
	return
}

/////////////
// TRADING //
/////////////

// GetTradingBalance gets the account trading balance.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#trading-balance
func (client *Client) GetTradingBalance(ctx context.Context) (result []models.Balance, err error) {
	err = client.privateGet(ctx, endpointTradingBalance, nil, &result)
	return
}

// GetActiveOrders gets the account active orders.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#get-active-orders
//
// Arguments:
//  Symbol(string) // Optional. A symbol for filtering active orders
func (client *Client) GetActiveOrders(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointOrder, params, &result)
	return
}

// GetActiveOrder gets an active order by its client order id.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#get-active-orders
//
// Arguments:
//  ClientOrderId(string // The clientOrderId of the order
//  Wait(int)            // Optional. Time in milliseconds Max value is 60000. Default value is None. While using long polling request: if order is filled, cancelled or expired order info will be returned instantly. For other order statuses, actual order info will be returned after specified wait time.
func (client *Client) GetActiveOrder(ctx context.Context, arguments ...args.Argument) (result *models.Order, err error) {
	params, _ := args.BuildParams(arguments, "clientOrderId")
	err = client.privateGet(ctx, endpointOrder+"/"+params["clientOrderId"].(string), nil, &result)
	return
}

// CreateOrder Creates a new order.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#create-new-order
//
// Arguments:
//  Symbol(string)               // Trading symbol
//  Side(SideType)               // SideTypeBuy or SideTypeSell
//  Quantity(string)             // Order quantity
//  ClientOrderID(string)        // Optional. If given must be unique within the trading day, including all active orders. If not given, is generated by the server
//  Type(OrderType)              // Optional. Default is OrderTypeLimit
//  TimeInForce(TimeInForceType) // Optional. Default is TimeInForceTypeGTC
//  Price(string)                // Required for OrderTypelimit and OrderTypeStopLimit. limit price of the order
//  StopPrice(string)            // Required for OrderTypeStopLimit and OrderTypeStopMarket orders. stop price of the order
//  ExpireTime(string)           // Required for orders with TimeInForceTypeGDT
//  StrictValidate(bool)         // Optional. If False, the server rounds half down for tickerSize and quantityIncrement. Example of ETHBTC: tickSize = '0.000001', then price '0.046016' is valid, '0.0460165' is invalid
//  PostOnly(bool)               // Optional. If True, your post_only order causes a match with a pre-existing order as a taker, then the order will be cancelled
func (client *Client) CreateOrder(ctx context.Context, arguments ...args.Argument) (result *models.Order, err error) {
	params, err := args.BuildParams(arguments, "symbol", "side", "quantity")
	if err != nil {
		return
	}
	if clientOrderID, ok := params["clientOrderId"]; ok {
		err = client.put(ctx, endpointOrder+"/"+clientOrderID.(string), params, &result)
	} else {
		err = client.post(ctx, endpointOrder, params, &result)
	}
	return
}

// CancelAllOrders cancel all active orders, or all active orders for a specified symbol.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#cancel-orders
func (client *Client) CancelAllOrders(ctx context.Context) (result []models.Order, err error) {
	err = client.delete(ctx, endpointOrder, nil, &result)
	return
}

// CancelOrder cancel the order with clientOrderId.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#cancel-order-by-clientorderid
//
// Arguments:
//  clientOrderId(string) // the client id of the order to cancel
func (client *Client) CancelOrder(ctx context.Context, arguments ...args.Argument) (result *models.Order, err error) {
	params, err := args.BuildParams(arguments, "clientOrderId")
	if err != nil {
		return
	}
	err = client.delete(ctx, endpointOrder+"/"+params["clientOrderId"].(string), nil, &result)
	return
}

// GetTradingFee gets personal trading commission rates for a symbol.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#get-trading-commission
//
// Arguments:
//  Symbol(string) The symbol of the comission rates
func (client *Client) GetTradingFee(ctx context.Context, arguments ...args.Argument) (result *models.TradingFee, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTradingFee+"/"+params["symbol"].(string), nil, &result)
	return
}

/////////////////////
// Trading history //
/////////////////////

// GetOrderHistory gets the account order history.
//
// All not active orders older than 24 are deleted.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#orders-history
//
// Arguments:
//  Symbol(string) // Optional. Filter orders by symbol
//  From(string)   // Optional. Initial value of the queried interval
//  Till(string)   // Optional. Last value of the queried interval
//  Limit(int)     // Optional. Candles per query. Defaul is 100. Max is 1000
//  Offset(int)    // Optional. Default is 0. Max is 100000
func (client *Client) GetOrderHistory(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointOrderHistory, params, &result)
	return
}

// GetOrders gets orders with the clientOrderId.
//
// All not active orders older than 24 are deleted.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#orders-history
//
// Arguments:
//  ClientOrderID(string) // the clientOrderId of the orders
func (client *Client) GetOrders(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, err := args.BuildParams(arguments, "clientOrderId")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointOrderHistory, params, &result)
	return
}

// GetTradeHistory gets the user's trading history.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#orders-history
//
// Arguments:
//  Symbol(string)     // Optional. Filter trades by symbol
//  Sort(SortType)     // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  SortBy(SortByType) // Optional. Defines the sorting type. SortByTimestamp or SortByID
//  From(string)       // Optional. Initial value of the queried interval. Id or datetime
//  Till(string)       // Optional. Last value of the queried interval. Id or datetime
//  Limit(int)         // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)        // Optional. Default is 0. Max is 100000
//  Margin(string)     // Optional. Default is MarginTypeInclude
func (client *Client) GetTradeHistory(ctx context.Context, arguments ...args.Argument) (result []models.Trade, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointTradeHistory, params, &result)
	return
}

// GetTradesByOrderID gets the account's trading orders of a specified order id
//
// Requires authentication
//
// https://api.exchange.cryptomarket.com/#trades-by-order
//
// Arguments:
//  OrderId(int64) // Order unique identifier assigned by exchange
func (client *Client) GetTradesByOrderID(ctx context.Context, arguments ...args.Argument) (result []models.Trade, err error) {
	params, err := args.BuildParams(arguments, "orderId")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointOrderHistory+"/"+strconv.FormatInt(params["orderId"].(int64), 10)+"/trades", nil, &result)
	return
}

////////////////////////
// ACCOUNT MANAGAMENT //
////////////////////////

// GetAccountBalance gets the user account balance.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#account-balance
func (client *Client) GetAccountBalance(ctx context.Context) (result []models.Balance, err error) {
	err = client.privateGet(ctx, endpointAccountBalance, nil, &result)
	return
}

// GetDepositCryptoAddress gets the current address of a currency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#deposit-crypto-address
//
// Arguments:
//  Currency(string) // currency to get the address
func (client *Client) GetDepositCryptoAddress(ctx context.Context, arguments ...args.Argument) (result *models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointCryptoAdress+"/"+params["currency"].(string), nil, &result)
	return
}

// CreateDepositCryptoAddress Creates a new address for the currency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#deposit-crypto-address
//
// Arguments:
//  Currency(string) // currency to create a new address
func (client *Client) CreateDepositCryptoAddress(ctx context.Context, arguments ...args.Argument) (result *models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointCryptoAdress+"/"+params["currency"].(string), nil, &result)
	return
}

// GetLast10DepositCryptoAddresses gets the last 10 addresses used for deposit by currency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#last-10-deposit-crypto-address
//
// Arguments:
//  Currency(string) // currency to get the list of addresses
func (client *Client) GetLast10DepositCryptoAddresses(ctx context.Context, arguments ...args.Argument) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointCryptoAdresses+"/"+params["currency"].(string), nil, &result)
	return
}

// GetLast10UsedCryptoAddresses gets the last 10 unique addresses used for withdraw by currency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#last-10-used-crypto-address
//
// Arguments:
//  Currency(string) // currency to get the list of addresses
func (client *Client) GetLast10UsedCryptoAddresses(ctx context.Context, arguments ...args.Argument) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointUsedAddressed+"/"+params["currency"].(string), nil, &result)
	return
}

// WithdrawCrypto withdraws cryptocurrency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#withdraw-crypto
//
// Arguments:
//  Currency(string)  // currency code of the crypto to withdraw
//  Amount(string)    // the amount to be sent to the specified address
//  Address(string)   // the address identifier
//  PaymentID(string) // Optional.
//  IncludeFee(bool)  // Optional. If true then the total spent amount includes fees. Default false
//  AutoCommit(bool)  // Optional. If false then you should commit or rollback transaction in an hour. Used in two phase commit schema. Default true
func (client *Client) withdrawCrypto(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "currency", "amount", "address")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointWithdrawCrypto, params, &result)
	return
}

// TransferConvert converts between currencies.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#transfer-convert-between-currencies
//
// Arguments:
//  FromCurrency(string) // currency code of origin
//  ToCurrency(string)   // currency code of destiny
//  Amount(string)       // the amount to be sent
func (client *Client) TransferConvert(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "fromCurrency", "toCurrency", "amount")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointTransferConvert, params, &result)
	return
}

// CommitWithdrawCrypto commit a withdrawal of cryptocurrency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#withdraw-crypto-commit-or-rollback
//
// Arguments:
//  ID(string) // the withdrawal transaction identifier
func (client *Client) CommitWithdrawCrypto(ctx context.Context, arguments ...args.Argument) (result bool, err error) {
	params, err := args.BuildParams(arguments, "id")
	if err != nil {
		return
	}
	data := make(map[string]bool)
	err = client.put(ctx, endpointWithdrawCrypto+"/"+params["id"].(string), nil, &data)
	if err != nil {
		return
	}
	if res, ok := data["result"]; ok {
		result = res
		return
	}
	err = errors.New("CryptomarketSDKError: invalid response format")
	return
}

// RollbackWithdrawCrypto rollback a withdrawal of cryptocurrency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#withdraw-crypto-commit-or-rollback
//
// Arguments:
//  ID(string) // the withdrawal transaction identifier
func (client *Client) RollbackWithdrawCrypto(ctx context.Context, arguments ...args.Argument) (result bool, err error) {
	params, err := args.BuildParams(arguments, "id")
	if err != nil {
		return
	}
	data := make(map[string]bool)
	err = client.delete(ctx, endpointWithdrawCrypto+"/"+params["id"].(string), nil, &data)
	if err != nil {
		return
	}
	if res, ok := data["result"]; ok {
		result = res
		return
	}
	err = errors.New("CryptomarketSDKError: invalid response format")
	return
}

// GetEstimatesWithdrawFee gets an estimate of the withdrawal fee.
//
// Requires authetication.
//
// https://api.exchange.cryptomarket.com/#estimate-withdraw-fee
//
// Arguments:
//  Currency(string) // the currency code for withdraw
//  Amount(string)   // the expected withdraw amount
func (client *Client) GetEstimatesWithdrawFee(ctx context.Context, arguments ...args.Argument) (result string, err error) {
	params, err := args.BuildParams(arguments, "currency", "amount")
	if err != nil {
		return
	}
	data := make(map[string]string)
	err = client.privateGet(ctx, endpointEstimateWithdraw, params, &data)
	if err != nil {
		return
	}
	if res, ok := data["fee"]; ok {
		result = res
		return
	}
	err = errors.New("CryptomarketSDKError: invalid response format")
	return
}

// CheckIfCryptoAddressIsMine check if an address is from this account.
//
// Requires authentication
//
// https://api.exchange.cryptomarket.com/#check-if-crypto-address-belongs-to-current-account
//
// Arguments:
//  Address(string) // The address to check
func (client *Client) CheckIfCryptoAddressIsMine(ctx context.Context, arguments ...args.Argument) (result bool, err error) {
	params, err := args.BuildParams(arguments, "address")
	if err != nil {
		return
	}
	data := make(map[string]bool)
	err = client.delete(ctx, endpointCryptoAddressIsMine+"/"+params["address"].(string), nil, &data)
	if err != nil {
		return
	}
	if res, ok := data["result"]; ok {
		result = res
		return
	}
	err = errors.New("CryptomarketSDKError: invalid response format")
	return
}

// TransferMoneyFromTradingToAccountBalance transfer money from the trading balance to the account balance.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#transfer-money-between-trading-account-and-bank-account
//
// Arguments:
//  Currency(string) // Currency code for transfering
//  Amount(string)   // Amount to be transfered
func (client *Client) TransferMoneyFromTradingToAccountBalance(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	arguments = append(arguments, args.TransferType(transferTypeExchangeToBank))
	params, err := args.BuildParams(arguments, "currency", "amount")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointAccountTranser, params, &result)
	return
}

// TransferMoneyFromAccountToTradingBalance transfer money from the account balance to the trading balance.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#transfer-money-between-trading-account-and-bank-account
//
// Arguments:
//  Currency(string) // Currency code for transfering
//  Amount(string)   // Amount to be transfered
func (client *Client) TransferMoneyFromAccountToTradingBalance(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	arguments = append(arguments, args.TransferType(transferTypeBankToExchange))
	params, err := args.BuildParams(arguments, "currency", "amount")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointAccountTranser, params, &result)
	return
}

// TransferMoneyToAnotherUser transfers money to another user.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#transfer-money-to-another-user-by-email-or-username
//
// Arguments:
//  Currency(string)   // currency code
//  Amount(string)     // amount to be transfered between balances
//  TransferBy(string) // TransferByEmail or TransferByUsername
//  Identifier(string) // the email or the username
func (client *Client) TransferMoneyToAnotherUser(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "currency", "amount", "by", "identifier")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointAccountTranserInternal, params, &result)
	return
}

// GetTransactionHistory gets the transactions of the account by currency.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#get-transactions-history
//
// Arguments:
//  Currency(string)   // Currency code to get the transaction history
//  Sort(SortType)     // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  SortBy(SortByType) // Optional. Defines the sorting type. SortByTimestamp or SortByID
//  From(string)       // Optional. Initial value of the queried interval. Id or datetime
//  Till(string)       // Optional. Last value of the queried interval. Id or datetime
//  Limit(int)         // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)        // Optional. Default is 0. Max is 100000
func (client *Client) GetTransactionHistory(ctx context.Context, arguments ...args.Argument) (result []models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTransactionHistory, params, &result)
	return
}

// GetTransaction gets the transactions of the account by its identifier.
//
// Requires authentication.
//
// https://api.exchange.cryptomarket.com/#get-transactions-history
//
// Arguments:
//  ID(string) // The identifier of the transaction
func (client *Client) GetTransaction(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "id")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTransactionHistory+"/"+params["id"].(string), nil, &result)
	return
}
