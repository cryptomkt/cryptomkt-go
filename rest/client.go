package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/cryptomkt/go-api/args"

	"github.com/cryptomkt/go-api/models"
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

// NewClient creates a new rest client to communicate with the exchange
func NewClient(apiKey, apiSecret string) (client Client) {
	client.hclient = newHTTPClient(apiKey, apiSecret)
	return
}

func (client Client) publicGet(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodGet, publicCall, endpoint, params, model)
}

func (client Client) privateGet(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodGet, privateCall, endpoint, params, model)
}

func (client Client) post(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodPost, privateCall, endpoint, params, model)
}

func (client Client) put(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodPut, privateCall, endpoint, params, model)
}

func (client Client) delete(ctx context.Context, endpoint string, params map[string]interface{}, model interface{}) error {
	return client.doRequest(ctx, methodDelete, privateCall, endpoint, params, model)
}

func (client Client) doRequest(ctx context.Context, method string, public bool, endpoint string, params map[string]interface{}, model interface{}) error {
	data, err := client.hclient.doRequest(ctx, method, endpoint, params, public)
	if err != nil {
		return err
	}
	return client.handleResponseData(data, model)
}

func (client Client) handleResponseData(data []byte, model interface{}) error {
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

func (client Client) getCurrencies(ctx context.Context, arguments ...args.Argument) (result []models.Currency, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointCurrency, params, &result)
	return
}

func (client Client) getCurrency(ctx context.Context, arguments ...args.Argument) (result *models.Currency, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointCurrency+"/"+params["currency"].(string), nil, &result)
	return
}

func (client Client) getSymbols(ctx context.Context, arguments ...args.Argument) (result []models.Symbol, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointSymbol, params, &result)
	return
}

func (client Client) getSymbol(ctx context.Context, arguments ...args.Argument) (result *models.Symbol, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointSymbol+"/"+params["symbol"].(string), nil, &result)
	return
}

func (client Client) getTickers(ctx context.Context, arguments ...args.Argument) (result []models.Ticker, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointTicker, params, &result)
	return
}

func (client Client) getTicker(ctx context.Context, arguments ...args.Argument) (result *models.Ticker, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointTicker+"/"+params["symbol"].(string), nil, &result)
	return
}

func (client Client) getTrades(ctx context.Context, arguments ...args.Argument) (result map[string][]models.PublicTrade, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointTrade, params, &result)
	return
}

func (client Client) getTradesOfSymbol(ctx context.Context, arguments ...args.Argument) (result []models.PublicTrade, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointTrade+"/"+params["symbol"].(string), params, &result)
	return
}

func (client Client) getOrderbooks(ctx context.Context, arguments ...args.Argument) (result map[string]models.OrderBook, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointOrderbook, params, &result)
	return
}

func (client Client) getOrderbook(ctx context.Context, arguments ...args.Argument) (result *models.OrderBook, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointOrderbook+"/"+params["symbol"].(string), params, &result)
	return
}

func (client Client) getCandles(ctx context.Context, arguments ...args.Argument) (result map[string][]models.Candle, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.publicGet(ctx, endpointCandle, params, &result)
	return
}

func (client Client) getCandlesOfSymbol(ctx context.Context, arguments ...args.Argument) (result []models.Candle, err error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return
	}
	err = client.publicGet(ctx, endpointCandle+"/"+params["symbol"].(string), params, &result)
	return
}

// /////////////
// // TRADING //
// /////////////

func (client Client) getTradingBalance(ctx context.Context) (result []models.Balance, err error) {
	err = client.privateGet(ctx, endpointTradingBalance, nil, &result)
	return
}

func (client Client) getActiveOrders(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointOrder, params, &result)
	return
}

func (client Client) getActiveOrder(ctx context.Context, arguments ...args.Argument) (result *models.Order, err error) {
	params, _ := args.BuildParams(arguments, "clientOrderId")
	err = client.privateGet(ctx, endpointOrder+"/"+params["clientOrderId"].(string), nil, &result)
	return
}

func (client Client) createOrder(ctx context.Context, arguments ...args.Argument) (result *models.Order, err error) {
	// params, err := args.BuildParams(arguments, "symbol", "side", "quantity")
	params, _ := args.BuildParams(arguments)
	// if err != nil {
	// 	return
	// }
	if clientOrderID, ok := params["clientOrderId"]; ok {
		err = client.put(ctx, endpointOrder+"/"+clientOrderID.(string), params, &result)
	} else {
		err = client.post(ctx, endpointOrder, params, &result)
	}
	return
}

func (client Client) cancelAllOrders(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.delete(ctx, endpointOrder, params, &result)
	return
}

func (client Client) cancelOrder(ctx context.Context, arguments ...args.Argument) (result *models.Order, err error) {
	params, err := args.BuildParams(arguments, "clientOrderId")
	if err != nil {
		return
	}
	err = client.delete(ctx, endpointOrder+"/"+params["clientOrderId"].(string), nil, &result)
	return
}

func (client Client) getTradingFee(ctx context.Context, arguments ...args.Argument) (result *models.TradingFee, err error) {
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

func (client Client) getOrderHistory(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointOrderHistory, params, &result)
	return
}

func (client Client) getOrders(ctx context.Context, arguments ...args.Argument) (result []models.Order, err error) {
	params, err := args.BuildParams(arguments, "clientOrderId")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointOrderHistory, params, &result)
	return
}

func (client Client) getTradeHistory(ctx context.Context, arguments ...args.Argument) (result []models.Trade, err error) {
	params, _ := args.BuildParams(arguments)
	err = client.privateGet(ctx, endpointTradeHistory, params, &result)
	return
}

func (client Client) getTradesByOrder(ctx context.Context, arguments ...args.Argument) (result []models.Trade, err error) {
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

func (client Client) getAccountBalance(ctx context.Context, arguments ...args.Argument) (result []models.Balance, err error) {
	err = client.privateGet(ctx, endpointAccountBalance, nil, &result)
	return
}

func (client Client) getDepositCryptoAddress(ctx context.Context, arguments ...args.Argument) (result *models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointCryptoAdress+"/"+params["currency"].(string), nil, &result)
	return
}

func (client Client) createDepositCryptoAddress(ctx context.Context, arguments ...args.Argument) (result *models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointCryptoAdress+"/"+params["currency"].(string), nil, &result)
	return
}

func (client Client) getLast10DepositCryptoAddresses(ctx context.Context, arguments ...args.Argument) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointCryptoAdresses+"/"+params["currency"].(string), nil, &result)
	return
}

func (client Client) getLast10UsedCryptoAddresses(ctx context.Context, arguments ...args.Argument) (result []models.CryptoAddress, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointUsedAddressed+"/"+params["currency"].(string), nil, &result)
	return
}

// not working
func (client Client) withdrawCrypto(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "currency", "amount", "address")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointWithdrawCrypto, params, &result)
	return
}

// not working
func (client Client) transferConvert(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "fromCurrency", "toCurrency", "amount")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointTransferConvert, params, &result)
	return
}


// not working
func (client Client) commitWithdrawCrypto(ctx context.Context, arguments ...args.Argument) (result bool, err error) {
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

// not working
func (client Client) rollbackWithdrawCrypto(ctx context.Context, arguments ...args.Argument) (result bool, err error) {
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

func (client Client) getEstimatesWithdrawFee(ctx context.Context, arguments ...args.Argument) (result string, err error) {
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

func (client Client) checkIfCryptoAddressIsMine(ctx context.Context, arguments ...args.Argument) (result bool, err error) {
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

func (client Client) transferFromTradingToAccountBalance(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	arguments = append(arguments, args.TransferType(transferTypeExchangeToBank))
	params, err := args.BuildParams(arguments, "currency", "amount")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointAccountTranser, params, &result)
	return
}

func (client Client) transferFromAccountToTradingBalance(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	arguments = append(arguments, args.TransferType(transferTypeBankToExchange))
	params, err := args.BuildParams(arguments, "currency", "amount")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointAccountTranser, params, &result)
	return
}

func (client Client) transferMoneyToAnotherUser(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "currency", "amount", "by", "identifier")
	if err != nil {
		return
	}
	err = client.post(ctx, endpointAccountTranserInternal, params, &result)
	return
}

func (client Client) getTransactionHistory(ctx context.Context, arguments ...args.Argument) (result []models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "currency")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTransactionHistory, params, &result)
	return
}

func (client Client) getTransaction(ctx context.Context, arguments ...args.Argument) (result *models.Transaction, err error) {
	params, err := args.BuildParams(arguments, "id")
	if err != nil {
		return
	}
	err = client.privateGet(ctx, endpointTransactionHistory+"/"+params["id"].(string), nil, &result)
	return
}