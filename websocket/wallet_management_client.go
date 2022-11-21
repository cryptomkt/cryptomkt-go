package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/internal"
	"github.com/cryptomarket/cryptomarket-go/models"
)

// WalletManagementClient connects via websocket to cryptomarket to get wallet information of the user. uses SHA256 as auth method and authenticates automatically.

type WalletManagementClient struct {
	clientBase
}

// NewWalletManagementClient returns a new wallet client if the connection with the
// cryptomarket server is successful and if the authentication is successful.
// return error otherwise.
// Arguments:
//  apiKey // The API key
//  apiSecret // The API secret
//  window // Maximum difference between the creation of the request and the moment of request processing in milliseconds. Max is 60_000. Defaul is 10_000 (use 0 as argument for default)
func NewWalletManagementClient(apiKey, apiSecret string, window int) (*WalletManagementClient, error) {
	client := &WalletManagementClient{
		clientBase: clientBase{
			wsManager: newWSManager("/api/3/ws/wallet"),
			chanCache: newChanCache(),
			window:    window,
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

// GetWalletBalance gets the user's wallet balance for all currencies with balance
//
// https://api.exchange.cryptomkt.com/#request-wallet-balance
func (client *WalletManagementClient) GetWalletBalances(
	ctx context.Context,
) ([]models.Balance, error) {
	var resp struct {
		Result []models.Balance
	}
	err := client.doRequest(ctx, methodWalletBalances, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// GetWalletBalanceOfCurrency gets the user's wallet balance of a currency
//
// https://api.exchange.cryptomkt.com/#request-wallet-balance
//
// Arguments:
//  Currency(string)  // The currency code to query the balance
func (client *WalletManagementClient) GetWalletBalanceOfCurrency(
	ctx context.Context,
	arguments ...args.Argument,
) (*models.Balance, error) {
	var resp struct {
		Result models.Balance
	}
	err := client.doRequest(
		ctx,
		methodWalletBalance,
		arguments,
		[]string{internal.ArgNameCurrency},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// GetTransactions gets the transactions of the account
//
// Important:
//
//  - The list of supported transaction types may be expanded in future versions
//
//  - Some transaction subtypes are reserved for future use and do not purport to provide any functionality on the platform
//
//  - The list of supported transaction subtypes may be expanded in future versions
//
// https://api.exchange.cryptomkt.com/#get-transactions
//
// Arguments:
//  TransactionIds([]string)  // Optional. List of transaction identifiers to query
//  TransactionTypes([]TransactionType)  // Optional. List of types to query. valid types are: TransactionDeposit, TransactionWithdraw, TransactionTransfer and TransactionSwap
//  TransactionSubTypes([]TransactionSubType)  // Optional. List of subtypes to query. valid subtypes are: TransactionSubTypeUnclassified, TransactionSubTypeBlockchain,  TransactionSubTypeAffiliate,  TransactionSubtypeOffchain, TransactionSubTypeFiat, TransactionSubTypeSubAccount, TransactionSubTypeWalletToSpot, TransactionSubTypeSpotToWallet, TransactionSubTypeChainSwitchFrom and TransactionSubTypeChainSwitchTo
//  TransactionStatuses([]TransactionStatusType)  // Optional. List of statuses to query. valid subtypes are: TransactionStatusCreated, TransactionStatusPending, TransactionStatusFailed, TransactionStatusSuccess and TransactionStatusRolledBack
//  SortBy(SortByType)  // Optional. sorting parameter. SortByCreatedAt or SortByID. Default is SortByCreatedAt
//  From(string)  // Optional. Interval initial value when ordering by CreatedAt. As Datetime.
//  Till(string)  // Optional. Interval end value when ordering by CreatedAt. As Datetime.
//  IDFrom(string)  // Optional. Interval initial value when ordering by id. Min is 0
//  IDTill(string)  // Optional. Interval end value when ordering by id. Min is 0
//  Sort(SortType)  // Optional. Sort direction. SortASC or SortDESC. Default is SortDESC
//  Limit(int64)  // Optional. Transactions per query. Defaul is 100. Max is 1000
//  Offset(int64)  // Optional. Default is 0. Max is 100000
func (client *WalletManagementClient) GetTransactions(
	ctx context.Context,
	arguments ...args.Argument,
) ([]models.Transaction, error) {
	var resp struct {
		Result []models.Transaction
	}
	err := client.doRequest(
		ctx,
		methodGetTransactions,
		arguments,
		nil,
		&resp,
	)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

///////////////////
// Subscriptions //
///////////////////

// SubscribeToTransactions A transaction notification occurs each time a transaction has been changed, such as creating a transaction, updating the pending state (e.g., the hash assigned) or completing a transaction
//
// this subscriptions only sends updates, no snapshot, with a transaction
//
// https://api.exchange.cryptomkt.com/#subscribe-to-transactions
func (client *WalletManagementClient) SubscribeToTransactions() (notificationCh chan models.Notification[models.Transaction], err error) {
	dataCh, err := client.doSubscription(methodSubscribeTransactions, nil, nil)
	if err != nil {
		return nil, err
	}
	notificationCh = make(chan models.Notification[models.Transaction])
	go func() {
		defer close(notificationCh)
		var resp struct {
			Params models.Transaction
		}
		for data := range dataCh {
			json.Unmarshal(data, &resp)
			notificationCh <- models.Notification[models.Transaction]{Data: resp.Params, NotificationType: args.NotificationUpdate}
		}
	}()
	return notificationCh, nil
}

// UnsubscribeToTransactions stop recieving the feed of transactions changes
//
// https://api.exchange.cryptomkt.com/#subscribe-to-transactions
//
func (client *WalletManagementClient) UnsubscribeToTransactions() error {
	return client.doUnsubscription(methodUnsubscribeTransactions, nil, nil)
}

// SubscribeToWalletBalances subscribe to a feed of the user's wallet balances
//
// only non-zero values are present
//
// the subscriptions sends a snapshot for with all balances and one update with one balance at a time
//
// https://api.exchange.cryptomkt.com/#subscribe-to-wallet-balance
func (client *WalletManagementClient) SubscribeToWalletBalance() (notificationCh chan models.Notification[[]models.Balance], err error) {
	dataCh, err := client.doSubscription(methodSubscribeWalletBalances, nil, nil)
	if err != nil {
		return nil, err
	}
	notificationCh = make(chan models.Notification[[]models.Balance])
	go func() {
		defer close(notificationCh)
		sendMap := map[string]func([]byte){
			methodWalletBalances: func(data []byte) {
				var resp struct {
					Params []models.Balance
				}
				json.Unmarshal(data, &resp)
				notificationCh <- models.Notification[[]models.Balance]{
					Data:             resp.Params,
					NotificationType: args.NotificationSnapshot,
				}
			},
			methodWalletBalanceUpdate: func(data []byte) {
				var resp struct {
					Params models.Balance
				}
				json.Unmarshal(data, &resp)
				notificationCh <- models.Notification[[]models.Balance]{
					Data:             []models.Balance{resp.Params},
					NotificationType: args.NotificationUpdate,
				}
			},
		}
		var method struct {
			Method string
		}
		for data := range dataCh {
			json.Unmarshal(data, &method)
			sendMap[method.Method](data)
		}
	}()
	return notificationCh, nil
}

// UnsubscribeToWalletBalances stop recieving the feed of balances changes
//
// https://api.exchange.cryptomkt.com/#subscribe-to-wallet-balance
func (client *WalletManagementClient) UnsubscribeToWalletBalance() error {
	return client.doUnsubscription(methodUnsubscribeWalletBalances, nil, nil)
}
