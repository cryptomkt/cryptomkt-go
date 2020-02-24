package conn

import (
	"fmt"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
)

// GetTradesAllPages returns a pointer to a Trades struct with the data given
// by the api and an error message. It returns (nil, error) when an error
// is raised and (*Trades, nil) when the operation is successful. If end
// argument is not provided, the maximum data amount will be trucated when
// it raises more than 100 elements. It is not sure it will give you
// exactly 100 TradeData Data.
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Start, End
// https://developers.cryptomkt.com/es/#trades
func (client *Client) GetTradesAllPages(arguments ...args.Argument) ([]TradeData, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTradesAllPages: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	neededArguments = append(neededArguments, args.Market(argsMap["market"]))
	if val, ok := argsMap["start"]; ok {
		neededArguments = append(neededArguments, args.Start(val))
	}
	if val, ok := argsMap["end"]; ok {
		neededArguments = append(neededArguments, args.End(val))
	}

	tPage, err := client.GetTrades(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTrades: %s", err)
	}
	allt := make([]TradeData, len(tPage.Data))
	copy(allt, tPage.Data)
	for tPage, err = tPage.GetNext(); err == nil; tPage, err = tPage.GetNext() {
		time.Sleep(2 * time.Second) //because the server only accepts 30 calls per minute.
		allt = append(allt, tPage.Data...)
		// When the data length raises 100 elements or more when "end" parameter is not provided,
		// it breaks. This block limit the number of pages
		if _, ok := argsMap["end"]; !ok {
			if len(allt) > 100 {
				break
			}
		}
	}
	return allt, nil
}

// GetActiveOrdersAllPages gets all the actives orders of the client in a given market.
// returns an array of orders.
//
// List of accepted Arguments:
//   - required: Market
//   - optional: none
// https://developers.cryptomkt.com/es/#ordenes-activas
func (client *Client) GetActiveOrdersAllPages(arguments ...args.Argument) ([]Order, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllActiveOrders: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["market"]
	neededArguments = append(neededArguments, args.Market(val))

	oList, err := client.GetActiveOrders(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetActiveOrdersAllPages: %s", err)
	}
	return getAllOrders(oList), nil
}

// GetExecutedOrdersAllPages gets all executed orders of the client in a given market
//
// List of accepted Arguments:
//   - required: Market
//   - optional: none
// https://developers.cryptomkt.com/es/#ordenes-ejecutadas
func (client *Client) GetExecutedOrdersAllPages(arguments ...args.Argument) ([]Order, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetExecutedOrdersAllPages: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["market"]
	neededArguments = append(neededArguments, args.Market(val))

	oList, err := client.GetExecutedOrders(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllExecutedOrders: %s", err)
	}
	return getAllOrders(oList), nil
}

// PaymentOrdersAllPages get all the payment orders between the two given dates.
// Returns an array of PaymentOrder
//
// List of accepted Arguments:
//   - required: StartDate, EndDate
//   - optional: none
func (client *Client) PaymentOrdersAllPages(arguments ...args.Argument) ([]PaymentOrder, error) {
	req, err := makeReq([]string{"start_date", "end_date"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in PaymentOrdersAllPages: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["start_date"]
	neededArguments = append(neededArguments, args.StartDate(val))
	val = argsMap["end_date"]
	neededArguments = append(neededArguments, args.EndDate(val))

	poList, err := client.GetPaymentOrders(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetPaymentOrders: %s", err)
	}
	allpo := make([]PaymentOrder, len(poList.Data))
	copy(allpo, poList.Data)
	for poList, err = poList.GetNext(); err == nil; poList, err = poList.GetNext() {
		time.Sleep(2 * time.Second)
		allpo = append(allpo, poList.Data...)
		// When the data length raises 100 elements or more,
		// it breaks. This "if" block limit the number of pages
		if len(allpo) > 100 {
			break
		}
	}
	return allpo, nil
}

func getAllOrders(oList *OrderList) []Order {
	allo := make([]Order, len(oList.Data))
	copy(allo, oList.Data)
	for oList, err := oList.GetNext(); err == nil; oList, err = oList.GetNext() {
		time.Sleep(2 * time.Second)
		oList.setClientInOrders()
		allo = append(allo, oList.Data...)

		// When the data pages length raises 100 elements or more,
		// it breaks. This "if" block limit the number of pages
		if len(allo) > 100 {
			break
		}
	}
	return allo
}

func (client *Client) GetAllTransactions(argus ...args.Argument) ([]Transaction, error) {
	req, err := makeReq([]string{"currency"}, argus...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllTransactions: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	neededArguments = append(neededArguments, args.Currency(argsMap["currency"]))

	trans, err := client.GetTransactions(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetTransactions: %s", err)
	}
	allTrans := make([]Transaction, len(trans.Data))
	copy(allTrans, trans.Data)
	for trans, err = trans.GetNext(); err == nil; trans, err = trans.GetNext() {
		time.Sleep(2 * time.Second)
		allTrans = append(allTrans, trans.Data...)
		if len(allTrans) > 100 {
			break
		}
	}
	return allTrans, nil
}
