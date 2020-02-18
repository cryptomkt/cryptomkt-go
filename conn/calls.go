package conn

import (
	"encoding/json"
	"fmt"

	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/requests"
)

func makeReq(required []string, args ...args.Argument) (*requests.Request, error) {
	req := requests.NewReq(required)
	for _, argument := range args {
		err := argument(req)
		if err != nil {
			return nil, fmt.Errorf("argument error: %s", err)
		}
	}
	err := req.AssertRequired()
	if err != nil {
		return nil, fmt.Errorf("required arguments not meeted:%s", err)
	}
	return req, nil
}

func (client *Client) postReq(endpoint string, caller string, required []string, args ...args.Argument) (string, error) {
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", caller, err)
	}
	return client.post(endpoint, req)
}

func (client *Client) getReq(endpoint string, caller string, required []string, args ...args.Argument) (string, error) {
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", caller, err)
	}
	return client.get(endpoint, req)
}

// Account gives the account information of the client.
// https://developers.cryptomkt.com/es/#cuenta
func (client *Client) GetAccount() (*Account, error) {
	resp, err := client.get("account", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var accountResp AccountResponse
	json.Unmarshal([]byte(resp), &accountResp)
	if accountResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", accountResp.Message)
	}
	return &accountResp.Data, nil
}

// Balance returns the actual balance of the wallets of the client in Cryptomarket
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) GetBalance() (*[]Balance, error) {
	resp, err := client.get("balance", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var balanceResp BalancesResponse
	json.Unmarshal([]byte(resp), &balanceResp)
	if balanceResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", balanceResp.Message)
	}
	return &balanceResp.Data, nil
}

// Wallets is an alias for Balance
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) GetWallets() (*[]Balance, error) {
	return client.GetBalance()
}

// Transactions returns the movements of the wallets of the client.
//
// List of accepted Arguments:
//   - required: Currency
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#obtener-movimientos
func (client *Client) GetTransactions(args ...args.Argument) (*[]Transaction, error) {
	resp, err := client.getReq("transactions", "GetTransaction", []string{"currency"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var transactionsResp TransactionsResponse
	json.Unmarshal([]byte(resp), &transactionsResp)
	if transactionsResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", transactionsResp.Message)
	}
	return &transactionsResp.Data, nil
}

// ActiveOrders returns the list of active orders of the client
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-activas
func (client *Client) GetActiveOrders(args ...args.Argument) (*[]Order, error) {
	resp, err := client.getReq("orders/active", "GetActiveOrders", []string{"market"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var activeOrdersResp OrderListResp
	json.Unmarshal([]byte(resp), &activeOrdersResp)
	if activeOrdersResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", activeOrdersResp.Message)
	}
	return &activeOrdersResp.Data, nil
}

// ExecutedOrders return a list of the executed orders of the client
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-ejecutadas
func (client *Client) GetExecutedOrders(args ...args.Argument) (*[]Order, error) {
	resp, err := client.getReq("orders/executed", "GetExecutedOrders", []string{"market"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var activeOrdersResp OrderListResp
	json.Unmarshal([]byte(resp), &activeOrdersResp)
	if activeOrdersResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", activeOrdersResp.Message)
	}
	return &activeOrdersResp.Data, nil
}

// OrderStatus gives the status of an order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden
func (client *Client) GetOrderStatus(args ...args.Argument) (*[]Order, error) {
	resp, err := client.getReq("orders/status", "GetOrderStatus", []string{"id"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var activeOrdersResp OrderListResp
	json.Unmarshal([]byte(resp), &activeOrdersResp)
	if activeOrdersResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", activeOrdersResp.Message)
	}
	return &activeOrdersResp.Data, nil
}

// Instant emulates an order in the current state of the Instant Exchange of CryptoMarket
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#obtener-cantidad
func (client *Client) GetInstant(args ...args.Argument) (*Quantity, error) {
	resp, err := client.getReq("orders/instant/get", "GetInstant", []string{"market", "type", "amount"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var instantResp InstantResponse
	json.Unmarshal([]byte(resp), &instantResp)
	if instantResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", instantResp.Message)
	}
	return &instantResp.Data, nil
}

// CreateOrder signal to create an order of buy or sell in CryptoMarket
//
// List of accepted Arguments:
//   - required: Amount, Market, Price, Type
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden
func (client *Client) CreateOrder(args ...args.Argument) (*Order, error) {
	resp, err := client.postReq("orders/create", "CreateOrder", []string{"amount", "market", "price", "type"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var orderResp OrderResponse
	json.Unmarshal([]byte(resp), &orderResp)
	if orderResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", orderResp.Message)
	}
	return &orderResp.Data, nil
}

// CancelOrder signal to cancel an order in CryptoMarket
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#cancelar-una-orden
func (client *Client) CancelOrder(args ...args.Argument) (*Order, error) {
	resp, err := client.postReq("orders/cancel", "CancelOrder", []string{"id"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var orderResp OrderResponse
	json.Unmarshal([]byte(resp), &orderResp)
	if orderResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", orderResp.Message)
	}
	return &orderResp.Data, nil
}

// CreateInstant makes an order in the Instant Exchange of CryptoMarket
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden-2
func (client *Client) CreateInstant(args ...args.Argument) error {
	resp, err := client.postReq("orders/instant/create", "CreateInstant", []string{"market", "type", "amount"}, args...)
	if err != nil {
		return fmt.Errorf("error while making the request: %s", err)
	}
	var instantResp InstantResponse
	json.Unmarshal([]byte(resp), &instantResp)
	if instantResp.Status == "error" {
		return fmt.Errorf("error in the response: %s", instantResp.Message)
	}
	return nil
}

// RequestDeposit notifies a deposit to a wallet of local currency
//
// List of accepted Arguments:
//   - required: Amount, BankAccount
// -only for México, Brasil and European Union: Voucher
// -only for México: Date, TrackingCode
// https://developers.cryptomkt.com/es/#notificar-deposito
func (client *Client) RequestDeposit(args ...args.Argument) error {
	resp, err := client.postReq("request/deposit", "RequestDeposit", []string{"amount", "bank_account"}, args...)
	if err != nil {
		return fmt.Errorf("error while making the request: %s", err)
	}
	var instantResp InstantResponse
	json.Unmarshal([]byte(resp), &instantResp)
	if instantResp.Status == "error" {
		return fmt.Errorf("error in the response: %s", instantResp.Message)
	}
	return nil
}

// Request withdrawal notifies a withdrawal from a bank account of the client
//
// List of accepted Arguments:
//   - required: Amount, BankAccount
//   - optional: none
// https://developers.cryptomkt.com/es/#notificar-retiro
func (client *Client) RequestWithdrawal(args ...args.Argument) error {
	resp, err := client.postReq("request/withdrawal", "RequestWithdrawal", []string{"amount", "bank_account"}, args...)
	if err != nil {
		return fmt.Errorf("error while making the request: %s", err)
	}
	var instantResp InstantResponse
	json.Unmarshal([]byte(resp), &instantResp)
	if instantResp.Status == "error" {
		return fmt.Errorf("error in the response: %s", instantResp.Message)
	}
	return nil
}

// Transfer move crypto between wallets
//
// List of accepted Arguments:
//   - required: Address, Amount, Currency
//   - optional: Memo
// https://developers.cryptomkt.com/es/#transferir
func (client *Client) Transfer(args ...args.Argument) error {
	resp, err := client.postReq("transfer", "Transfer", []string{"address", "amount", "currency"}, args...)
	if err != nil {
		return fmt.Errorf("error while making the request: %s", err)
	}
	var instantResp InstantResponse
	json.Unmarshal([]byte(resp), &instantResp)
	if instantResp.Status == "error" {
		return fmt.Errorf("error in the response: %s", instantResp.Message)
	}
	return nil

}

// NewOrder enables a payment order, and gives a QR and urls
//
// List of accepted Arguments:
//   - required: ToReceive, ToReceiveCurrency, PaymentReceiver
//   - optional: ExternalId, CallbackUrl, ErrorUrl, SuccessUrl, RefundEmail, Language
// https://developers.cryptomkt.com/es/#crear-orden-de-pago
func (client *Client) NewOrder(args ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/new_order", "NewOrder", []string{"to_receive", "to_receive_currency", "payment_receiver"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var paymentResp PaymentResponse
	json.Unmarshal([]byte(resp), &paymentResp)
	if paymentResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", paymentResp.Message)
	}
	return &paymentResp.Data, nil
}

// CreateWallet creates a wallet to pay a payment order
//
// List of accepted Arguments:
//   - required: Id, Token, Wallet
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-billetera-de-orden-de-pago
func (client *Client) CreateWallet(args ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/create_wallet", "CreateWallet", []string{"id", "token", "wallet"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var paymentResp PaymentResponse
	json.Unmarshal([]byte(resp), &paymentResp)
	if paymentResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", paymentResp.Message)
	}
	return &paymentResp.Data, nil
}

// PaymentOrders returns all the generated payment orders
//
// List of accepted Arguments:
//   - required: StartDate, EndDate
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#listado-de-ordenes-de-pago
func (client *Client) PaymentOrders(args ...args.Argument) (*[]PaymentOrder, error) {
	resp, err := client.postReq("payment/orders", "PaymentOrders", []string{"start_date", "end_date"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var paymentOrdsResp PaymentOrdersResponse
	json.Unmarshal([]byte(resp), &paymentOrdsResp)
	if paymentOrdsResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", paymentOrdsResp.Message)
	}
	return &paymentOrdsResp.Data, nil
}

// PaymentStatus gives the status of a pyment order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden-de-pago
func (client *Client) PaymentStatus(args ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/status", "PaymentStatus", []string{"id"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var paymentResp PaymentResponse
	json.Unmarshal([]byte(resp), &paymentResp)
	if paymentResp.Status == "error" {
		return nil, fmt.Errorf("error in the response: %s", paymentResp.Message)
	}
	return &paymentResp.Data, nil
}

func (client *Client) MarketList(args ...args.Argument) (*MarketStruct, error) {
	required := []string{"market"}
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in MakeMarket: %s", err)
	}
	resp, err := client.getPublic("market", req)
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}
	// estructuar el output
	var response map[string]interface{}
	var respu MarketStruct

	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return nil, err
	} else if response["status"].(string) == "success" {
		var i int = 0
		var largo int = len(response["data"].([]interface{}))
		respu.Data = make([]string, largo)
		for i < largo {
			respu.Data[i] = response["data"].([]interface{})[i].(string)
			i += 1
		}
		return &respu, nil
	} else {
		panic("Response from server failed")
	}
}

func makeArrayMap(respString string, response map[string]interface{}, data []map[string]string) ([]map[string]string, error) {
	if err := json.Unmarshal([]byte(respString), &response); err != nil {
		return data, err
	} else if response["status"].(string) == "success" {
		var i int = 0
		var largo int = len(response["data"].([]interface{}))
		paraConvertir := make([]interface{}, largo)
		paraConvertir = response["data"].([]interface{})
		data = make([]map[string]string, len(paraConvertir))
		for i < largo {
			data[i] = make(map[string]string)
			for key, value := range paraConvertir[i].(map[string]interface{}) {
				data[i][key] = value.(string)
			}
			i += 1
		}
		return data, nil
	} else {
		return data, fmt.Errorf("Response from server failed")
	}
}

//generalizar makeArrayMap
func (client *Client) MakeTicker(args ...args.Argument) (*Ticker, error) {
	resp, err := client.getPublic("ticker", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}

	var response map[string]interface{}
	var respu Ticker
	data, err := makeArrayMap(resp, response, respu.Data)
	if err == nil {
		respu.Data = data
		return &respu, err
	} else {
		return nil, fmt.Errorf("Cannot create Ticker object")
	}
}

func (client *Client) getBook(args ...args.Argument) (*Book, error) {
	required := []string{"market", "type"}
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in MakeMarket: %s", err)
	}
	resp, err := client.getPublic("book", req)
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}

	var response map[string]interface{}
	var respu Book
	data, err := makeArrayMap(resp, response, respu.Data)
	if err == nil {
		respu.Data = data
		return &respu, nil
	} else {
		return nil, fmt.Errorf("Cannot create Order object")
	}
}

func (client *Client) getTrades(args ...args.Argument) (*Trades, error) {
	required := []string{"market"}
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in MakeMarket: %s", err)
	}
	resp, err := client.getPublic("trades", req)
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}

	var response map[string]interface{}
	var respu Trades
	data, err := makeArrayMap(resp, response, respu.Data)
	if err == nil {
		respu.Data = data
		return &respu, nil
	} else {
		return nil, fmt.Errorf("Cannot create Trades object")
	}
}

func (client *Client) getPrices(args ...args.Argument) (*Prices, error) {
	required := []string{"market", "timeframe"}
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in MakeMarket: %s", err)
	}
	resp, err := client.getPublic("prices", req)
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}

	var response map[string]interface{}
	var respu Prices
	//unmarshal string to response
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return nil, err
	}

	//data marshaled to []byte
	if dataMarsh, err := json.Marshal(response["data"]); err != nil {
		return nil, err
	} else {
		//data unmershaled to struct Prices: map[string][]Field where Field is a struct with
		//its first field as integer and the others as string. If you wanna get the Candle_id from
		//ask label from Prices struct you must call as (objectPrice).Data["ask"].Candle_id.
		//this is because golang doesnt support dict with a value as interger and the rest string.
		err := json.Unmarshal(dataMarsh, &respu)
		if err != nil {
			return nil, err
		}
	}
	return &respu, nil
}
