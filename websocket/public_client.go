package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

// PublicClient connects via websocket to cryptomarket to get market information of the exchange.
type PublicClient struct {
	clientBase
	obCache *orderbookCache
}

// NewPublicClient returns a new chan client if the connection with the
// cryptomarket server is successful, and error otherwise.
func NewPublicClient() (*PublicClient, error) {
	methodMapping := map[string]string{
		// tickers
		"subscribeTicker":   "tickers",
		"unsubscribeTicker": "tickers",
		"ticker":            "tickers",
		// orderbooks
		"subscribeOrderbook":   "orderbooks",
		"unsubscribeOrderbook": "orderbooks",
		"snapshotOrderbook":    "orderbooks",
		"updateOrderbook":      "orderbooks",
		// trades
		"subscribeTrades":   "trades",
		"unsubscribeTrades": "trades",
		"snapshotTrades":    "trades",
		"updateTrades":      "trades",
		// candles
		"subscribeCandles":   "candles",
		"unsubscribeCandles": "candles",
		"snapshotCandles":    "candles",
		"updateCandles":      "candles",
	}
	keyFunc := func(method string, params map[string]interface{}) string {
		methodKey := methodMapping[method]
		period, _ := params["period"].(string)

		if methodKey == "candles" && period == "" { // default period
			period = string(args.PeriodType30Minutes)
		}
		symbol, _ := params["symbol"].(string)
		key := methodKey + ":" + symbol + ":" + period
		return strings.ToUpper(key)
	}
	client := &PublicClient{
		clientBase: clientBase{
			wsManager: newWSManager("/api/2/ws/public"),
			chanCache: newChanCache(),
			subscriptionKeysFunc: func(method string, params map[string]interface{}) (string, bool) {
				return keyFunc(method, params), true
			},
			keyFromResponse: func(response wsResponse) string {
				return keyFunc(
					response.Method,
					map[string]interface{}{
						"symbol": response.Params.Symbol,
						"period": response.Params.Period,
					},
				)
			},
		},
		obCache: newOrderbookCache(),
	}

	// connect to streaming
	err := client.wsManager.connect()
	if err != nil {
		return nil, fmt.Errorf("Error in websocket client connection: %s", err)
	}
	// handle incomming data
	go client.handle(client.wsManager.rcv)
	return client, nil
}

// GetCurrencies gets a list all available currencies on the exchange
//
//https://api.exchange.cryptomarket.com/#Get-currencies
func (client *PublicClient) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	var resp struct {
		Result []models.Currency
	}
	err := client.doRequest(ctx, methodGetCurrencies, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// GetCurrency gets the data of a currency
//
// https://api.exchange.cryptomarket.com/#get-currencies
//
// Arguments:
//  Currency(string) // the currency id
func (client *PublicClient) GetCurrency(ctx context.Context, arguments ...args.Argument) (*models.Currency, error) {
	var resp struct {
		Result models.Currency
	}
	err := client.doRequest(ctx, methodGetCurrency, arguments, []string{"currency"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

//GetSymbols gets a list of the specified symbols or all of them if no symbols are specified.
// A symbol is the combination of the base currency (first one) and quote currency (second one).
//
// https://api.exchange.cryptomarket.com/#get-symbols
func (client *PublicClient) GetSymbols(ctx context.Context) ([]models.Symbol, error) {
	var resp struct {
		Result []models.Symbol
	}
	err := client.doRequest(ctx, methodGetSymbols, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// GetSymbol gets a symbol by its id.
// A symbol is the combination of the base currency (first one) and quote currency (second one).
//
// https://api.exchange.cryptomarket.com/#get-symbols
//
// Arguments:
//  Symbol(string) // The symbol id
func (client *PublicClient) GetSymbol(ctx context.Context, arguments ...args.Argument) (*models.Symbol, error) {
	var resp struct {
		Result models.Symbol
	}
	err := client.doRequest(ctx, methodGetSymbol, arguments, []string{"symbol"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// GetTrades gets trades of the specified symbol.
//
// https://api.exchange.cryptomarket.com/#get-trades
//
// Arguments:
//  Symbol(string) // The symbol to get the trades
//  Sort(SortType) // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  From(string)   // Optional. Initial value of the queried interval
//  Till(string)   // Optional. Last value of the queried interval
//  Limit(int)     // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)    // Optional. Default is 0. Max is 100000
func (client *PublicClient) GetTrades(ctx context.Context, arguments ...args.Argument) (trades []models.PublicTrade, err error) {
	var resp struct {
		Result struct {
			Data []models.PublicTrade
		}
	}
	err = client.doRequest(ctx, methodGetTrades, arguments, []string{"symbol"}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result.Data, nil
}

///////////////////
// Subscriptions //
///////////////////

// SubscribeToTicker subscribes to a ticker of a symbol.
//
// https://api.exchange.cryptomarket.com/#subscribe-to-ticker
//
// Arguments:
//  Symbol(string) // The symbol of the ticker to subscribe
func (client *PublicClient) SubscribeToTicker(arguments ...args.Argument) (feedCh chan models.Ticker, err error) {
	dataCh, err := client.doSubscription(methodSubscribeTicker, arguments, []string{"symbol"})
	if err != nil {
		return nil, err
	}
	feedCh = make(chan models.Ticker)
	go func() {
		defer close(feedCh)
		var resp struct {
			Params models.Ticker
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			feedCh <- resp.Params
		}
	}()
	return feedCh, nil
}

// UnsubscribeToTicker unsubscribes to a ticker of a symbol.
// It also closes the feedCh of the subscription.
//
// https://api.exchange.cryptomarket.com/#subscribe-to-ticker
//
// Arguments:
//  Symbol(string) // The symbol of the ticker to unsubscribe
func (client *PublicClient) UnsubscribeToTicker(arguments ...args.Argument) error {
	return client.doUnsubscription(methodUnsubcribeTicker, arguments, []string{"symbol"})
}

// SubscribeToOrderbook subscribes to the order book of a symbol.
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level.
//
// https://api.exchange.cryptomarket.com/#subscribe-to-order-book
//
// Arguments:
//  Symbol(string) // The symbol of the orderbook to subscribe
func (client *PublicClient) SubscribeToOrderbook(arguments ...args.Argument) (chan models.OrderBook, error) {
	dataCh, err := client.doSubscription(methodSubscribeOrderbook, arguments, []string{"symbol"})
	if err != nil {
		return nil, err
	}
	feedCh := make(chan models.OrderBook)
	go func() {
		obCache := newOrderbookCache()
		defer close(feedCh)
		var resp struct {
			Method string
			Params struct {
				Symbol string
			}
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			if resp.Method == methodSnapshotOrderbook {
				obCache.obSnapshot(data)
			} else {
				obCache.obUpdate(data)
			}
			if obCache.obBroken() {
				obCache.waitOBSnapshot()
				notification := wsNotification{
					ID:     client.chanCache.nextID(),
					Method: methodSubscribeOrderbook,
					Params: map[string]interface{}{"symbol": resp.Params.Symbol},
				}
				requestData, _ := json.Marshal(notification)
				client.wsManager.snd <- requestData
			}
			if obCache.obWaiting() {
				continue
			}
			var orderbook models.OrderBook
			json.Unmarshal(obCache.getOrderbookData(), &orderbook)
			feedCh <- orderbook
		}
	}()
	return feedCh, nil
}

// UnsubscribeToOrderbook unsubscribes to an order book of a symbol.
// It also closes the feedCh of the subscription.
//
// An Order Book is an electronic list of buy and sell orders for a specific symbol, structured by price level
//
// https://api.exchange.cryptomarket.com/#subscribe-to-order-book
//
// Arguments:
//  Symbol(string) // The symbol of the orderbook to unsubscribe
func (client *PublicClient) UnsubscribeToOrderbook(arguments ...args.Argument) error {
	return client.doUnsubscription(methodUnsubscribeOrderbook, arguments, []string{"symbol"})
}

// SubscribeToTrades subscribes to the trades of a symbol
//
// https://api.exchange.cryptomarket.com/#subscribe-to-trades
//
// Arguments:
//  Symbol(string) // The symbol of the trades to subscribe
//  Limit(int)     // Optional. Maximum number of trades in the first feed.
func (client *PublicClient) SubscribeToTrades(arguments ...args.Argument) (feedCh chan []models.PublicTrade, err error) {
	dataCh, err := client.doSubscription(methodSubscribeTrades, arguments, []string{"symbol"})
	if err != nil {
		return nil, err
	}
	feedCh = make(chan []models.PublicTrade)
	go func() {
		defer close(feedCh)
		var resp struct {
			Params struct {
				Data []models.PublicTrade
			}
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			feedCh <- resp.Params.Data
		}
	}()
	return feedCh, nil
}

// UnsubscribeToTrades unsubscribes to a trades of a symbol.
// It also closes the feedCh of the subscription.
//
// https://api.exchange.cryptomarket.com/#subscribe-to-trades
//
// Arguments:
//  Symbol(string) // The symbol of the trades to unsubscribe
func (client *PublicClient) UnsubscribeToTrades(arguments ...args.Argument) error {
	return client.doUnsubscription(methodUnsubscribeTrades, arguments, []string{"symbol"})
}

// SubscribeToCandles subscribes to the candles of a symbol, at the given period
//
// Candels are used for OHLC representation
//
//
// https://api.exchange.cryptomarket.com/#subscribe-to-candles
//
// Arguments:
//  Symbol(string)     // The symbol of the candles to subscribe
//  Period(PeriodType) // A valid tick interval. A PeriodType
//  Limit(int)         // Optional. Maximum number of trades in the first feed.
func (client *PublicClient) SubscribeToCandles(arguments ...args.Argument) (feedCh chan []models.Candle, err error) {
	dataCh, err := client.doSubscription(methodSubscribeCandles, arguments, []string{"symbol", "period"})
	if err != nil {
		return nil, err
	}
	feedCh = make(chan []models.Candle)
	go func() {
		defer close(feedCh)
		var resp struct {
			Params struct {
				Data []models.Candle
			}
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			feedCh <- resp.Params.Data
		}
	}()
	return feedCh, nil
}

// UnsubscribeToCandles unsubscribes to the candles of a symbol at a given period.
// It also closes the feedCh of the subscription
//
// https://api.exchange.cryptomarket.com/#subscribe-to-candles
//
// Arguments:
//  Symbol(string)     // The symbol of the candles to unsubscribe
//  Period(PeriodType) // A valid tick interval. A PeriodType
func (client *PublicClient) UnsubscribeToCandles(arguments ...args.Argument) error {
	return client.doUnsubscription(methodUnsubscribeCandles, arguments, []string{"symbol", "period"})
}
