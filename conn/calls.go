package conn

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/requests"
)

// GetAccount gives the account information of the client.
//
// https://developers.cryptomkt.com/es/#cuenta
func (client *Client) GetAccount() (*Account, error) {
	resp, err := client.get("account", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var aResp AccountResponse
	json.Unmarshal(resp, &aResp)
	if aResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", aResp.Message)
	}
	return &aResp.Data, nil
}

// GetBalance returns the actual balance of the wallets of the client in CryptoMarket
//
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) GetBalance() (*[]Balance, error) {
	resp, err := client.get("balance", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var bResp BalancesResponse
	json.Unmarshal(resp, &bResp)
	if bResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", bResp.Message)
	}
	return &bResp.Data, nil
}

// GetWallets is an alias for Balance, returns the actual balance of wallets of the client in CryptoMarket
//
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) GetWallets() (*[]Balance, error) {
	return client.GetBalance()
}

// GetTransactions returns the movements of the wallets of the client.
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
	var tResp TransactionsResponse
	json.Unmarshal(resp, &tResp)
	if tResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", tResp.Message)
	}
	return &tResp.Data, nil
}

// GetActiveOrders returns the list of active orders of the client
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
	var oListResp OrderListResp
	json.Unmarshal(resp, &oListResp)
	if oListResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", oListResp.Message)
	}
	return &oListResp.Data, nil
}

// GetExecutedOrders return a list of the executed orders of the client
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
	var oListResp OrderListResp
	json.Unmarshal(resp, &oListResp)
	if oListResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", oListResp.Message)
	}
	return &oListResp.Data, nil
}

// GetOrderStatus gives the status of an order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden
func (client *Client) GetOrderStatus(args ...args.Argument) (*Order, error) {
	resp, err := client.getReq("orders/status", "GetOrderStatus", []string{"id"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var oResp OrderResponse
	json.Unmarshal(resp, &oResp)
	if oResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", oResp.Message)
	}
	return &oResp.Data, nil
}

// GetInstant emulates an order in the current state of the Instant Exchange of CryptoMarket
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
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", iResp.Message)
	}
	return &iResp.Data, nil
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
	var oResp OrderResponse
	json.Unmarshal(resp, &oResp)
	if oResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", oResp.Message)
	}
	return &oResp.Data, nil
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
	var oResp OrderResponse
	json.Unmarshal(resp, &oResp)
	if oResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", oResp.Message)
	}
	return &oResp.Data, nil
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
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from in the  server side: %s", iResp.Message)
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
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from in the  server side: %s", iResp.Message)
	}
	return nil
}

// RequestWithdrawal notifies a withdrawal from a bank account of the client
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
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from in the  server side: %s", iResp.Message)
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
	var iResp InstantResponse
	json.Unmarshal(resp, &iResp)
	if iResp.Status == "error" {
		return fmt.Errorf("error from in the  server side: %s", iResp.Message)
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
	var poResp PaymentResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
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
	var poResp PaymentResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
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
	var poResp PaymentOrdersResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
}

// GetPaymentStatus gives the status of a pyment order
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden-de-pago
func (client *Client) GetPaymentStatus(args ...args.Argument) (*PaymentOrder, error) {
	resp, err := client.postReq("payment/status", "PaymentStatus", []string{"id"}, args...)
	if err != nil {
		return nil, fmt.Errorf("error while making the request: %s", err)
	}
	var poResp PaymentResponse
	json.Unmarshal(resp, &poResp)
	if poResp.Status == "error" {
		return nil, fmt.Errorf("error from in the  server side: %s", poResp.Message)
	}
	return &poResp.Data, nil
}

// Public Endpoints:

//

// GetMarket returns a pointer to a MarketStruct with the field "data" given by the api. The data given is
// an array of strings indicating the markets in cryptomkt. This function returns two values.
// The first is a reference to the struct created and the second is a error message. It returns (nil, error)
// when an error is raised.
// This method does not accept any arguments.
func (client *Client) GetMarkets() (*MarketStruct, error) {
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

// GetTicker returns a pointer to a Ticker struct with the data given by the api and an error message. It returns (nil,error)
//when an error is raised and (*Ticker, nil) when the operation is successful. The data fields are: High, Low, Ask, Bid,
//LastPrice, Volume, Market and Timestamp
//
// List of accepted Arguments:
//
// 		- optional: Market
func (client *Client) GetTicker(args ...args.Argument) (*Ticker, error) {
	resp, err := client.getPublic("ticker", requests.NewEmptyReq())
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}

	var tempTicker TemporalTicker

	err = json.Unmarshal([]byte(resp), &tempTicker)
	if err != nil {
		return nil, fmt.Errorf("Unable to make data")
	}

	var ticker Ticker

	if tempTicker.Status == "success" {
		ticker.Data = tempTicker.Data
		return &ticker, nil
	} else {
		return nil, fmt.Errorf("Response from server failed")
	}
}

// GetBook returns a pointer to a Book struct with the data given by
// the api and an error message. It returns (nil, error) when an error
// is raised and (*Book, nil) when the operation is successful.
// The data fields are: price, amount and timestamp.
//
// List of accepted Arguments:
//
// 		- required: Market , Type
//		- optional: Page, Limit
func (client *Client) GetBook(args ...args.Argument) (*Book, error) {
	required := []string{"market", "type"}
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in MakeMarket: %s", err)
	}
	resp, err := client.getPublic("book", req)
	if err != nil {
		return nil, fmt.Errorf("error at client: %s", err)
	}

	var response TemporalBook
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return nil, err
	} else {
		var resp Book

		if response.Status == "success" {
			resp.args = req.GetArguments()
			resp.pagination = response.Pagination
			resp.client = client
			resp.Data = response.Data
			return &resp, nil
		} else {
			return nil, fmt.Errorf("Response from server failed")
		}
	}
}

// Here you have methods to interact with the Book's pagination

// GetPrevious lets you go to the previous page if it exists, returns (*Book, nil) if
// it is successfull and (nil, error) otherwise
func (b *Book) GetPrevious() (*Book, error) {
	if b.pagination.Previous != nil {
		_, okPage := p.args["page"]
		_, okLimit := b.args["limit"]
		limit, _ := strconv.Atoi(b.args["limit"])
		pageToPut := int(b.pagination.Page - 1)
		if okPage && !okLimit {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut))
		} else if !okPage && okLimit {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut), args.Limit(limit))
		} else if okPage && okLimit {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut), args.Limit(limit))
		} else {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut))
		}
	} else {
		return nil, fmt.Errorf("Cannot go to the previous page because it does not exist")
	}
}

// GetNext lets you go to the next page if it exists, returns (*Book, nil) if
// it is successfull and (nil, error) otherwise
func (b *Book) GetNext() (*Book, error) {
	if b.pagination.Next != nil {
		_, okPage := b.args["page"]
		_, okLimit := b.args["limit"]
		limit, _ := strconv.Atoi(b.args["limit"])
		pageToPut := int(b.pagination.Page + 1)
		if okPage && !okLimit {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut))
		} else if !okPage && okLimit {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut), args.Limit(limit))
		} else if okPage && okLimit {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(pageToPut), args.Limit(limit))
		} else {
			return b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]))
		}
	} else {
		return nil, fmt.Errorf("Cannot go to the next page, because it does not exist")
	}
}

// GetPage returns the page you have
func (b *Book) GetPage() int {
	return b.pagination.Page
}

// GetLimit returns the limit you have provided, but if you have not, it provides the default
func (b *Book) GetLimit() int {
	return b.pagination.Limit
}

func (b *Book) GetAllBooks() ([]BookData, error) {
	var resp []BookData = make([]BookData, 0)
	var bookPointer *Book
	var err error
	if _, ok := b.args["limit"]; ok {
		value, err2 := strconv.Atoi(b.args["limit"])
		if err2 == nil {
			bookPointer, err = b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(0), args.Limit(value))
		}
	} else {
		bookPointer, err = b.client.GetBook(args.Market(b.args["market"]), args.Type(b.args["type"]), args.Page(0))
	}

	for bookPointer.pagination.Next != nil {
		if err != nil {
			return resp, fmt.Errorf("Error getting next page, %s", err)
		} else {
			resp = append(resp, bookPointer.Data...)
		}
		bookPointer, err = bookPointer.GetNext()
		time.Sleep(3 * time.Second)
	}
	resp = append(resp, bookPointer.Data...) // at the end of the loop, there is one left to add
	return resp, nil
}

// GetTrades returns a pointer to a Trades struct with the data given
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

	var response TemporalTrades
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return nil, err
	} else {
		var resp Trades
		if response.Status == "success" {
			resp.args = req.GetArguments()
			resp.pagination = response.Pagination
			resp.client = client
			resp.Data = response.Data
			return &resp, nil
		} else {
			return nil, fmt.Errorf("Response from server failed")
		}
	}
}

// getArgsList returns the args list from a *Trades object that has been created. Only
// used internally by this SDK
func getArgsList(t *Trades, tipe string) ([]args.Argument, error) {
	var newArgs []args.Argument = make([]args.Argument, len(t.args))
	var i int = 0
	for k, v := range t.args {
		switch k {
		case "start":
			newArgs[i] = args.Start(v)
		case "end":
			newArgs[i] = args.End(v)
		case "page":
			page, err := strconv.Atoi(v)
			if err == nil {
				if tipe == "p" {
					newArgs[i] = args.Page(int(page - 1))
				} else if tipe == "n" {
					newArgs[i] = args.Page(int(page + 1))
				}
			} else {
				return nil, fmt.Errorf("Cannot convert page to int")
			}
		case "limit":
			limit, err := strconv.Atoi(v)
			if err == nil {
				newArgs[i] = args.Limit(int(limit))
			} else {
				return nil, fmt.Errorf("Cannot convert limit to int")
			}
		case "market":
			newArgs[i] = args.Market(v)
		default:
			return nil, fmt.Errorf("Unknown argument")
		}
		i++
	}
	return newArgs, nil
}

// executeTradeRandomArgs returns the object with the arguments provided by the object
// Used internally by this SDK
func executeTradeRandomArgs(newArgs []args.Argument, t *Trades) (*Trades, error) {
	switch len(newArgs) {
	case 1:
		return t.client.GetTrades(newArgs[0])
	case 2:
		return t.client.GetTrades(newArgs[0], newArgs[1])
	case 3:
		return t.client.GetTrades(newArgs[0], newArgs[1], newArgs[2])
	case 4:
		return t.client.GetTrades(newArgs[0], newArgs[1], newArgs[2], newArgs[3])
	case 5:
		return t.client.GetTrades(newArgs[0], newArgs[1], newArgs[2], newArgs[3], newArgs[4])
	default:
		return nil, fmt.Errorf("Need one to five args")
	}
}

// executePrevNext returns the previous or the next object page depending on the
// string passed. Used internally by this SDK
func executePrevNext(t *Trades, tipe string) (*Trades, error) {
	var newArgs []args.Argument
	var err error

	if t.pagination.Previous != nil {
		newArgs, err = getArgsList(t, tipe)
		if err != nil {
			return nil, fmt.Errorf("Cannot get args, because %s ", err)
		} else {
			//You have the optional args so far.
			return executeTradeRandomArgs(newArgs, t)
		}
		return nil, fmt.Errorf("Cannot go to previous page, because it does not exist")
	}
}

// Here you find methods to interact with the Trades's pagination

// GetPrevious lets you go to the previous page if it exists, returns (*Trades, nil) if
// it is successfull and (nil, error) otherwise
func (t *Trades) GetPrevious() (*Trades, error) {
	return executePrevNext(t, "p")
}

// GetNext lets you go to the next page if it exists, returns (*Trades, nil) if
// it is successfull and (nil, error) otherwise
func (t *Trades) GetNext() (*Trades, error) {
	return executePrevNext(t, "n")
}

// GetPage returns the page you have
func (t *Trades) GetPage() int {
	return t.pagination.Page
}

// GetLimit returns the limit you have provided, but if you have not, it provides the default
func (t *Trades) GetLimit() int {
	return t.pagination.Limit
}

func createMapStringArgument(argus map[string]string) []args.Argument {
	slice := make([]args.Argument, len(argus))
	var i int = 0
	for k, v := range argus {
		switch k {
		case "start":
			slice[i] = args.Start(v)
		case "end":
			slice[i] = args.End(v)
		case "limit":
			value, err := strconv.Atoi(v)
			if err == nil {
				slice[i] = args.Limit(value)
			} else {
				fmt.Errorf("Cannot conver %s to int, %s", v, err)
			}
		case "page": //page is given
			slice[i] = args.Page(0)
		case "timeframe":
			slice[i] = args.Timeframe(v)
		}
		i++
	}
	//check if "page" argument is not given
	if _, ok := argus["page"]; !ok {
		slice[len(slice)] = args.Page(0)
	}
	return slice
}

// Check if the error is nil when is used, because if it has an error, the response is wrong
func (t *Trades) GetAllTrades() ([]TradesData, error) {
	var resp []TradesData
	var newArgs []args.Argument = createMapStringArgument(t.args) // set the page to zero value

	tradesPointer, err := executeTradeRandomArgs(newArgs, t)
	for tradesPointer.pagination.Next != nil {
		if err != nil {
			return resp, fmt.Errorf("%s", err)
		} else {
			resp = append(resp, tradesPointer.Data...)
		}
		tradesPointer, err = tradesPointer.GetNext()
	}

	//append the last element left
	resp = append(resp, tradesPointer.Data...)
	return resp, nil
}

// GetPrices return a pointer to a Prices struct with the data given by
// the api and an error message. It returns (nil,error) when an error
// is raised and (*Prices, nil) when the operation is successful.
// The data field is a map[string][]Field, where the Field structure contains all the
// information. To consult these fields you must call *Prices.Data.Ask[index].fieldYouWant or
// *Prices.Data.Bid[index].fieldYouWant
//
// List of accepted Arguments:
//
//		- required: Market, Timeframe
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

	var response TemporalPrices
	//unmarshal string to response
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return nil, err
	} else {
		if response.Status == "success" {
			var resp Prices
			resp.args = req.GetArguments()
			resp.pagination = response.Pagination
			resp.client = client
			resp.Data = response.Data
			return &resp, nil
		} else {
			return nil, fmt.Errorf("Response from server failed")
		}
	}
}

// Here you have methods to interact with Prices's pagination

// GetPrevious lets you go to the previous page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (p *Prices) GetPrevious() (*Prices, error) {
	if p.pagination.Next != nil {
		_, okPage := p.args["page"]
		_, okLimit := p.args["limit"]
		limit, _ := strconv.Atoi(p.args["limit"])
		pageToPut := int(p.pagination.Page - 1)
		if okPage && !okLimit {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]), args.Page(pageToPut))
		} else if !okPage && okLimit {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]), args.Page(pageToPut), args.Limit(limit))
		} else if okPage && okLimit {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]), args.Page(pageToPut), args.Limit(limit))
		} else {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]))
		}
	} else {
		return nil, fmt.Errorf("Cannot go to the next page, because it does not exist")
	}
}

// GetNext lets you go to the next page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (p *Prices) GetNext() (*Prices, error) {
	if p.pagination.Next != nil {
		_, okPage := p.args["page"]
		_, okLimit := p.args["limit"]
		limit, _ := strconv.Atoi(p.args["limit"])
		pageToPut := int(p.pagination.Page + 1)
		if okPage && !okLimit {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]), args.Page(pageToPut))
		} else if !okPage && okLimit {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]), args.Page(pageToPut), args.Limit(limit))
		} else if okPage && okLimit {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]), args.Page(pageToPut), args.Limit(limit))
		} else {
			return p.client.GetPrices(args.Market(p.args["market"]), args.Timeframe(p.args["timeframe"]))
		}
	} else {
		return nil, fmt.Errorf("Cannot go to the next page, because it does not exist")
	}
}

// GetPage returns the page you have
func (p *Prices) GetPage() int {
	return p.pagination.Page
}

// GetLimit returns the limit you have provided, but if you have not, it provides the default
func (p *Prices) GetLimit() int {
	return p.pagination.Limit
}

func (p *Prices) GetAllPrices() ([]DataPrices1, error) {
	var resp []DataPrices1 = make([]DataPrices1, 0)
	var pricesPointer *Prices
	var err error
	if _, ok := p.args["limit"]; ok {
		value, err2 := strconv.Atoi(p.args["limit"])
		if err2 == nil {
			pricesPointer, err = p.client.GetPrices(args.Market(p.args["market"]), args.Type(p.args["type"]), args.Page(0), args.Limit(value))
		}
	} else {
		pricesPointer, err = p.client.GetPrices(args.Market(p.args["market"]), args.Type(p.args["type"]), args.Page(0))
	}

	for pricesPointer.pagination.Next != nil {
		if err != nil {
			return resp, fmt.Errorf("Error getting next page, %s", err)
		} else {
			resp = append(resp, pricesPointer.Data)
		}
		pricesPointer, err = pricesPointer.GetNext()
		time.Sleep(3 * time.Second)
	}
	resp = append(resp, pricesPointer.Data) // at the end of the loop, there is one left to add
	return resp, nil
}
