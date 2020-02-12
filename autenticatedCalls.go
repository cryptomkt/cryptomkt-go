package client

func (client *Client) Account() string {
	return client.get("account", nil)
}

func (client *Client) Balance() string {
	return client.get("balance", nil)
}

func (client *Client) Wallet() string {
	return client.Balance()
}

func (client *Client) Transactions(args map[string]interface{}) string {
	return client.get("transactions", args)
}

func (client *Client) ActiveOrders(args map[string]interface{}) string {
	return client.get("orders/active", args)
}

func (client *Client) OrderStatus(args map[string]interface{}) string {
	return client.get("orders/status", args)
}

func (client *Client) Instant(args map[string]interface{}) string {
	return client.get("orders/instant/get", args)
}

//post methods
func (client *Client) CreateOrder(args map[string]interface{}) string {
	return client.post("orders/create", args)
}

func (client *Client) CancelOrder(args map[string]interface{}) string {
	return client.post("orders/cancel", args)
}

func (client *Client) CreateInstant(args map[string]interface{}) string {
	return client.post("orders/instant/create", args)
}

func (client *Client) RequestDeposit(args map[string]interface{}) string {
	return client.post("request/deposit", args)
}

func (client *Client) RequestWithdrawal(args map[string]interface{}) string {
	return client.post("request/withdrawal", args)
}

func (client *Client) Transfer(args map[string]interface{}) string {
	return client.post("transfer", args)
}

//crypto compra
//post
func (client *Client) NewOrder(args map[string]interface{}) string {
	return client.post("payment/new_order", args)
}

func (client *Client) CreateWallet(args map[string]interface{}) string {
	return client.post("payment/create_wallet", args)
}

//get
func (client *Client) PaymentOrders(args map[string]interface{}) string {
	return client.get("payment/orders", args)
}

func (client *Client) PaymentStatus(args map[string]interface{}) string {
	return client.get("payment/status", args)
}
