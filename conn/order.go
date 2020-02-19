package conn

import (
	"github.com/cryptomkt/cryptomkt-go/args"
)

func (order *Order) Refresh() error {
	_, err := order.client.GetOrderStatus(args.Id(order.Id))
	return err
}

func (order *Order) Close() error {
	_, err := order.client.CancelOrder(args.Id(order.Id))
	return err
}
