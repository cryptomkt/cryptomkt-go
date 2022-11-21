package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/internal"

	"github.com/cryptomarket/cryptomarket-go/models"
)

// http methods
const (
	methodGet    = "GET"
	methodPut    = "PUT"
	methodPost   = "POST"
	methodPatch  = "PATCH"
	methodDelete = "DELETE"
)

// for authentication porpouses
const (
	publicCall  = true
	privateCall = false
)

// Client handles all the comunication with the rest API
type Client struct {
	hclient httpclient
}

// NewClient creates a new rest client to communicate with the exchange.
// Requests to the exchange via this client use the args package for aguments.
// All requests accept contexts for cancellation.
// arguments:
//  apiKey // the key of the user
//  apiSecret // the secret key of the user
//  window // the window of execution for requests to the server in milliseconds. Max is 60_000 (miliseconds). Use 0 for default window (10 seconds)
func NewClient(apiKey, apiSecret string, window int) *Client {
	return &Client{
		hclient: newHTTPClient(apiKey, apiSecret, window),
	}
}

// NewPublicClient creates a new rest client with no credentials to communicate with the exchange.
// Only works with public calls.
// Requests to the exchange via this client use the args package for arguments.
// All requests accept contexts for cancellation
func NewPublicClient() *Client {
	return &Client{
		hclient: newHTTPClient("", "", 0),
	}
}

func (client *Client) publicGet(
	ctx context.Context,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	return client.doRequest(
		ctx,
		methodGet,
		publicCall,
		endpoint,
		params,
		model,
	)
}

func (client *Client) privateGet(
	ctx context.Context,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	return client.doRequest(
		ctx,
		methodGet,
		privateCall,
		endpoint,
		params,
		model,
	)
}

func (client *Client) post(
	ctx context.Context,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	return client.doRequest(
		ctx,
		methodPost,
		privateCall,
		endpoint,
		params,
		model,
	)
}

func (client *Client) put(
	ctx context.Context,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	return client.doRequest(
		ctx,
		methodPut,
		privateCall,
		endpoint,
		params,
		model,
	)
}

func (client *Client) patch(
	ctx context.Context,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	return client.doRequest(
		ctx,
		methodPatch,
		privateCall,
		endpoint,
		params,
		model,
	)
}

func (client *Client) delete(
	ctx context.Context,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	return client.doRequest(
		ctx,
		methodDelete,
		privateCall,
		endpoint,
		params,
		model,
	)
}

func (client *Client) doRequest(
	ctx context.Context,
	method string,
	public bool,
	endpoint string,
	params map[string]interface{},
	model interface{},
) error {
	data, err := client.hclient.doRequest(
		ctx,
		method,
		endpoint,
		params,
		public,
	)
	if err != nil {
		return err
	}
	return client.handleResponseData(data, model)
}

func (client *Client) handleResponseData(
	data []byte,
	model interface{},
) error {
	errorResponse := models.ErrorMetadata{}
	json.Unmarshal(data, &errorResponse)
	apiError := errorResponse.APIError
	if apiError != nil { // is a real error
		return fmt.Errorf(
			"CryptomarketAPIError: (code=%v) %v. %v",
			apiError.Code,
			apiError.Message,
			apiError.Description,
		)
	}
	err := json.Unmarshal(data, model)
	if err != nil {
		return errors.New(
			"CryptomarketSDKError: Failed to parse response data: " + err.Error(),
		)
	}
	return nil
}

// GetCurrencies gets a map of all currencies or specified currencies. indexed by id
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#currencies
//
// Arguments:
//  Currencies([]CurrenciesType)  // Optional. A list of currencies ids
func (client *Client) GetCurrencies(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.Currency, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointCurrency, params, &result)
	return
}

// GetCurrency gets the data of a currency
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#currencies
//
// Arguments:
//  Currency(string)  // A currency id
func (client *Client) GetCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Currency, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	err = client.publicGet(
		ctx,
		endpointCurrency+"/"+params["currency"].(string),
		nil,
		&result,
	)
	return
}

// GetSymbols gets a map of all symbols or for specified symbols
//
// A symbol is the combination of the base currency (first one) and quote currency (second one)
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#symbols
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *Client) GetSymbols(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.Symbol, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointSymbol, params, &result)
	return
}

// GetSymbol gets a symbol by its id
//
// A symbol is the combination of the base currency (first one) and quote currency (second one)
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#symbols
//
// Arguments:
//  Symbol(string)  // A symbol id
func (client *Client) GetSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Symbol, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	err = client.publicGet(
		ctx,
		endpointSymbol+"/"+params["symbol"].(string),
		nil,
		&result,
	)
	return
}

// GetTickers gets a map of tickers for all symbols or for specified symbols. indexed by symbol id
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#tickers
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *Client) GetTickers(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.Ticker, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointTicker, params, &result)
	return
}

// GetTickerOfSymbol gets the ticker of a symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#tickers
//
// Arguments:
//  Symbol(string)  // A symbol id
func (client *Client) GetTickerOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Ticker, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	err = client.publicGet(
		ctx,
		endpointTicker+"/"+params["symbol"].(string),
		nil,
		&result,
	)
	return
}

// GetPrices gets a map of quotation prices of currencies
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#prices
//
// Arguments:
//  To(string)  // Target currency code
//  From(string)  // Optional. Source currency rate
func (client *Client) GetPrices(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.Price, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameTo)
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointPrices, params, &result)
	return
}

// GetPricesHistory get the quotation prices history
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#prices
//
// Arguments:
//  To(string)  // Target currency code
//  From(string)  // Optional. Source currency rate
//  Period(PeriodType)  // Optional. A valid tick interval. Period1Minute, Period3Minutes, Period5Minutes, Period15Minutes, Period30Minutes, Period1Hour, Period4Hours, Period1Day, Period7Days, Period1Month. Default is Period30Minutes
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  Since(string)  // Optional. Initial value of the queried interval. As Datetime
//  Until(string)  // Optional. Last value of the queried interval. As Datetime
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 1. Min is 1. Max is 1000
func (client *Client) GetPricesHistory(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.PriceHistory, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameTo)
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointPriceHistory, params, &result)
	return
}

// GetTickerLastPrices gets a map of the ticker's last prices for all symbols or for the specified symbols
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#prices
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *Client) GetTickerLastPrices(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.Price, err error) {
	params, err := args.BuildParams(arguments)
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointPriceTicker, params, &result)
	return
}

// GetTickerLastPriceOfSymbol gets the ticker's last price of a symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#prices
//
// Arguments:
//  Symbol(string)  // A symbol id
func (client *Client) GetTickerLastPricesOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Price, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	err = client.publicGet(
		ctx,
		endpointPriceTicker+"/"+params["symbol"].(string),
		params,
		&result,
	)
	return
}

// GetTrades gets a map of trades for all symbols or for specified symbols. indexed by symbol
//
// From param and Till param must have the same format, both id or both timestamp
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#trades
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
//  SortBy(SortByType)  // Optional. Sorting parameter. SortByID or SortByTimestamp. Default is SortByTimestamp
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  From(string)  // Optional. Initial value of the queried interval
//  Till(string)  // Optional. Last value of the queried interval
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 10. Min is 1. Max is 1000
func (client *Client) GetTrades(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string][]models.PublicTrade, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointTrade, params, &result)
	return
}

// GetTradesOfSymbol gets trades of a symbol
//
// From param and Till param must have the same format, both index of both timestamp
//
// https://api.exchange.cryptomarket.com/#trades
//
// Arguments:
//  Symbol(string)  // A symbol id
//  SortBy(SortByType)  // Optional. Sorting parameter. SortByID or SortByTimestamp. Default is SortByTimestamp
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  Since(string)  // Optional. Initial value of the queried interval
//  Until(string)  // Optional. Last value of the queried interval
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 10. Min is 1. Max is 1000
//  Offset(int)  // Optional. Default is 0. Min is 0. Max is 100000
func (client *Client) GetTradesOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.PublicTrade, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	err = client.publicGet(
		ctx,
		endpointTrade+"/"+params["symbol"].(string),
		params,
		&result,
	)
	return
}

// GetOrderBooks gets a map of orderbooks for all symbols or for the specified symbols
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#order-books
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
//  Depth(int)  // Optional. Order Book depth. Default value is 100. Set to 0 to view the full Order Book
func (client *Client) GetOrderbooks(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string]models.OrderBook, err error) {
	params, _ := args.BuildParams(arguments)
	response := make(map[string]models.OrderBookJson)
	err = client.publicGet(ctx, endpointOrderbook, params, &response)
	if err != nil {
		return
	}
	result = make(map[string]models.OrderBook)
	for symbol, obJson := range response {
		result[symbol] = models.FromOrderbookJsonToOrderbook(obJson)
	}
	return
}

// GetOrderBookOfSymbol get order book of a symbol
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#order-books
//
// Arguments:
//  Symbol(string)  // A symbol id
//  Depth(int)  // Optional. Order Book depth. Default value is 100. Set to 0 to view the full Order Book
func (client *Client) GetOrderBookOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.OrderBook, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	response := models.OrderBookJson{
		Ask:       make([][]string, 0),
		Bid:       make([][]string, 0),
		Timestamp: "",
	}
	err = client.publicGet(
		ctx,
		endpointOrderbook+"/"+params["symbol"].(string),
		params,
		&response,
	)
	ob := models.FromOrderbookJsonToOrderbook(response)
	result = &ob
	return
}

// GetOrderBookVolumeOfSymbol get order book of a symbol with the desired volume for market depth search
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#order-books
//
// Arguments:
//  Symbol(string)  // A symbol id
//  Volume(string)  // Desired volume for market depth search
func (client *Client) GetOrderBookVolumeOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.OrderBook, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameSymbol,
		internal.ArgNameVolume,
	)
	if err != nil {
		return
	}
	response := models.OrderBookJson{
		Ask:       make([][]string, 0),
		Bid:       make([][]string, 0),
		Timestamp: "",
	}
	err = client.publicGet(
		ctx,
		endpointOrderbook+"/"+params["symbol"].(string),
		params,
		&response,
	)
	ob := models.FromOrderbookJsonToOrderbook(response)
	result = &ob
	return
}

// GetCandles gets a map of candles for all symbols or for specified symbols
//
// Candels are used for OHLC representation
//
// The result contains candles with non-zero volume only (no trades = no candles)
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#candles
//
// Arguments:
//  Symbols([]string)  // A list of symbol ids
//  Period(PeriodType)  // Optional. A valid tick interval. Period1Minute, Period3Minutes, Period5Minutes, Period15Minutes, Period30Minutes, Period1Hour, Period4Hours, Period1Day, Period7Days, Period1Month. Default is Period30Minutes
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  From(string)  // Optional. Initial value of the queried interval. As DateTime
//  Till(string)  // Optional. Last value of the queried interval. As DateTime
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 10. Min is 1. Max is 1000
func (client *Client) GetCandles(
	ctx context.Context,
	arguments ...args.Argument,
) (result map[string][]models.Candle, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointCandle, params, &result)
	return
}

// GetCandlesOfSymbol get candles of a symbol
//
// Candels are used for OHLC representation
//
// The result contains candles with non-zero volume only (no trades = no candles)
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#candles
//
// Arguments:
//  Symbol(string)  // A symbol id
//  Period(PeriodType)  // Optional. A valid tick interval. Period1Minute, Period3Minutes, Period5Minutes, Period15Minutes, Period30Minutes, Period1Hour, Period4Hours, Period1Day, Period7Days, Period1Month. Default is Period30Minutes
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  From(string)  // Optional. Initial value of the queried interval. As DateTime
//  Till(string)  // Optional. Last value of the queried interval. As DateTime
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 100. Min is 1. Max is 1000
//  Offset(int)  // Optional. Default is 0. Min is 0. Max is 100000
func (client *Client) GetCandlesOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Candle, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	err = client.publicGet(
		ctx,
		endpointCandle+"/"+params["symbol"].(string),
		nil,
		&result,
	)
	return
}

//////////////////
// SPOT TRADING //
//////////////////

// GetSpotTradingBalances gets the user's spot trading balance for all currencies with balance
//
// Requires the "Orderbook, History, Trading balance" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-spot-trading-balance
func (client *Client) GetSpotTradingBalances(
	ctx context.Context,
) (result []models.Balance, err error) {
	err = client.privateGet(ctx, endpointTradingBalance, nil, &result)
	return
}

// GetSpotTradingBalanceOfCurrency gets the user spot trading balance of a currency
//
// Requires the "Orderbook, History, Trading balance" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-spot-trading-balance
//
// Arguments:
//  Currency(string)  // The currency code to query the balance
func (client *Client) GetSpotTradingBalanceOfCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Balance, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	err = client.privateGet(
		ctx,
		endpointTradingBalance+"/"+params[internal.ArgNameCurrency].(string),
		nil,
		&result,
	)
	if err != nil {
		return
	}
	result.Currency = params[internal.ArgNameCurrency].(string)
	return
}

// GetAllActiveSpotOrders gets the user's active spot orders
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-all-active-spot-orders
//
// Arguments:
//  Symbol(string)  // Optional. A symbol for filtering the active spot orders
func (client *Client) GetAllActiveSpotOrders(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Order, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointOrder, params, &result)
	return
}

// GetActiveSpotOrder gets an active spot order by its client order id
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-active-spot-orders
//
// Arguments:
//  ClientOrderID(string)  // The client order id of the order
func (client *Client) GetActiveSpotOrder(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Order, err error) {
	params, _ := args.BuildParams(arguments, internal.ArgNameClientOrderID)
	err = client.privateGet(
		ctx,
		endpointOrder+"/"+params["client_order_id"].(string),
		nil,
		&result,
	)
	return
}

// CreateSpotOrder creates a new spot order
//
// For fee, for price accuracy and quantity, and for order status information see the api docs
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#create-new-spot-order
//
// Arguments:
//  Symbol(string)  // Trading symbol
//  Side(SideType)  // Either SideBuy or SideSell
//  Quantity(string)  // Order quantity
//  ClientOrderID(string)  // Optional. If given must be unique within the trading day, including all active orders. If not given, is generated by the server
//  Type(OrderType)  // Optional. OrderLimit, OrderMarket, OrderStopLimit, OrderStopMarket, OrderTakeProfitLimit or OrderTakeProfitMarket. Default is OrderLimit
//  TimeInForce(TimeInForceType)  // Optional. TimeInForceGTC, TimeInForceIOC, TimeInForceFOK, TimeInForceDay, TimeInForceGTD. Default to TimeInForceGTC
//  Price(string)  // Optional. Required for OrderLimit and OrderStopLimit. limit price of the order
//  StopPrice(string)  // Optional. Required for OrderStopLimit and OrderStopMarket orders. stop price of the order
//  ExpireTime(string)  // Optional. Required for orders with timeInForceGDT
//  StrictValidate(bool)  // Optional. If False, the server rounds half down for tickerSize and quantityIncrement. Example of ETHBTC: tickSize = '0.000001', then price '0.046016' is valid, '0.0460165' is invalid
//  PostOnly(bool)  // Optional. If True, your postOnly order causes a match with a pre-existing order as a taker, then the order will be cancelled
//  TakeRate(string)  // Optional. Liquidity taker fee, a fraction of order volume, such as 0.001 (for 0.1% fee). Can only increase the fee. Used for fee markup.
//  MakeRate(string)  // Optional. Liquidity provider fee, a fraction of order volume, such as 0.001 (for 0.1% fee). Can only increase the fee. Used for fee markup.
func (client *Client) CreateSpotOrder(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Order, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameSymbol,
		internal.ArgNameSide,
		internal.ArgNameQuantity,
	)
	if err != nil {
		return
	}
	err = client.post(ctx, endpointOrder, params, &result)
	return
}

// ReplaceSpotOrder replaces a spot order
//
// For fee, for price accuracy and quantity, and for order status information see the api docs
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#replace-spot-order
//
// Arguments:
//  ClientOrderID(string)  // client order id of the old order
//  NewClientOrderID(string)  // client order id for the new order
//  Quantity(string)  // Order quantity
//  StrictValidate(bool)  // Price and quantity will be checked for incrementation within the symbol’s tick size and quantity step. See the symbol's TickSize and QuantityIncrement
//  Price(string)  // Required for OrderLimit, OrderStopLimit, or OrderTakeProfitLimit. Order price
func (client *Client) ReplaceSpotOrder(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Order, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameClientOrderID,
		internal.ArgNameNewClientOrderID,
		internal.ArgNameQuantity,
	)
	if err != nil {
		return
	}
	ClientOrderID := params["client_order_id"].(string)
	delete(params, "client_order_id")
	err = client.patch(ctx, endpointOrder+"/"+ClientOrderID, params, &result)
	return
}

// CreateSpotOrderList creates a list of spot orders and returns a list of the created spot orders or a possible error
//
// Types or contingency:
//
//  - ContingencyTypeAllOrNone (ContingencyTypeAON) (AON)
//  - ContingencyTypeOneCancelOther (ContingencyTypeOCO) (OCO)
//  - ContingencyOneTriggerOneCancelOther (ContingencyTypeOTOCO) (OTOCO)
//
// Restriction in the number of orders:
//
//  - An AON list must have 2 or 3 orders
//  - An OCO list must have 2 or 3 orders
//  - An OTOCO must have 3 or 4 orders
//
// Symbol restrictions:
//
//  - For an AON order list, the symbol code of orders must be unique for each order in the list.
//  - For an OCO order list, there are no symbol code restrictions.
//  - For an OTOCO order list, the symbol code of orders must be the same for all orders in the list (placing orders in different order books is not supported).
//
// ORDER_TYPE restrictions:
//  - For an AON order list, orders must be OrderLimit or OrderMarket
//  - For an OCO order list, orders must be OrderLimit, OrderStopLimit, OrderStopMarket, OrderTakeProfitLimit or OrderTakeProfitMarket.
//  - An OCO order list cannot include more than one limit order (the same applies to secondary orders in an OTOCO order list).
//  - For an OTOCO order list, the first order must be OrderLimit, OrderMarket, OrderStopLimit, OrderStopMarket, OrderTakeProfitLimit or OrderTakeProfitMarket.
//  - For an OTOCO order list, the secondary orders have the same restrictions as an OCO order
//  - Default is OrderTypeLimit
//
// https://api.exchange.cryptomkt.com/#create-new-spot-order-list-2
//
// Arguments:
//  ContingencyType(ContingencyTypeType) order list type.
//  Orders(OrderRequest[]) the list of OrderRequests in the order list
//  OrderListID(string) order list identifier. If not provided, it will be generated by the system. Must be equal to the client order id of the first order in the requests list.
func (client *Client) CreateSpotOrderList(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Order, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameContingencyType,
		internal.ArgNameOrders,
	)
	if err != nil {
		return
	}
	err = client.post(ctx, endpointOrderList, params, &result)
	return
}

// CancelAllSpotOrders Cancel all active spot orders, or all active orders for a specified symbol
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#cancel-all-spot-orders
func (client *Client) CancelAllSpotOrders(
	ctx context.Context,
) (result []models.Order, err error) {
	err = client.delete(ctx, endpointOrder, nil, &result)
	return
}

// CancelSpotOrder Cancel the order with the client order id
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#cancel-spot-order
//
// Arguments:
//  ClientOrderID(string)  // client order id of the order to cancel
func (client *Client) CancelSpotOrder(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Order, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameClientOrderID,
	)
	if err != nil {
		return
	}
	err = client.delete(
		ctx,
		endpointOrder+"/"+params["client_order_id"].(string),
		nil,
		&result,
	)
	return
}

// GetAllTradingCommissions gets the personal trading commission rates for all symbols
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-all-trading-commission
func (client *Client) GetAllTradingCommissions(
	ctx context.Context,
) (result []models.TradingCommission, err error) {
	err = client.privateGet(ctx,
		endpointTradingCommission,
		nil,
		&result,
	)
	return
}

// GetTradingCommissionOfSymbol gets the personal trading commission rate of a symbol
//
// Requires the "Place/cancel orders" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-trading-commission
//
// Arguments:
//  Symbol(string)  // The symbol of the commission rate
func (client *Client) GetTradingCommissionOfSymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.TradingCommission, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbol)
	if err != nil {
		return
	}
	err = client.privateGet(
		ctx,
		endpointTradingCommission+"/"+params[internal.ArgNameSymbol].(string),
		nil,
		&result,
	)
	if err != nil {
		return
	}
	result.Symbol = params[internal.ArgNameSymbol].(string)
	return
}

/////////////////////
// Trading history //
/////////////////////

// GetSpotOrderHistory gets all the spot orders
//
// Orders without executions are deleted after 24 hours
//
// 'from' param and 'till' param must have the same format, both id or both timestamp
//
// Requires the "Orderbook, History, Trading balance" API key Access Right
//
// https://api.exchange.cryptomkt.com/#spot-orders-history
//
// Arguments:
//  Symbol(string)  // Optional. Filter orders by symbol
//  SortBy(SortByType)  // Optional. Sorting parameter. SortByID or SortByTimestamp. Default is SortByTimestamp
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  From(string)  // Optional. Initial value of the queried interval
//  Till(string)  // Optional. Last value of the queried interval
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 100. Max is 1000
//  Offset(int)  // Optional. Default is 0. Max is 100000
func (client *Client) GetSpotOrdersHistory(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Order, err error) {
	params, err := args.BuildParams(arguments)
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointOrderHistory, params, &result)
	return
}

// GetSpotTradesHistory gets the user's spot trading history
//
// Requires the "Orderbook, History, Trading balance" API key Access Right
//
// https://api.exchange.cryptomkt.com/#spot-trades-history
//
// Arguments:
//  OrderID(string)  // Optional. Order unique identifier as assigned by the exchange
//  Symbol(string)  // Optional. Filter orders by symbol
//  SortBy(SortByType)  // Optional. Sorting parameter. SortByID or SortByTimestamp. Default is SortByTimestamp
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  From(string)  // Optional. Initial value of the queried interval
//  Till(string)  // Optional. Last value of the queried interval
//  Limit(int)  // Optional. Prices per currency pair. Defaul is 100. Max is 1000
//  Offset(int)  // Optional. Default is 0. Max is 100000
func (client *Client) GetSpotTradesHistory(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Trade, err error) {
	params, err := args.BuildParams(arguments)
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTradeHistory, params, &result)
	return
}

///////////////////////
// WALLET MANAGAMENT //
///////////////////////

// GetWalletBalance gets the user's wallet balance for all currencies with balance
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#wallet-balance
func (client *Client) GetWalletBallances(
	ctx context.Context,
) (result []models.Balance, err error) {
	err = client.privateGet(ctx, endpointWalletBalance, nil, &result)
	return
}

// GetWalletBalanceOfCurrency gets the user's wallet balance of a currency
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#wallet-balance
//
// Arguments:
//  Currency(string)  // The currency code to query the balance
func (client *Client) GetWalletBalanceOfCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Balance, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	err = client.privateGet(
		ctx,
		endpointWalletBalance+"/"+params["currency"].(string),
		nil,
		&result)
	if err != nil {
		return
	}
	result.Currency = params[internal.ArgNameCurrency].(string)
	return
}

// GetDepositCryptoAddresses gets a list of addresses with the current addresses of all currencies
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#deposit-crypto-address
func (client *Client) GetDepositCryptoAddresses(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments)
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointCryptoAdress, params, &result)
	return
}

// GetDepositCryptoAddressOfCurrency gets the current addresses of a currency of the user
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#deposit-crypto-address
//
// Arguments:
//  Currency(string)  // Currency to gets the address
func (client *Client) GetDepositCryptoAddressOfCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	resposne := make([]models.CryptoAddress, 0)
	err = client.privateGet(ctx, endpointCryptoAdress, params, &resposne)
	if len(resposne) < 1 {
		return result, errors.New("CryptomarketSDKError: no such address")
	}
	result = &resposne[0]
	return
}

// CreateDepositCryptoAddress Creates a new address for a currency
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#deposit-crypto-address
//
// Arguments:
//  Currency(string)  // currency to create a new address
func (client *Client) CreateDepositCryptoAddress(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	err = client.post(
		ctx,
		endpointCryptoAdress,
		params,
		&result,
	)
	return
}

// GetLast10DepositCryptoAddresses gets the last 10 unique addresses used for deposit, by currency
//
// Addresses used a long time ago may be omitted, even if they are among the last 10 unique addresses
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#last-10-deposit-crypto-address
//
// Arguments:
//  Currency(string)  // currency to gets the list of addresses
func (client *Client) GetLast10DepositCryptoAddresses(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	err = client.privateGet(
		ctx,
		endpointCryptoAdressRecentDeposit,
		params,
		&result,
	)
	return
}

// GetLast10WithdrawalCryptoAddresses gets the last 10 unique addresses used for withdrawals, by currency
//
// Addresses used a long time ago may be omitted, even if they are among the last 10 unique addresses
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#last-10-withdrawal-crypto-addresses
//
// Arguments:
//  Currency(string)  // currency to gets the list of addresses
func (client *Client) GetLast10WithdrawalCryptoAddresses(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameCurrency)
	if err != nil {
		return
	}
	err = client.privateGet(
		ctx,
		endpointCryptoAdressRecentWithdraw,
		params,
		&result,
	)
	return
}

// WithdrawCrypto withdraws crypto to the given address, and returns the transaction identifier
//
// Please take note that changing security settings affects withdrawals:
//
// - It is impossible to withdraw funds without enabling the two-factor authentication (2FA)
//
// - Password reset blocks withdrawals for 72 hours
//
// - Each time a new address is added to the whitelist, it takes 48 hours before that address becomes active for withdrawal
//
// Successful response to the request does not necessarily mean the resulting transaction got executed immediately. It has to be processed first and may eventually be rolled back
//
// To see whether a transaction has been finalized, call GetTransaction() with the corresponding ID
//
// Requires the "Withdraw cryptocurrencies" API key Access Right
//
// https://api.exchange.cryptomkt.com/#withdraw-crypto
//
// Arguments:
//  Currency(string)  // currency code of the crypto to withdraw
//  Amount(string)  // amount to be sent to the specified address
//  Address(string)  // address identifier
//  PaymentID(string)  // Optional.
//  IncludeFee(bool)  // Optional. If true then the amount includes fees. Default is false
//  AutoCommit(bool)  // Optional. If false then you should commit or rollback the transaction in an hour. Used in two phase commit schema. Default is true
//  UseOffchain(UseOffchainType)  // Optional. Whether the withdrawal may be comitted offchain. Accepted values are UseOffchainNever, UseOffchainOptionaly and UseOffChainRequired
//  PublicComment(string)  // Optional. Maximum lenght is 255
func (client *Client) withdrawCrypto(
	ctx context.Context,
	arguments ...args.Argument,
) (result string, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameCurrency,
		internal.ArgNameAmount,
		internal.ArgNameAddress,
	)
	if err != nil {
		return
	}
	response := models.IDResponse{}
	err = client.post(ctx, endpointCryptoWithdraw, params, &response)
	result = response.ID
	return
}

// WithdrawCryptoCommit commits a withdrawal
//
// Requires the "Withdraw cryptocurrencies" API key Access Right
//
// https://api.exchange.cryptomkt.com/#withdraw-crypto-commit-or-rollback
//
// Arguments:
//  ID(string)  // the withdrawal transaction identifier
func (client *Client) WithdrawCryptoCommit(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameID)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.put(
		ctx,
		endpointCryptoWithdraw+"/"+params["id"].(string),
		nil,
		&response,
	)
	result = response.Result
	return
}

// WithdrawCryptoRollback Rollback a withdrawal
//
// Requires the "Withdraw cryptocurrencies" API key Access Right
//
// https://api.exchange.cryptomkt.com/#withdraw-crypto-commit-or-rollback
//
// Arguments:
//  ID(string)  // the withdrawal transaction identifier
func (client *Client) WithdrawCryptoRollback(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameID)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.delete(
		ctx,
		endpointCryptoWithdraw+"/"+params["id"].(string),
		nil,
		&response,
	)
	result = response.Result
	return
}

// GetEstimateWithdrawalFee gets an estimate of the withdrawal fee
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#estimate-withdraw-fee
//
// Arguments:
//  Currency(string)  // the currency code for withdrawal
//  Amount(string)  // the expected withdraw amount
func (client *Client) GetEstimateWithdrawFee(
	ctx context.Context,
	arguments ...args.Argument,
) (result string, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameCurrency,
		internal.ArgNameAmount,
	)
	if err != nil {
		return
	}
	response := models.FeeResponse{}
	err = client.privateGet(ctx, endpointEstimateWithdrawFee, params, &response)
	result = response.Fee
	return
}

// ConvertBetweenCurrencies Converts between currencies
//
// Successful response to the request does not necessarily mean the resulting transaction got executed immediately. It has to be processed first and may eventually be rolled back
//
// To see whether a transaction has been finalized, call
//  getTransaction(id string)
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#convert-between-currencies
//
// Arguments:
//  FromCurrency(string)  // currency code of origin
//  ToCurrency(string)  // currency code of destiny
//  Amount(string)  // the amount to be converted
func (client *Client) ConvertBetweenCurrencies(
	ctx context.Context,
	arguments ...args.Argument,
) (result []string, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameFromCurrency,
		internal.ArgNameToCurrency,
		internal.ArgNameAmount,
	)
	if err != nil {
		return
	}
	response := models.ResultListResponse{}
	err = client.post(ctx, endpointConvert, params, &response)
	result = response.IDs
	return
}

// CheckIfCryptoAddressBelongsToCurrentAccount checks if an address is from this account
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#check-if-crypto-address-belongs-to-current-account
//
// Arguments:
//  Address(string)  // address to check
func (client *Client) CheckIfCryptoAddressBelongsToCurrentAccount(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameAddress)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.privateGet(
		ctx,
		endpointCryptoAdressCheckMine,
		params,
		&response,
	)
	result = response.Result
	return
}

// TransferBetweenWalletAndExchange Transfer funds between account types
//
// Source param and Destination params must be different account types
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#transfer-between-wallet-and-exchange
//
// Arguments:
//  Currency(string)  // currency code for transfering
//  Amount(string)  // amount to be transfered
//  Source(AccountType)  // transfer source account type. Either AccountWallet or AccountSpot
//  Destination(AccountType)  // transfer source account type. Either AccountWallet or AccountSpot
func (client *Client) TransferBetweenWalletAndExchange(
	ctx context.Context,
	arguments ...args.Argument,
) (transactionID string, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameCurrency,
		internal.ArgNameAmount,
		internal.ArgNameSource,
		internal.ArgNameDestination,
	)
	if err != nil {
		return
	}
	response := []string{}
	err = client.post(
		ctx,
		endpointWalletTranser,
		params,
		&response,
	)
	if err != nil {
		return
	}
	if len(response) < 1 {
		return transactionID, errors.New("CryptomarketSDKError: Bad response format")
	}
	transactionID = response[0]
	return
}

// TransferMoneyToAnotherUser Transfer funds to another user
//
// Requires the "Withdraw cryptocurrencies" API key Access Right
//
// https://api.exchange.cryptomkt.com/#transfer-money-to-another-user
//
// Arguments:
//  Currency(string)  // currency code
//  Amount(string)  // amount to be transfered
//  IdentifyBy(IdentifyByType)  // type of identifier. Either IdentifyByEmail or IdentifyByUsername
//  Identifier(string)  // the email or username of the recieving user
func (client *Client) TransferMoneyToAnotherUser(
	ctx context.Context,
	arguments ...args.Argument,
) (result string, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameCurrency,
		internal.ArgNameAmount,
		internal.ArgNameIdentifyBy,
		internal.ArgNameIdentifier)
	if err != nil {
		return
	}
	response := models.ResultResponse{}
	err = client.post(ctx, endpointInternalWithdraw, params, &response)
	result = response.ID
	return
}

// GetTransactionHistory gets the transaction history of the account
//
// Important:
//
//  - The list of supported transaction types may be expanded in future versions
//
//  - Some transaction subtypes are reserved for future use and do not purport to provide any functionality on the platform
//
//  - The list of supported transaction subtypes may be expanded in future versions
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-transactions-history
//
// Arguments:
//  TransactionIds([]string)  // Optional. List of transaction identifiers to query
//  TransactionTypes([]TransactionType)  // Optional. List of types to query. valid types are: TransactionDeposit, TransactionWithdraw, TransactionTransfer and TransactionSwap
//  TransactionSubTypes([]TransactionSubType)  // Optional. List of subtypes to query. valid subtypes are: TransactionSubTypeUnclassified, TransactionSubTypeBlockchain,  TransactionSubTypeAffiliate,  TransactionSubtypeOffchain, TransactionSubTypeFiat, TransactionSubTypeSubAccount, TransactionSubTypeWalletToSpot, TransactionSubTypeSpotToWallet, TransactionSubTypeChainSwitchFrom and TransactionSubTypeChainSwitchTo
//  TransactionStatuses([]TransactionStatusType)  // Optional. List of statuses to query. valid subtypes are: TransactionStatusCreated, TransactionStatusPending, TransactionStatusFailed, TransactionStatusSuccess and TransactionStatusRolledBack
//  SortBy(SortByType)  // Optional. sorting parameter. SortByCreatedAt or SortByID. Default is SortByCreatedAt
//  From(string)  // Optional. Interval initial value when ordering by CreatedAt. As Datetime.
//  Till(string)  // Optional. Interval end value when ordering by CreatedAt. As Datetime.
//  IDFrom(string)  // Optional. Interval initial value when ordering by id. Min is 0
//  IDTill(string)  // Optional. Interval end value when ordering by id. Min is 0
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  Limit(int)  // Optional. Transactions per query. Defaul is 100. Max is 1000
//  Offset(int)  // Optional. Default is 0. Max is 100000
func (client *Client) GetTransactionHistory(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Transaction, err error) {
	params, err := args.BuildParams(arguments)
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTransactions, params, &result)
	return
}

// GetTransaction gets a transaction by its identifier
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-transactions-history
//
// Arguments:
//  ID(string)  // The identifier of the transaction
func (client *Client) GetTransaction(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameID)
	if err != nil {
		return
	}
	err = client.privateGet(
		ctx,
		endpointTransactions+"/"+params["id"].(string),
		nil,
		&result,
	)
	return
}

// CheckIfOffchainIsAvailable gets the status of the offchain
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#check-if-offchain-is-available
//
// Arguments:
//  Currency(string)  // currency code
//  Address(string)  // address identifier
//  PaymentID(string)  // Optional.
func (client *Client) CheckIfOffchainIsAvailable(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(
		arguments,
		internal.ArgNameCurrency,
		internal.ArgNameAddress,
	)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.post(
		ctx,
		endpointCryptoCheckOffchain,
		params,
		&response,
	)
	result = response.Result
	return
}

// GetAmountLocks gets the list of amount locks
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-amount-locks
//
// Arguments:
//  Currency(string)  // Optional. Currency code
//  Active(bool)  // Optional. value showing whether the lock is active
//  Limit(int)  // Optional. Dafault is 100. Min is 0. Max is 1000
//  Offset(int)  // Optional. Default is 0. Min is 0
//  From(string)  // Optional. Interval initial value. As Datetime
//  Till(string)  // Optional. Interval end value. As Datetime
func (client *Client) GetAmountLocks(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.AmountLock, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(
		ctx,
		endpointAmountLocks,
		params,
		&result,
	)
	return
}

//////////////////
// sub accounts //
//////////////////

// GetSubAccounts gets the list of sub-accounts
//
// Requires no API key Access Rights.
//
// https://api.exchange.cryptomkt.com/#get-sub-accounts-list
func (client *Client) GetSubAccounts(
	ctx context.Context,
) (result []models.SubAccount, err error) {
	var response struct {
		Result []models.SubAccount `json:"result"`
	}
	err = client.privateGet(ctx, endpointSubAccountList, nil, &response)
	return response.Result, err
}

// FreezeSubAccounts freezes sub-accounts listed and returnes a boolean indicating whether the sub-accounts where frozen
//
// A frozen wouldn't be able to:
//
// - login
//
// - withdraw funds
//
// - trade
//
// - complete pending orders
//
// - use API keys
//
//
// For any sub-account listed, all orders will be canceled and all funds will be transferred form the Trading balance
//
// Requires no API key Access Rights. Requires to be authenticated
//
// https://api.exchange.cryptomkt.com/#freeze-sub-account
//
// Arguments:
//  SubAccountIDs(...string)  // A list of sub account ids. Ids as hexadecimal code
func (client *Client) FreezeSubAccounts(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountIDs)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.post(ctx, endpointFreezeSubAccount, params, &response)
	return response.Result, err
}

// ActivateSubAccounts  unfreezes sub-accounts listed and returns a bool indicating whether the sub accounts where activated
//
// Requires no API key Access Rights. Requires to be authenticated
//
// https://api.exchange.cryptomkt.com/#activate-sub-account
//
// Arguments:
//  SubAccountIDs(...string)  // currency code of the sub-accounts to activate. Ids as hexadecimal code
func (client *Client) ActivateSubAccounts(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountIDs)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.post(ctx, endpointActivateSubAccount, params, &response)
	return response.Result, err
}

// TransferFunds transfers funds from the super-account to a sub-account or from a sub-account to the super-account, and returns the transaction id
//
// Requires the "Withdraw cryptocurrencies" API key Access Right
//
// https://api.exchange.cryptomkt.com/#transfer-funds
//
// Arguments:
//  SubAccountID(string)  // id of the sub-account to transfer with the super-account
//  Amount(string)  // amount of funds to transfer
//  Currency(string)  // currency of transfer
//  TransferType(TransferTypeType)  // TransferToSubAccount or TransferFromSubAccount
func (client *Client) TransferFunds(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountID)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.post(ctx, endpointSubAccountTransferFunds, params, &response)
	return response.Result, err
}

// GetACLSettings get a list of withdrawal settings for all sub-accounts or for the specified sub-accounts and returns A list of ACL settings for subaccounts
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#get-acl-settings
//
// Arguments:
//  SubAccountIDs(...string)  // A list of sub account ids. Ids as hexadecimal code
func (client *Client) GetACLSettings(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.ACLSettings, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountIDs)
	if err != nil {
		return
	}
	type Response struct {
		Result []models.ACLSettings `json:"result"`
	}
	var response Response
	err = client.privateGet(ctx, endpointSubaccountACLSettings, params, &response)
	return response.Result, err
}

// ChangeACLSettings changes the ACL settings of subaccounts and returns a list of acl settings of the changed sub-accounts
//
// Disables or enables withdrawals for a sub-account
//
// Requires the "Payment information" API key Access Right
//
// https://api.exchange.cryptomkt.com/#change-acl-settings
//
// Arguments:
//  SubAccountIDs(...string)  // currency code for transfering
//  DepositAddressGenerationEnabled(bool)	// Optional. value indicaiting permission for deposits
//  WithdrawEnabled(bool)  // Optional. value indicating permission for withdrawals
//  Description(string)  // Optional. Textual description.
//  CreatedAt(string)  // Optional. ACL creation time
//  UpdatedAt(string)  // Optional. ACL update time
func (client *Client) ChangeACLSettings(
	ctx context.Context,
	arguments ...args.Argument,
) (result bool, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountID)
	if err != nil {
		return
	}
	response := models.BooleanResponse{}
	err = client.post(ctx, endpointSubaccountACLSettings, params, &response)
	return response.Result, err
}

// GetSubAccountBalances Get the non-zero balances of a sub-account
//
// Report will include the wallet and Trading balances.
//
// Works independent of account state.
//
// Requires the "Payment information" API key Access Right.
//
// https://api.exchange.cryptomkt.com/#get-sub-account-balance
//
// Arguments:
//  SubAccountID(string)  // id of the sub-account to get the balances
func (client *Client) GetSubAccountBalances(
	ctx context.Context,
	arguments ...args.Argument,
) (result models.SubAccountBalances, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountID)
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointSubaccountBalance+"/"+params[internal.ArgNameSubAccountID].(string), nil, &result)
	return
}

// GetSubAccountCryptoAddress get the crypto address of the sub-account
//
// Requires the "Payment information" API key Access Right.
//
// https://api.exchange.cryptomkt.com/#get-sub-account-crypto-address
//
// Arguments:
//  SubAccountID(string)  // the sub-account id
//  Currency(string)  // the currency code of the crypto address
func (client *Client) GetSubAccountCryptoAddress(
	ctx context.Context,
	arguments ...args.Argument,
) (result string, err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubAccountID)
	if err != nil {
		return
	}
	type Result struct {
		Address string `json:"address"`
	}
	type Response struct {
		Result Result `json:"result"`
	}
	response := Response{}
	err = client.privateGet(
		ctx,
		endpointSubaccountCryptoAddress+"/"+params[internal.ArgNameSubAccountID].(string)+"/"+params[internal.ArgNameCurrency].(string),
		nil,
		&response,
	)
	return response.Result.Address, err
}
