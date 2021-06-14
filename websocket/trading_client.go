package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

// TradingClient connects via websocket to cryptomarket to enable the user to manage orders. uses SHA256 as auth method and authenticates automatically.
type TradingClient struct {
	clientBase
}

// NewTradingClient returns a new chan client if the connection with the
// cryptomarket server is successful and if the authentication is successfull.
// return error otherwise.
func NewTradingClient(apiKey, apiSecret string) (*TradingClient, error) {
	methodMapping := map[string]string{
		// reports
		"subscribeReports":   "reports",
		"unsubscribeReports": "reports",
		"activeOrders":       "reports",
		"report":             "reports",
	}
	client := &TradingClient{
		clientBase: clientBase{
			wsManager: newWSManager("/api/2/ws/trading"),
			chanCache: newChanCache(),
			subscriptionKeysFunc: func(method string, params map[string]interface{}) (string, bool) {
				val, ok := methodMapping[method]
				return val, ok
			},
			keyFromResponse: func(response wsResponse) string {
				return methodMapping[response.Method]
			},
		},
	}

	// connect to streaming
	if err := client.wsManager.connect(); err != nil {
		return nil, fmt.Errorf("Error in websocket client connection: %s", err)
	}
	// handle incomming data
	go client.handle(client.wsManager.rcv)

	if err := client.authenticate(apiKey, apiSecret); err != nil {
		return nil, err
	}
	return client, nil
}

// GetTradingBalance Get the user trading balance.
//
// https://api.exchange.cryptomarket.com/#get-trading-balance
//
func (client *TradingClient) GetTradingBalance(ctx context.Context) ([]models.Balance, error) {
	var resp struct {
		Result []models.Balance
	}
	err := client.doRequest(ctx, methodGetTradingBalance, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// GetActiveOrders gets the account active orders.
//
// https://api.exchange.cryptomarket.com/#get-active-orders-2
func (client *TradingClient) GetActiveOrders(ctx context.Context) ([]models.Report, error) {
	var resp struct {
		Result []models.Report
	}
	err := client.doRequest(ctx, methodGetOrders, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// CreateOrder creates a new order.
//
// https://api.exchange.cryptomarket.com/#place-new-order
//
// Arguments:
//  ClientOrderID(string)        // Must be unique within the trading day, including all active orders.
//  Symbol(string)               // Trading symbol
//  Side(SideType)               // SideTypeBuy or SideTypeSell
//  Quantity(string)             // Order quantity
//  Type(OrderType)              // Optional. Default is OrderTypeLimit
//  TimeInForce(TimeInForceType) // Optional. Default is TimeInForceTypeGTC
//  Price(string)                // Required for OrderTypelimit and OrderTypeStopLimit. limit price of the order
//  StopPrice(string)            // Required for OrderTypeStopLimit and OrderTypeStopMarket orders. stop price of the order
//  ExpireTime(string)           // Required for orders with TimeInForceTypeGDT
//  StrictValidate(bool)         // Optional. If False, the server rounds half down for tickerSize and quantityIncrement. Example of ETHBTC: tickSize = '0.000001', then price '0.046016' is valid, '0.0460165' is invalid
//  PostOnly(bool)               // Optional. If True, your post_only order causes a match with a pre-existing order as a taker, then the order will be cancelled
func (client *TradingClient) CreateOrder(ctx context.Context, arguments ...args.Argument) (*models.Report, error) {
	var resp struct {
		Result models.Report
	}
	err := client.doRequest(ctx, methodNewOrder, arguments, []string{"clientOrderId", "symbol", "side", "quantity"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// CancelOrder cancels the order with ClientOrderID
//
// Requires authentication
//
// https://api.exchange.cryptomarket.com/#cancel-order
//
// Arguments:
//  ClientOrderID(string) // The ClientOrderId of the order to cancel
func (client *TradingClient) CancelOrder(ctx context.Context, arguments ...args.Argument) (*models.Report, error) {
	var resp struct {
		Result models.Report
	}
	err := client.doRequest(ctx, methodCancelOrder, arguments, []string{"clientOrderId"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// ReplaceOrder rewrites an order, canceling it or replacing it.
//
// The Cancel/Replace request is used to change the parameters of an existing order and to change the quantity or price attribute of an open order.
// Do not use this request to cancel the quantity remaining in an outstanding order. Use the cancelOrder(...) for this purpose.
// It is stipulated that a newly entered order cancels a prior order that has been entered, but not yet executed.
//
// https://api.exchange.cryptomarket.com/#cancel-replace-order
//
// Arguments:
//  ClientOrderID(string)        // The ClientOrderId of the order to modify
//  RequestClientOrderID(string) // The new ClientOrderId of the order to modify
//  Price(string)                // the new price of the order
//  Quantity(string)             // the new quantity of the order
//  StrictValidate(bool)         // Optional. If False, the server rounds half down for tickerSize and quantityIncrement. Example of ETHBTC: tickSize = '0.000001', then price '0.046016' is valid, '0.0460165' is invalid
func (client *TradingClient) ReplaceOrder(ctx context.Context, arguments ...args.Argument) (*models.Report, error) {
	var resp struct {
		Result models.Report
	}
	err := client.doRequest(ctx, methodReplaceOrder, arguments, []string{"clientOrderId", "requestClientId", "price", "quantity"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

///////////////////
// Subscriptions //
///////////////////

// SubscribeToReports subscribes to a feed of trading events of the account.
//
// https://api.exchange.cryptomarket.com/#subscribe-to-reports
func (client *TradingClient) SubscribeToReports() (feedCh chan models.Report, err error) {
	dataCh, err := client.doSubscription(methodSubscribeReports, nil, nil)
	if err != nil {
		return nil, err
	}
	feedCh = make(chan models.Report)
	go func() {
		defer close(feedCh)
		// the first time it recieves a list of reports
		var reports struct {
			Params []models.Report
		}
		data := <-dataCh
		json.Unmarshal(data, &reports)
		for _, report := range reports.Params {
			feedCh <- report
		}
		// then recieves one report at a time
		var resp struct {
			Params models.Report
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			feedCh <- resp.Params
		}
	}()
	return feedCh, nil
}
