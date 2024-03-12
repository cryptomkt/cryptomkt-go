package rest

import (
	"context"

	"github.com/cryptomkt/cryptomkt-go/v3/args"
	"github.com/cryptomkt/cryptomkt-go/v3/models"
)

// Public - Market data

// GetTickerBySymbol is an alias of GetTickerOfSymbol
func (client *Client) GetTickerBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Ticker, err error) {
	return client.GetTickerOfSymbol(ctx, arguments...)
}

// GetTicker is an alias of GetTickerOfSymbol
func (client *Client) GetTicker(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Ticker, err error) {
	return client.GetTickerOfSymbol(ctx, arguments...)
}

// GetTickerLastPricesBySymbol is an alias of GetTickerLastPricesOfSymbol
func (client *Client) GetTickerLastPricesBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Price, err error) {
	return client.GetTickerLastPricesOfSymbol(ctx, arguments...)
}

// GetTickerLastPrice is an alias of GetTickerLastPricesOfSymbol
func (client *Client) GetTickerLastPrice(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Price, err error) {
	return client.GetTickerLastPricesOfSymbol(ctx, arguments...)
}

// GetTradesBySymbol is an alias of GetTradesOfSymbol
func (client *Client) GetTradesBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.PublicTrade, err error) {
	return client.GetTradesOfSymbol(ctx, arguments...)
}

// GetOrderBookBySymbol is an alias of GetOrderBookOfSymbol
func (client *Client) GetOrderBookBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.OrderBook, err error) {
	return client.GetOrderBookOfSymbol(ctx, arguments...)
}

// GetOrderBook is an alias of GetOrderBookOfSymbol
func (client *Client) GetOrderBook(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.OrderBook, err error) {
	return client.GetOrderBookOfSymbol(ctx, arguments...)
}

// GetOrderBookVolumeBySymbol is an alias of GetOrderBookVolumeOfSymbol
func (client *Client) GetOrderBookVolumeBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.OrderBook, err error) {
	return client.GetOrderBookVolumeOfSymbol(ctx, arguments...)
}

// GetOrderBookVolume is an alias of GetOrderBookVolumeOfSymbol
func (client *Client) GetOrderBookVolume(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.OrderBook, err error) {
	return client.GetOrderBookVolumeOfSymbol(ctx, arguments...)
}

// GetCandlesBySymbol is an alias of GetCandlesOfSymbol
func (client *Client) GetCandlesBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result []models.Candle, err error) {
	return client.GetCandlesOfSymbol(ctx, arguments...)
}

// GetConvertedCandlesBySymbol is an alias of GetConvertedCandlesOfSymbol
func (client *Client) GetConvertedCandlesBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result models.ConvertedCandlesOfSymbol, err error) {
	return client.GetConvertedCandlesOfSymbol(ctx, arguments...)
}

// Private - Spot Trading

// GetSpotTradingBalanceByCurrency is an alias of GetSpotTradingBalanceOfCurrency
func (client *Client) GetSpotTradingBalanceByCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Balance, err error) {
	return client.GetSpotTradingBalanceOfCurrency(ctx, arguments...)
}

// GetSpotTradingBalance is an alias of GetSpotTradingBalanceOfCurrency
func (client *Client) GetSpotTradingBalance(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Balance, err error) {
	return client.GetSpotTradingBalanceOfCurrency(ctx, arguments...)
}

// GetTradingCommissions is an alias of GetAllTradingCommissions
func (client *Client) GetTradingCommissions(
	ctx context.Context,
) (result []models.TradingCommission, err error) {
	return client.GetAllTradingCommissions(ctx)
}

// GetTradingCommissionBySymbol is an alias of GetTradingCommissionOfSymbol
func (client *Client) GetTradingCommissionBySymbol(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.TradingCommission, err error) {
	return client.GetTradingCommissionOfSymbol(ctx, arguments...)
}

// GetTradingCommissionBySymbol is an alias of GetTradingCommissionOfSymbol
func (client *Client) GetTradingCommission(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.TradingCommission, err error) {
	return client.GetTradingCommissionOfSymbol(ctx, arguments...)
}

// Private - Wallet Management

// GetWalletBalanceOfCurrency is an alias of GetWalletBalanceOfCurrency
func (client *Client) GetWalletBalanceByCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Balance, err error) {
	return client.GetWalletBalanceOfCurrency(ctx, arguments...)
}

// GetWalletBalance is an alias of GetWalletBalanceOfCurrency
func (client *Client) GetWalletBalance(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.Balance, err error) {
	return client.GetWalletBalanceOfCurrency(ctx, arguments...)
}

// GetDepositCryptoAddressByCurrency is an alias of GetDepositCryptoAddressOfCurrency
func (client *Client) GetDepositCryptoAddressByCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.CryptoAddress, err error) {
	return client.GetDepositCryptoAddressOfCurrency(ctx, arguments...)
}

// GetDepositCryptoAddress is an alias of GetDepositCryptoAddressOfCurrency
func (client *Client) GetDepositCryptoAddress(
	ctx context.Context,
	arguments ...args.Argument,
) (result *models.CryptoAddress, err error) {
	return client.GetDepositCryptoAddressOfCurrency(ctx, arguments...)
}
