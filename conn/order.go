package conn

import (
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/args"
)

type OrderListResp struct {
	Status     string
	Message    string
	Pagination Pagination
	Warnings   string
	Data       []Order
}

type OrderResponse struct {
	Status  string
	Message string
	Data    Order
}

type OrderList struct {
	client     *Client
	caller     string
	market     string
	Status     string
	pagination Pagination
	Warnings   string
	Data       []Order
}

type Order struct {
	client            *Client
	Id                string
	Status            string
	Type              string
	Price             string
	Amount            Amount
	ExecutionPrice    string `json:"execution_price"`
	AvgExecutionPrice int    `json:"avg_execution_price"`
	Market            string
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ExecutedAt        string `json:"executed_at"`
}

type Amount struct {
	Original  string
	Remaining string
	Executed  string
}

func (o *Order) Close() error {
	oClosed, err := o.client.CancelOrder(args.Id(o.Id))
	if err != nil {
		return fmt.Errorf("Close order %s failed: %s", o.Id, err)
	}
	o =  oClosed
	return nil
}

func (o *Order) Refresh() error {
	oRefreshed, err := o.client.GetOrderStatus(args.Id(o.Id))
	if err != nil {
		return fmt.Errorf("Refresh order %s failed: %s", o.Id, err)
	}
	o = oRefreshed
	return nil
}

func (oList *OrderList) Close() error {
	for i, order := range oList.Data {
		oClosed, err := oList.client.CancelOrder(args.Id(order.Id))
		if err != nil {
			return fmt.Errorf("Close order %s failed: %s", order.Id, err)
		}
		oList.Data[i] = *oClosed
	}
	return nil
}

func (oList *OrderList) Refresh() error {
	for i, order := range oList.Data {
		oRefreshed, err := oList.client.GetOrderStatus(args.Id(order.Id))
		if err != nil {
			return fmt.Errorf("Refresh order %s failed: %s", order.Id, err)
		}
		oList.Data[i] = *oRefreshed
	}
	return nil
}

func (o *OrderList) GetPrevious() (*OrderList, error) {
	if o.pagination.Next == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	if o.caller == "active_orders" {
		return o.client.GetActiveOrdersPage(
			args.Market(o.market),
			args.Page(int(o.pagination.Previous.(float64))),
			args.Limit(o.pagination.Limit))
	}
	return o.client.GetExecutedOrdersPage(
		args.Market(o.market),
		args.Page(int(o.pagination.Previous.(float64))),
		args.Limit(o.pagination.Limit))
}

// GetNext lets you go to the next page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (o *OrderList) GetNext() (*OrderList, error) {
	if o.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	if o.caller == "active_orders" {
		return o.client.GetActiveOrdersPage(
			args.Market(o.market),
			args.Page(int(o.pagination.Next.(float64))),
			args.Limit(o.pagination.Limit))
	}
	return o.client.GetExecutedOrdersPage(
		args.Market(o.market),
		args.Page(int(o.pagination.Next.(float64))),
		args.Limit(o.pagination.Limit))
}

// GetAllPaymentOrders get all the payment orders between the two given dates.
// Returns an array of PaymentOrder
//
// List of accepted Arguments:
//   - required: Market
//   - optional: none
func (client *Client) GetExecutedOrders(arguments... args.Argument) (*[]Order, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllExecutedOrders: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["market"]
	neededArguments = append(neededArguments, args.Market(val))

	oList, err := client.GetExecutedOrdersPage(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllExecutedOrders: %s", err)
	}
	return getAllOrders(oList), nil
}

func (client *Client) GetActiveOrders(arguments... args.Argument) (*[]Order, error) {
	req, err := makeReq([]string{"market"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllActiveOrders: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["market"]
	neededArguments = append(neededArguments, args.Market(val))

	oList, err := client.GetActiveOrdersPage(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllActiveOrders: %s", err)
	}
	return getAllOrders(oList), nil
}

func getAllOrders(oList *OrderList) (*[]Order) {
	allo := make([]Order, len(oList.Data))
	copy(allo, oList.Data)
	for oList, err := oList.GetNext(); err == nil; oList, err = oList.GetNext() {
		oList.setClientInOrders()
		allo = append(allo, oList.Data...)
	}
	return &allo
}

func (oList *OrderList) setClientInOrders() {
	for _, order := range oList.Data {
		order.client = oList.client
	}
}