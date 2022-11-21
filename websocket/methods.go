package websocket

const (
	methodSubscriptions = "subscriptions"
	methodSubscribe     = "subscribe"

	methodSubscribeSpotReports = "spot_subscribe"
	methodSpotUnsubscribe      = "spot_unsubscribe"
	methodSpotOrder            = "spot_order"
	methodSpotOrders           = "spot_orders"

	methodGetSpotOrders       = "spot_get_orders"
	methodCreateSpotOrder     = "spot_new_order"
	methodCreateSpotOrderList = "spot_new_order_list"
	methodCancelSpotOrder     = "spot_cancel_order"
	methodCancelSpotOrders    = "spot_cancel_orders"
	methodReplaceSpotOrder    = "spot_replace_order"
	methodSpotBalance         = "spot_balance"
	methodSpotBalances        = "spot_balances"
	methodSpotFee             = "spot_fee"
	methodSpotFees            = "spot_fees"

	methodSubscribeTransactions   = "subscribe_transactions"
	methodUnsubscribeTransactions = "unsubscribe_transactions"
	methodUpdateTransactions      = "transaction_update"

	methodSubscribeWalletBalances   = "subscribe_wallet_balances"
	methodUnsubscribeWalletBalances = "unsubscribe_wallet_balances"
	methodWalletBalances            = "wallet_balances"
	methodWalletBalanceUpdate       = "wallet_balance_update"
	methodWalletBalance             = "wallet_balance"

	methodGetTransactions = "get_transactions"
)

const (
	reports        = "reports"
	transactions   = "transactions"
	walletBalances = "walletBalances"
)

var subscriptionMapping = map[string]string{
	methodSubscribeSpotReports: reports,
	methodSpotUnsubscribe:      reports,
	methodSpotOrder:            reports,
	methodSpotOrders:           reports,

	methodSubscribeWalletBalances:   walletBalances,
	methodUnsubscribeWalletBalances: walletBalances,
	methodWalletBalances:            walletBalances,
	methodWalletBalanceUpdate:       walletBalances,

	methodSubscribeTransactions:   transactions,
	methodUnsubscribeTransactions: transactions,
	methodUpdateTransactions:      transactions,
}

func getSubscriptionKey(response wsResponse) (key string, ok bool) {
	if key, ok := subscriptionMapping[response.Method]; ok {
		return key, true
	}
	if response.Ch != "" {
		return response.Ch, true
	}
	return "", false
}
