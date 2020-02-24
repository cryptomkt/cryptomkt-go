package conn

import (
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/args"
	"time"
)

// GetTradesAllPages returns a pointer to a Trades struct with the data given
// by the api and an error message. It returns (nil, error) when an error
// is raised and (*Trades, nil) when the operation is successful.
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Start, End
// https://developers.cryptomkt.com/es/#trades
func (client *Client) GetTradesAllPages(arguments ...args.Argument) ([]TradeData, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetPaymentOrders: %s", err)
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
		return nil, fmt.Errorf("Error in GetPaymentOrders: %s", err)
	}
	allt := make([]TradeData, len(tPage.Data))
	copy(allt, tPage.Data)

	for tPage, err = tPage.GetNext(); err == nil; tPage, err = tPage.GetNext() {
		time.Sleep(2 * time.Second)
		allt = append(allt, tPage.Data...)
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
		return nil, fmt.Errorf("Error in GetAllActiveOrders: %s", err)
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
		return nil, fmt.Errorf("Error in GetAllExecutedOrders: %s", err)
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
		return nil, fmt.Errorf("Error in GetPaymentOrders: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["start_date"]
	neededArguments = append(neededArguments, args.StartDate(val))
	val = argsMap["end_date"]
	neededArguments = append(neededArguments, args.EndDate(val))

	poList, err := client.PaymentOrders(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetPaymentOrders: %s", err)
	}
	allpo := make([]PaymentOrder, len(poList.Data))
	copy(allpo, poList.Data)
	for poList, err = poList.GetNext(); err == nil; poList, err = poList.GetNext() {
		time.Sleep(2 * time.Second)
		allpo = append(allpo, poList.Data...)
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
	}
	return allo
}
