package conn

import (
	"github.com/cryptomkt/cryptomkt-go/args"
)

type Pagination struct {
	PreviousHolder interface{} `json:"previous"`
	NextHolder     interface{} `json:"next"`
	Previous       int
	Next           int
	Limit          int
	Page           int
}

type Amount struct {
	Original  string
	Remaining string
	Executed  string
}

type OrderListResp struct {
	apiClient  *Client
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

type Order struct {
	apiClient         *Client
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

func (order *Order) Refresh() {
	order.apiClient.GetOrderStatus(
		args.Id(order.Id))
}

func (order *Order) Close() {
	order.apiClient.CancelOrder(
		args.Id(order.Id))
}
