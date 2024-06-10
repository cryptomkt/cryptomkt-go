package rest

const (
	// public endpoints
	endpointCurrency         = "public/currency"
	endpointSymbol           = "public/symbol"
	endpointTicker           = "public/ticker"
	endpointPrices           = "public/price/rate"
	endpointPriceHistory     = "public/price/history"
	endpointPriceTicker      = "public/price/ticker"
	endpointTrade            = "public/trades"
	endpointOrderbook        = "public/orderbook"
	endpointCandles          = "public/candles"
	endpointConvertedCandles = "public/converted/candles"
	// trading endpoints
	endpointTradingBalance    = "spot/balance"
	endpointOrder             = "spot/order"
	endpointOrderList         = "spot/order/list"
	endpointTradingCommission = "spot/fee"
	// trading history endpoints
	endpointOrderHistory = "spot/history/order"
	endpointTradeHistory = "spot/history/trade"
	// wallet management
	endpointWalletBalance              = "wallet/balance"
	endpointCryptoAdress               = "wallet/crypto/address"
	endpointCryptoAdressRecentDeposit  = "wallet/crypto/address/recent-deposit"
	endpointCryptoAdressRecentWithdraw = "wallet/crypto/address/recent-withdraw"
	endpointCryptoAdressCheckMine      = "wallet/crypto/address/check-mine"
	endpointCryptoWithdraw             = "wallet/crypto/withdraw"
	endpointConvert                    = "wallet/convert"
	endpointWalletTranser              = "wallet/transfer"
	endpointInternalWithdraw           = "wallet/internal/withdraw"
	endpointTransactions               = "wallet/transactions"
	endpointCryptoCheckOffchain        = "wallet/crypto/check-offchain-available"
	endpointEstimateWithdrawalFees     = "wallet/crypto/fees/estimate"
	endpointEstimateWithdrawalFee      = "wallet/crypto/fee/estimate"
	endpointBulkEstimateWithdrawalFees = "wallet/crypto/fee/estimate/bulk"
	endpointEstimateDepositFee         = "wallet/crypto/fee/deposit/estimate"
	endpointBulkEstimateDepositFees    = "wallet/crypto/fee/deposit/estimate/bulk"
	endpointAmountLocks                = "wallet/amount-locks"
	// sub accounts
	endpointSubAccountList          = "sub-account"
	endpointFreezeSubAccount        = "sub-account/freeze"
	endpointActivateSubAccount      = "sub-account/activate"
	endpointSubAccountTransferFunds = "sub-account/transfer"
	endpointSubaccountACLSettings   = "sub-account/acl"
	endpointSubaccountBalance       = "sub-account/balance"
	endpointSubaccountCryptoAddress = "sub-account/crypto/address"
)
