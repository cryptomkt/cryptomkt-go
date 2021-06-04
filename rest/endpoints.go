package rest

const (
	// public endpoints
	endpointCurrency  = "public/currency"
	endpointSymbol    = "public/symbol"
	endpointTrade     = "public/trades"
	endpointOrderbook = "public/orderbook"
	endpointCandle    = "public/candles"
	endpointTicker    = "public/ticker"
	// trading endpoints
	endpointTradingBalance = "trading/balance"
	endpointOrder          = "order"
	endpointTradingFee     = "trading/fee"
	// trading history endpoints
	endpointOrderHistory = "history/order"
	endpointTradeHistory = "history/trades"
	// account management
	endpointAccountBalance         = "account/balance"
	endpointCryptoAdress           = "account/crypto/address"
	endpointCryptoAdresses         = "account/crypto/addresses"
	endpointUsedAddressed          = "account/crypto/used-addresses"
	endpointWithdrawCrypto         = "account/crypto/withdraw"
	endpointTransferConvert        = "account/crypto/transfer-convert"
	endpointEstimateWithdraw       = "account/crypto/estimate-withdraw"
	endpointCryptoAddressIsMine    = "account/crypto/is-mine"
	endpointAccountTranser         = "account/transfer"
	endpointAccountTranserInternal = "account/transfer/internal"
	endpointTransactionHistory     = "account/transactions"
)
