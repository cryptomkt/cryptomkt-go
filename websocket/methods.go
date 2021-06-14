package websocket

import (
	"strings"

	"github.com/cryptomarket/cryptomarket-go/args"
)

const (
	methodGetCurrencies = "getCurrencies"
	methodGetCurrency   = "getCurrency"
	methodGetSymbols    = "getSymbols"
	methodGetSymbol     = "getSymbol"
	methodGetTrades     = "getTrades"

	methodNewOrder          = "newOrder"
	methodCancelOrder       = "cancelOrder"
	methodReplaceOrder      = "cancelReplceOrder"
	methodGetOrders         = "getOrders"
	methodGetTradingBalance = "getTradingBalance"

	methodGetBalance       = "getBalance"
	methodFindTransactions = "findTransactions"
	methodLoadTransactions = "loadTransactions"

	methodSubscribeTicker  = "subscribeTicker"
	methodUnsubcribeTicker = "unsubscribeTicker"
	methodTicker           = "ticker"

	methodSubscribeOrderbook   = "subscribeOrderbook"
	methodUnsubscribeOrderbook = "unsubscribeOrderbook"
	methodUpdateOrderbook      = "updateOrderbook"
	methodSnapshotOrderbook    = "snapshotOrderbook"

	methodSubscribeTrades   = "subscribeTrades"
	methodUnsubscribeTrades = "unsubscribeTrades"
	methodSnapshotTrades    = "snapshotTrades"
	methodUpdateTrades      = "updateTrades"

	methodSubscribeCandles   = "subscribeCandles"
	methodUnsubscribeCandles = "unsubscribeCandles"
	methodSnapshotCandles    = "snapshotCandles"
	methodUpdateCandles      = "updateCandles"

	methodSubscribeReports = "subscribeReports"
	methodActiveOrders     = "activeOrders"
	methodReport           = "report"

	methodSubscribeTransactions   = "subscribeTransactions"
	methodUnsubscribeTransactions = "unsubscribeTransactions"
	methodUpdateTransaction       = "updateTransaction"
)

const (
	ticker       = "ticker"
	orderbook    = "orderbook"
	trades       = "trades"
	candles      = "candles"
	reports      = "reports"
	transactions = "transactions"
)

var methodMapping = map[string]string{
	methodSubscribeTicker:  ticker,
	methodUnsubcribeTicker: ticker,
	methodTicker:           ticker,

	methodSubscribeOrderbook:   orderbook,
	methodUnsubscribeOrderbook: orderbook,
	methodUpdateOrderbook:      orderbook,
	methodSnapshotOrderbook:    orderbook,

	methodSubscribeTrades:   trades,
	methodUnsubscribeTrades: trades,
	methodSnapshotTrades:    trades,
	methodUpdateTrades:      trades,

	methodSubscribeCandles:   candles,
	methodUnsubscribeCandles: candles,
	methodSnapshotCandles:    candles,
	methodUpdateCandles:      candles,

	methodSubscribeReports: reports,
	methodActiveOrders:     reports,
	methodReport:           reports,

	methodSubscribeTransactions:   transactions,
	methodUnsubscribeTransactions: transactions,
	methodUpdateTransaction:       transactions,
}

func orderbookFeed(method string) bool {
	return methodMapping[method] == orderbook
}

func tradesFeed(method string) bool {
	return methodMapping[method] == trades
}

func candlesFeed(method string) bool {
	return methodMapping[method] == candles
}

func buildKey(method, symbol, period string) string {
	methodKey := methodMapping[method]
	var key string
	if method == methodReport || method == methodActiveOrders {
		key = methodKey + "::"
	} else {
		if methodKey == candles && period == "" { // default period
			period = string(args.PeriodType30Minutes)
		}
		key = methodKey + ":" + symbol + ":" + period
	}
	return strings.ToUpper(key)
}

func buildKeyFromParams(method string, params map[string]interface{}) string {
	// the key part
	symbol := ""
	if s, ok := params["symbol"]; ok {
		symbol = s.(string)
	}
	period := ""
	if p, ok := params["period"]; ok {
		period = p.(string)
	}
	return buildKey(method, symbol, period)
}

func buildKeyFromResponse(response wsResponse) string {
	if response.Method == "" {
		return ""
	}
	return buildKey(response.Method, response.Params.Symbol, response.Params.Period)
}
