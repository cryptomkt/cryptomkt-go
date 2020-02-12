package client

// Account gives the account information of the client.
// https://developers.cryptomkt.com/es/#cuenta
func (client *Client) Account() (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{},
	}
	return client.get("account", req)
}

// Balance returns the actual balance of the wallets of the client in cryptomarket
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) Balance() (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{},
	}
	return client.get("balance", req)
}

// Wallets is an alias for Balance
// https://developers.cryptomkt.com/es/#obtener-balance
func (client *Client) Wallets() (string, error) {
	return client.Balance()
}

func (client *Client) Transactions(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"currency"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("transactions", req)
}

// ActiveOrders returns the active orders of the client
// given a market and an opitonal page number
// https://developers.cryptomkt.com/es/#ordenes-activas
func (client *Client) ActiveOrders(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"market"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("orders/active", req)
}

// ExecutedOrders return a list of the executed orders of the client
// given a market and an optional page
// https://developers.cryptomkt.com/es/#ordenes-ejecutadas
func (client *Client) ExecutedOrders(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"market"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("orders/executed", req)
}


// OrderStatus gives the status of an order
// given the order "id"
// https://developers.cryptomkt.com/es/#estado-de-orden
func (client *Client) OrderStatus(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"id"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("orders/status", req)
}

// Instant emulates an order in the current state of the Instant Exchange of CryptoMarket
// given a "market", a "type" ("buy" or "sell") and an "amount"
// https://developers.cryptomkt.com/es/#obtener-cantidad
func (client *Client) Instant(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"market", "type", "amount"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("orders/instant/get", req)
}

// CreateOrder signal to create an order of buy or sell in CryptoMarket
// given an "amount", a "market", a "price" and a "type" ("buy" or "sell")
// https://developers.cryptomkt.com/es/#crear-orden
func (client *Client) CreateOrder(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"amount", "market", "price", "type"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("orders/create", req)
}

// CancelOrder signal to cancel an order in CryptoMarket
// given the order "id"
// https://developers.cryptomkt.com/es/#cancelar-una-orden
func (client *Client) CancelOrder(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"id"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("orders/cancel", req)
}

// CreateInstant makes an order in the Instant Exchange of CryptoMarket
// given a "market", a "type" ("buy" or "sell") and an "amount"
// https://developers.cryptomkt.com/es/#crear-orden-2
func (client *Client) CreateInstant(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"market", "type", "amount"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("orders/instant/create", req)
}

// RequestDeposit notifies a deposit to a wallet of local currency
// given an "amount", a "bank_account", a "date",
// a "tracking_code" (only for mexico)
// and a "voucher" (for Mexico, Brasil and European Union only)
// https://developers.cryptomkt.com/es/#notificar-deposito
func (client *Client) RequestDeposit(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"amount", "bank_account"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("request/deposit", req)
}

// Request withdrawal notifies a withdrawal from a bank account of the client
// given a "bank_account" and an "amount"
// https://developers.cryptomkt.com/es/#notificar-retiro
func (client *Client) RequestWithdrawal(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"amount", "bank_account"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("request/withdrawal", req)
}

// Transfer move crypto between wallets
// given the "adress" of the destiny wallet,
// the "amount", a "currency" of the origin wallet,
// and an optional "memo"
// https://developers.cryptomkt.com/es/#transferir
func (client *Client) Transfer(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"address", "amount", "currency"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("transfer", req)
}

// NewOrder enables a payment order, and gives a QR and urls
// https://developers.cryptomkt.com/es/#crear-orden-de-pago
func (client *Client) NewOrder(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"to_receive", "to_receive_currency", "payment_receiver"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("payment/new_order", req)
}

// CreateWallet creates a wallet to pay a payment order
// https://developers.cryptomkt.com/es/#crear-billetera-de-orden-de-pago
func (client *Client) CreateWallet(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"id", "token", "wallet"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.post("payment/create_wallet", req)
}

// PaymentOrders returns all the generated payment orders
// https://developers.cryptomkt.com/es/#listado-de-ordenes-de-pago
func (client *Client) PaymentOrders(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"start_date", "end_date"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("payment/orders", req)
}

// PaymentStatus gives the status of a pyment order
// given the order "id"
// https://developers.cryptomkt.com/es/#estado-de-orden-de-pago
func (client *Client) PaymentStatus(args ...Argument) (string, error) {
	req := &Request{
		arguments: make(map[string]string),
		required: []string{"id"},
	}
	for _, argument := range args {
		argument(req)
	}
	err := req.assertRequired()
	if err != nil {
		return "", err
	}
	return client.get("payment/status", req)
}
