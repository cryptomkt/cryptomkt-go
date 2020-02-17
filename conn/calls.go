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

// Account gives the account information of the client.
// https://developers.cryptomkt.com/es/#cuenta
func (client *Client) Account() (string, error) {
	return client.get("account", requests.NewEmptyReq())
}

// Balance returns the actual balance of the wallets of the client in Cryptomarket
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) Balance() (string, error) {
	return client.get("balance", requests.NewEmptyReq())
}

// Wallets is an alias for Balance
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) Wallets() (string, error) {
	return client.Balance()
}

// Transactions returns the movements of the wallets of the client.
//
// List of accepted Arguments:
//   - required: Currency
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#obtener-movimientos
func (client *Client) Transactions(args ...args.Argument) (string, error) {
	required := []string{"currency"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in Transaction: %s", err)
	}
	return client.get("transactions", req)
}

// ActiveOrders returns the list of active orders of the client
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-activas
func (client *Client) ActiveOrders(args ...args.Argument) (string, error) {
	required := []string{"market"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in ActiveOrders: %s", err)
	}
	return client.get("orders/active", req)
}

// ExecutedOrders return a list of the executed orders of the client
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-ejecutadas
func (client *Client) ExecutedOrders(args ...args.Argument) (string, error) {
	required := []string{"market"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in ExecutedOrders: %s", err)
	}
	return client.get("orders/executed", req)
}

// OrderStatus gives the status of an order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden
func (client *Client) OrderStatus(args ...args.Argument) (string, error) {
	required := []string{"id"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in OrderStatus: %s", err)
	}
	return client.get("orders/status", req)
}

// Instant emulates an order in the current state of the Instant Exchange of CryptoMarket
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#obtener-cantidad
func (client *Client) Instant(args ...args.Argument) (string, error) {
	required := []string{"market", "type", "amount"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in Instant: %s", err)
	}
	return client.get("orders/instant/get", req)
}

// CreateOrder signal to create an order of buy or sell in CryptoMarket
//
// List of accepted Arguments:
//   - required: Amount, Market, Price, Type
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden
func (client *Client) CreateOrder(args ...args.Argument) (string, error) {
	required := []string{"amount", "market", "price", "type"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in CreateOrder: %s", err)
	}
	return client.post("orders/create", req)
}

// CancelOrder signal to cancel an order in CryptoMarket
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#cancelar-una-orden
func (client *Client) CancelOrder(args ...args.Argument) (string, error) {
	required := []string{"id"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in CancelOrder: %s", err)
	}
	return client.post("orders/cancel", req)
}

// CreateInstant makes an order in the Instant Exchange of CryptoMarket
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden-2
func (client *Client) CreateInstant(args ...args.Argument) (string, error) {
	required := []string{"market", "type", "amount"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in CreateInstant: %s", err)
	}
	return client.post("orders/instant/create", req)
}

// RequestDeposit notifies a deposit to a wallet of local currency
//
// List of accepted Arguments:
//   - required: Amount, BankAccount
// -only for México, Brasil and European Union: Voucher
// -only for México: Date, TrackingCode
// https://developers.cryptomkt.com/es/#notificar-deposito
func (client *Client) RequestDeposit(args ...args.Argument) (string, error) {
	required := []string{"amount", "bank_account"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in RequestDeposit: %s", err)
	}
	return client.post("request/deposit", req)
}

// Request withdrawal notifies a withdrawal from a bank account of the client
//
// List of accepted Arguments:
//   - required: Amount, BankAccount
//   - optional: none
// https://developers.cryptomkt.com/es/#notificar-retiro
func (client *Client) RequestWithdrawal(args ...args.Argument) (string, error) {
	required := []string{"amount", "bank_account"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in RequestWithdrawal: %s", err)
	}
	return client.post("request/withdrawal", req)
}

// Transfer move crypto between wallets
//
// List of accepted Arguments:
//   - required: Address, Amount, Currency
//   - optional: Memo
// https://developers.cryptomkt.com/es/#transferir
func (client *Client) Transfer(args ...args.Argument) (string, error) {
	required := []string{"address", "amount", "currency"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in Transfer: %s", err)
	}
	return client.post("transfer", req)
}

// NewOrder enables a payment order, and gives a QR and urls
//
// List of accepted Arguments:
//   - required: ToReceive, ToReceiveCurrency, PaymentReceiver
//   - optional: ExternalId, CallbackUrl, ErrorUrl, SuccessUrl, RefundEmail, Language
// https://developers.cryptomkt.com/es/#crear-orden-de-pago
func (client *Client) NewOrder(args ...args.Argument) (string, error) {
	required := []string{"to_receive", "to_receive_currency", "payment_receiver"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in NewOrder: %s", err)
	}
	return client.post("payment/new_order", req)
}

// CreateWallet creates a wallet to pay a payment order
//
// List of accepted Arguments:
//   - required: Id, Token, Wallet
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-billetera-de-orden-de-pago
func (client *Client) CreateWallet(args ...args.Argument) (string, error) {
	required := []string{"id", "token", "wallet"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in CreateWallet: %s", err)
	}
	return client.post("payment/create_wallet", req)
}

// PaymentOrders returns all the generated payment orders
//
// List of accepted Arguments:
//   - required: StartDate, EndDate
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#listado-de-ordenes-de-pago
func (client *Client) PaymentOrders(args ...args.Argument) (string, error) {
	required := []string{"start_date", "end_date"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in PaymentOrders: %s", err)
	}
	return client.get("payment/orders", req)
}

// PaymentStatus gives the status of a pyment order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden-de-pago
func (client *Client) PaymentStatus(args ...args.Argument) (string, error) {
	required := []string{"id"}
	req, err := makeReq(required, args...)
	if err != nil {
		return "", fmt.Errorf("Error in PaymentStatus: %s", err)
	}
	return client.get("payment/status", req)
}

// MarketList returns a pointer to a MarketStruct with the field "data" given by the api. The data given is
// an array of strings indicating the markets in cryptomkt. This function returns two values.
// The first is a reference to the struct created and the second is a error message. It returns (nil, error)
// when an error is raised.
// This method does not accept any arguments.
func (client *Client) GetMarkes() (*MarketStruct, error) {
	resp, err := client.get("market", requests.NewEmptyReq())
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

// makeArrayMap is used for transforming the data given by the api in the functions MakeTicker, MakeOrder and MakeTrades.
// The response given by the api is unmarshaled into the same structure: map[string]interface{}. The data field is
// originally unmarshaled into []interface{}, and it must be turned into []map[string]string.
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

// MakeTicker returns a pointer to a Ticker struct with the data given by the api and an error message. It returns (nil,error)
//when an error is raised and (*Ticker, nil) when the operation is successful. The data fields are: high, low, ask, bid,
//last_price, volume, market and  timestamp
//
// List of accepted Arguments:
//
// 		- optional: Market
func (client *Client) GetTicker(args ...args.Argument) (*Ticker, error) {
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

// MakeOrder returns a pointer to a Order struct with the data given by
// the api and an error message. It returns (nil, error) when an error
// is raised and (*Order, nil) when the operation is successful.
// The data fields are: price, amount and timestamp.
//
// List of accepted Arguments:
//
// 		- required: Market , Type
//		- optional: Page, Limit
func (client *Client) GetOrders(args ...args.Argument) (*Order, error) {
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
	var respu Order
	data, err := makeArrayMap(resp, response, respu.Data)
	if err == nil {
		respu.Data = data
		return &respu, nil
	} else {
		return nil, fmt.Errorf("Cannot create Order object")
	}
}

// MakeTrades returns a pointer to a Trades struct with the data given
// by the api and an error message. It returns (nil, error) when an error
// is raised and (*Trades, nil) when the operation is successful.
// The data fields are market_taker, price, amount, tid, timestamp and market.
//
// List of accepted Arguments:
//
//		- required: Market
//		- optional: Start, End, Page, Limit
func (client *Client) GetTrades(args ...args.Argument) (*Trades, error) {
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

// MakePrices return a pointer to a Prices struct with the data given by
// the api and an error message. It returns (nil,error) when an error
// is raised and (*Prices, nil) when the operation is successful.
// The data field is a map[string][]Field, where the Field structure contains all the
// information. To consult these fields you must call *Prices.Data["ask"][index].fieldYouWant or
// *Prices.Data["bid"][index].fieldYouWant
//
// List of accepted Arguments:
//
//		- required: Market, TimeFrame
//		- optional: Page, Limit
func (client *Client) GetPrices(args ...args.Argument) (*Prices, error) {
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
