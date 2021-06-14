package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cryptomarket/cryptomarket-go/args"

	"github.com/cryptomarket/cryptomarket-go/models"
)

// AccountClient connects via websocket to cryptomarket to get account information of the user. uses SHA256 as auth method and authenticates automatically.
type AccountClient struct {
	clientBase
}

// NewAccountClient returns a new chan client if the connection with the
// cryptomarket server is successful and if the authentication is successful.
// return error otherwise.
func NewAccountClient(apiKey, apiSecret string) (*AccountClient, error) {
	methodMapping := map[string]string{
		// transaction
		"unsubscribeTransactions": "transaction",
		"subscribeTransactions":   "transaction",
		"updateTransaction":       "transaction",
		// balance
		"unsubscribeBalance": "balance",
		"subscribeBalance":   "balance",
		"balance":            "balance",
	}
	client := &AccountClient{
		clientBase: clientBase{
			wsManager: newWSManager("/api/2/ws/account"),
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

// GetAccountBalance gets the account balance
//
// https://api.exchange.cryptomkt.com/#request-balance
func (client *AccountClient) GetAccountBalance(ctx context.Context) ([]models.Balance, error) {
	var resp struct {
		Result []models.Balance
	}
	err := client.doRequest(ctx, methodGetBalance, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// FindTransactions gets a list of transactions of the account. Accepts only filtering by Datetime
//
// https://api.exchange.cryptomkt.com/#find-transactions
//
// Arguments:
//  Currency(string)  // Optional. Currency to filter transactions by.
//  Sort(SortType)    // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeDESC
//  From(string)      // Optional. Initial value of the queried interval
//  Till(string)      // Optional. Last value of the queried interval
//  Limit(int)        // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)       // Optional. Default is 0. Max is 100000
//  ShowSenders(bool) // Optional. If true incluedes senders addresses. Default is false.
func (client *AccountClient) FindTransactions(ctx context.Context, arguments ...args.Argument) ([]models.Transaction, error) {
	temp := make(map[string]interface{})
	for idx, argument := range arguments {
		argument(temp)
		if val, ok := temp["sort"]; ok {
			arguments[idx] = func(params map[string]interface{}) {
				params["order"] = val
			}
			break
		}
	}
	var resp struct {
		Result []models.Transaction
	}
	err := client.doRequest(ctx, methodFindTransactions, arguments, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// LoadTransactions gets a list of transactions of the account. Accepts only filtering by Index
//
// https://api.exchange.cryptomkt.com/#load-transactions
//
// Arguments:
//  Currency(string)  // Optional. Currency to filter transactions by.
//  Sort(SortType)    // Optional. Sort direction. SortTypeASC or SortTypeDESC. Default is SortTypeASC
//  From(string)      // Optional. Initial value of the queried interval (Included)
//  Till(string)      // Optional. Last value of the queried interval (Excluded)
//  Limit(int)        // Optional. Trades per query. Defaul is 100. Max is 1000
//  Offset(int)       // Optional. Default is 0. Max is 100000
//  ShowSenders(bool) // Optional. If true incluedes senders addresses. Default is false.
func (client *AccountClient) LoadTransactions(ctx context.Context, arguments ...args.Argument) ([]models.Transaction, error) {

	temp := make(map[string]interface{})
	for idx, argument := range arguments {
		argument(temp)
		if val, ok := temp["sort"]; ok {
			arguments[idx] = func(params map[string]interface{}) {
				params["order"] = val
			}
			break
		}
	}
	var resp struct {
		Result []models.Transaction
	}
	err := client.doRequest(ctx, methodLoadTransactions, arguments, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

///////////////////
// Subscriptions //
///////////////////

// SubscribeToTransactions subscribes to a feed of transactions of the account.
//
// A transaction notification occurs each time the transaction has been changed:
// such as creating a transaction, updating the pending state (for example the hash assigned)
//or completing a transaction. This is the easiest way to track deposits or develop real-time asset monitoring.
//
// A combination of the recovery mechanism and transaction subscription provides reliable and consistent information
// regarding transactions. For that, you should store the latest processed index and
// requested possible gap using a "loadTransactions" method after connecting or reconnecting the Websocket.
//
// https://api.exchange.cryptomkt.com/#subscription-to-the-transactions
func (client *AccountClient) SubscribeToTransactions() (feedCh chan models.Transaction, err error) {
	dataCh, err := client.doSubscription("subscribeTransactions", nil, nil)
	if err != nil {
		return nil, err
	}
	feedCh = make(chan models.Transaction)
	go func() {
		defer close(feedCh)
		var resp struct {
			Params models.Transaction
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			feedCh <- resp.Params
		}
	}()
	return feedCh, nil
}

// UnsubscribeToTransactions unsubscribe to the transaction feed.
//
// It also closes the feedCh of the subscription
//
// https://api.exchange.cryptomkt.com/#subscription-to-the-transactions
func (client *AccountClient) UnsubscribeToTransactions() error {
	return client.doUnsubscription("unsubscribeTransactions", nil, nil)
}

// SubscribeToTransactions subscribes to a feed of transactions of the account.
//
// A transaction notification occurs each time the transaction has been changed:
// such as creating a transaction, updating the pending state (for example the hash assigned)
//or completing a transaction. This is the easiest way to track deposits or develop real-time asset monitoring.
//
// A combination of the recovery mechanism and transaction subscription provides reliable and consistent information
// regarding transactions. For that, you should store the latest processed index and
// requested possible gap using a "loadTransactions" method after connecting or reconnecting the Websocket.
//
// https://api.exchange.cryptomkt.com/#subscription-to-the-transactions
func (client *AccountClient) SubscribeToBalance() (feedCh chan []models.Balance, err error) {
	dataCh, err := client.doSubscription("subscribeBalance", nil, nil)
	if err != nil {
		return nil, err
	}
	feedCh = make(chan []models.Balance)
	go func() {
		defer close(feedCh)
		var resp struct {
			Params []models.Balance
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			feedCh <- resp.Params
		}
	}()
	return feedCh, nil
}

// UnsubscribeToTransactions unsubscribe to the transaction feed.
//
// It also closes the feedCh of the subscription
//
// https://api.exchange.cryptomkt.com/#subscription-to-the-transactions
func (client *AccountClient) UnsubscribeToBalance() error {
	return client.doUnsubscription("unsubscribeBalance", nil, nil)
}
