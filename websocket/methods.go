package websocket

const (
	methodGetCurrencies = "getCurrencies"
	methodGetCurrency   = "getCurrency"
	methodGetSymbols    = "getSymbols"
	methodGetSymbol     = "getSymbol"
	methodGetTrades     = "getTrades"

	methodNewOrder     = "newOrder"
	methodCancelOrder  = "cancelOrder"
	methodReplaceOrder = "cancelReplceOrder"
	methodGetOrders    = "getOrders"
	methodGetBalance   = "getTradingBalance"

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
)

var methodMapping = map[string]string{
	methodSubscribeTicker:  "ticker",
	methodUnsubcribeTicker: "ticker",
	methodTicker:           "ticker",

	methodSubscribeOrderbook:   "orderbook",
	methodUnsubscribeOrderbook: "orderbook",
	methodUpdateOrderbook:      "orderbook",
	methodSnapshotOrderbook:    "orderbook",

	methodSubscribeTrades:   "trades",
	methodUnsubscribeTrades: "trades",
	methodSnapshotTrades:    "trades",
	methodUpdateTrades:      "trades",

	methodSubscribeCandles:   "candles",
	methodUnsubscribeCandles: "candles",
	methodSnapshotCandles:    "candles",
	methodUpdateCandles:      "candles",

	methodSubscribeReports: "reports",
	methodActiveOrders:     "reports",
	methodReport:           "reports",
}
