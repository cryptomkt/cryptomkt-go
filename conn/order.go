package conn

import (
	"bytes"
	"fmt"
	"strconv"

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

// Close closes the calling order, and changes the order to reflect
// the new state of the order, after being closed.
// Calls CancelOrder with the asociated client of the order.
// https://developers.cryptomkt.com/es/#cancelar-una-orden
func (o *Order) Close() (*Order, error) {
	oClosed, err := o.client.CancelOrder(args.Id(o.Id))
	if err != nil {
		return nil, fmt.Errorf("Close order %s failed: %s", o.Id, err)
	}
	oClosed.client = o.client
	return oClosed, nil
}

// Refresh refreshes the calling order, and changes it to be the actual
// state of the order.
// Calls GetOrderStatus with the asociated client of the order.
// https://developers.cryptomkt.com/es/#estado-de-orden
func (o *Order) Refresh() (*Order, error) {
	oRefreshed, err := o.client.GetOrderStatus(args.Id(o.Id))
	if err != nil {
		return nil, fmt.Errorf("Refresh order %s failed: %s", o.Id, err)
	}
	oRefreshed.client = o.client
	return oRefreshed, nil
}

// Close closes every order in the order list.
func (oList *OrderList) Close() error {
	for i, order := range oList.Data {
		oClosed, err := oList.client.CancelOrder(args.Id(order.Id))
		if err != nil {
			return fmt.Errorf("Close order %s failed: %s", order.Id, err)
		}
		oList.Data[i] = *oClosed
	}
	oList.setClientInOrders()
	return nil
}

// Refresh refreshes every order in the order list.
// its an iterative implementation, so if an error is rised refreshing
// some order, the preciding orders end refreshed.
func (oList *OrderList) Refresh() error {
	for i, order := range oList.Data {
		oRefreshed, err := oList.client.GetOrderStatus(args.Id(order.Id))
		if err != nil {
			return fmt.Errorf("Refresh order %s failed: %s", order.Id, err)
		}
		oList.Data[i] = *oRefreshed
	}
	oList.setClientInOrders()
	return nil
}

// GetPrevious get the previous page of the List of orders.
// If there is no previous page, rise an error.
func (o *OrderList) GetPrevious() (*OrderList, error) {
	if o.pagination.Previous == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	var call func(args ...args.Argument) (*OrderList, error)
	if o.caller == "active_orders" {
		call = o.client.GetActiveOrders
	} else { // caller is execute_order
		call = o.client.GetExecutedOrders
	}
	oList, err := call(
		args.Market(o.market),
		args.Page(int(o.pagination.Previous.(float64))),
		args.Limit(o.pagination.Limit))
	if err != nil {
		return nil, fmt.Errorf("error getting the previous page: %s", err)
	}
	oList.setClientInOrders()
	return oList, nil
}

// GetNext lets you go to the next page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (o *OrderList) GetNext() (*OrderList, error) {
	if o.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	var call func(args ...args.Argument) (*OrderList, error)
	if o.caller == "active_orders" {
		call = o.client.GetActiveOrders
	} else { // caller is execute_order
		call = o.client.GetExecutedOrders
	}
	oList, err := call(
		args.Market(o.market),
		args.Page(int(o.pagination.Next.(float64))),
		args.Limit(o.pagination.Limit))
	if err != nil {
		return nil, fmt.Errorf("error getting the next page: %s", err)
	}
	oList.setClientInOrders()
	return oList, nil
}

func (oList *OrderList) setClientInOrders() {
	for i, _ := range oList.Data {
		oList.Data[i].client = oList.client
	}
}

func (oList *OrderList) String() string {
	var b bytes.Buffer
	b.WriteString("OrderList{")

	b.WriteString("\n\tclientAddr:")
	b.WriteString(oList.client.String())

	b.WriteString("\n\tcaller:")
	b.WriteString(oList.caller)

	b.WriteString("\n\tmarket:")
	b.WriteString(oList.market)

	b.WriteString("\n\tpagination:")
	b.WriteString(oList.pagination.String())

	b.WriteString("\n\torders:")

	for _, order := range oList.Data {
		b.WriteString("\n")
		b.WriteString(order.String())
	}
	b.WriteString("\n}")
	return b.String()
}
func (order *Order) String() string {
	var b bytes.Buffer
	b.WriteString("Order{")

	b.WriteString("\n\tclientAddr:")
	b.WriteString(order.client.String())

	b.WriteString("\n\tId:")
	b.WriteString(order.Id)

	b.WriteString("\n\tstatus:")
	b.WriteString(order.Status)

	b.WriteString("\n\ttype:")
	b.WriteString(order.Type)

	b.WriteString("\n\tprice:")
	b.WriteString(order.Price)

	b.WriteString("\n\tamount:")
	b.WriteString(order.Amount.String())

	b.WriteString("\n\texecutionPrice:")
	b.WriteString(order.ExecutionPrice)

	b.WriteString("\n\tavgExecutionPrice:")
	b.WriteString(strconv.Itoa(order.AvgExecutionPrice))

	b.WriteString("\n\tmarket:")
	b.WriteString(order.Market)

	b.WriteString("\n\tcreatedAt:")
	b.WriteString(order.CreatedAt)

	b.WriteString("\n\tupdatedAt:")
	b.WriteString(order.UpdatedAt)

	b.WriteString("\n\texecutedAt:")
	b.WriteString(order.ExecutedAt)

	b.WriteString("\n}")
	return b.String()
}

func (amount *Amount) String() string {
	return "Amount{original:" +
		amount.Original +
		" remaining:" +
		amount.Remaining +
		" executed:" +
		amount.Executed +
		"}"
}
