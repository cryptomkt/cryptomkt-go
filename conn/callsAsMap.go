package conn

import (
	"github.com/cryptomkt/cryptomkt-go/args"
)

// GetTransactionsAsMapList returns the movements of the wallets of the client.
// Returns the data in a slice of maps: []map[string]string.
//
// List of accepted Arguments:
//   - required: Currency
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#obtener-movimientos
func (client *Client) GetTransactionsAsMapList(args ...args.Argument) ([]map[string]string, error) {
	transactions, err := client.GetTransactions(args...)
	if err != nil {
		return nil, err
	}
	mapList := make([]map[string]string, 0)
	for _, transaction := range *transactions {
		mapList = append(mapList, TransactionToMap(transaction))
	}
	return mapList, nil

}

// GetActiveOrdersAsMapList returns the list of active orders of the client.
// Returns the data in a slice of maps: []map[string]string.
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-activas
func (client *Client) GetActiveOrdersAsMapList(args ...args.Argument) ([]map[string]string, error) {
	activeOrders, err := client.GetActiveOrders(args...)
	if err != nil {
		return nil, err
	}
	mapList := make([]map[string]string, 0)
	for _, order := range activeOrders {
		mapList = append(mapList, OrderToMap(order))
	}
	return mapList, nil
}

// GetExecutedOrdersAsMapList return a list of the executed orders of the client.
// Returns the data in a slice of maps: []map[string]string.
//
// List of accepted Arguments:
//   - required: Market
//   - optional: Page, Limit
// https://developers.cryptomkt.com/es/#ordenes-ejecutadas
func (client *Client) GetExecutedOrdersAsMapList(args ...args.Argument) ([]map[string]string, error) {
	activeOrders, err := client.GetExecutedOrders(args...)
	if err != nil {
		return nil, err
	}
	mapList := make([]map[string]string, 0)
	for _, order := range activeOrders {
		mapList = append(mapList, OrderToMap(order))
	}
	return mapList, nil
}

// GetOrderStatusAsMap gives the status of an order.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden
func (client *Client) GetOrderStatusAsMap(args ...args.Argument) (map[string]string, error) {
	order, err := client.GetOrderStatus(args...)
	if err != nil {
		return nil, err
	}
	return OrderToMap(*order), nil
}

// GetInstantAsMap emulates an order in the current state of the Instant Exchange of CryptoMarket.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: Market, Type, Amount
//   - optional: none
// https://developers.cryptomkt.com/es/#obtener-cantidad
func (client *Client) GetInstantAsMap(args ...args.Argument) (map[string]string, error) {
	instant, err := client.GetInstant(args...)
	if err != nil {
		return nil, err
	}
	return InstantToMap(*instant), nil
}

// CreateOrderAsMap signal to create an order of buy or sell in CryptoMarket.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: Amount, Market, Price, Type
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-orden
func (client *Client) CreateOrderAsMap(args ...args.Argument) (map[string]string, error) {
	order, err := client.CreateOrder(args...)
	if err != nil {
		return nil, err
	}
	return OrderToMap(*order), nil
}

// CancelOrderAsMap signal to cancel an order in CryptoMarket.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#cancelar-una-orden
func (client *Client) CancelOrderAsMap(args ...args.Argument) (map[string]string, error) {
	order, err := client.CancelOrder(args...)
	if err != nil {
		return nil, err
	}
	return OrderToMap(*order), nil
}

// NewOrderAsMap enables a payment order, and gives a QR and urls.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: ToReceive, ToReceiveCurrency, PaymentReceiver
//   - optional: ExternalId, CallbackUrl, ErrorUrl, SuccessUrl, RefundEmail, Language
// https://developers.cryptomkt.com/es/#crear-orden-de-pago
func (client *Client) NewOrderAsMap(args ...args.Argument) (map[string]string, error) {
	payment, err := client.NewOrder(args...)
	if err != nil {
		return nil, err
	}
	return PaymentOrderToMap(*payment), nil
}

// CreateWalletAsMap creates a wallet to pay a payment order.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: Id, Token, Wallet
//   - optional: none
// https://developers.cryptomkt.com/es/#crear-billetera-de-orden-de-pago
func (client *Client) CreateWalletAsMap(args ...args.Argument) (map[string]string, error) {
	payment, err := client.CreateWallet(args...)
	if err != nil {
		return nil, err
	}
	return PaymentOrderToMap(*payment), nil
}

// PaymentOrdersAsMapList returns all the requested payment orders.
// Returns the data in a slice of maps: []map[string]string.
//
// List of accepted Arguments:
//   - required: StartDate, EndDate
//   - optional: none
// https://developers.cryptomkt.com/es/#listado-de-ordenes-de-pago
func (client *Client) PaymentOrdersAsMapList(args ...args.Argument) ([]map[string]string, error) {
	paymentList, err := client.PaymentOrders(args...)
	if err != nil {
		return nil, err
	}
	mapList := make([]map[string]string, 0)
	for _, payment := range paymentList{
		mapList = append(mapList, PaymentOrderToMap(payment))
	}
	return mapList, nil
}

// PaymentStatusAsMap gives the status of a pyment order.
// Returns the data in a map: map[string]string.
//
// List of accepted Arguments:
//   - required: Id
//   - optional: none
// https://developers.cryptomkt.com/es/#estado-de-orden-de-pago
func (client *Client) GetPaymentStatusAsMap(args ...args.Argument) (map[string]string, error) {
	payment, err := client.GetPaymentStatus(args...)
	if err != nil {
		return nil, err
	}
	return PaymentOrderToMap(*payment), nil
}
