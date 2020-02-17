package conn 

import (
	"github.com/cryptomkt/cryptomkt-go/args"
)

type Pagination struct {
	previous int
	next int
	limit int
	page int
}

type Amount  struct{
	Original string
	Remaining string
	Executed string
}

type AnOrder struct {
	apiClient *Client
	status string
	pagination Pagination
	warnings string
	data OrderData
}
type OrderData struct {
	
	Id string
	Status string
	Type string
	Price string
	Amount Amount
	ExecutionPrice string `json:"execution_price"`
	AvgExecutionPrice int `json:"avg_execution_price"`
	Market string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ExecutedAt string `json:"executed_at"`

}


func (order *AnOrder) Refresh() {
	order.apiClient.OrderStatus(
		args.Id(order.data.Id))
}

func (order *AnOrder) Cancel() {
	order.apiClient.CancelOrder(
		args.Id(order.data.Id))
}