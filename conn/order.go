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

func (o *OrderList) GetPrevious() (*OrderList, error) {
	if o.pagination.Next == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	if o.caller == "active_orders" {
		return o.client.GetActiveOrders(
			args.Market(o.market),
			args.Page(o.pagination.Previous.(int)),
			args.Limit(o.pagination.Limit))
	}
	return o.client.GetExecutedOrders(
		args.Market(o.market),
		args.Page(o.pagination.Previous.(int)),
		args.Limit(o.pagination.Limit))
}

// GetNext lets you go to the next page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (o *OrderList) GetNext() (*OrderList, error) {
	if o.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	if o.caller == "active_orders" {
		return o.client.GetActiveOrders(
			args.Market(o.market),
			args.Page(o.pagination.Next.(int)),
			args.Limit(o.pagination.Limit))
	}
	return o.client.GetExecutedOrders(
		args.Market(o.market),
		args.Page(o.pagination.Next.(int)),
		args.Limit(o.pagination.Limit))
}