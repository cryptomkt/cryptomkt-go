package client

import (
    "os"
    "log"
	"bufio"
	"testing"
	"encoding/json"
)

func assertSuccess(response string, t *testing.T) {
	var resp map[string] interface{}
		json.Unmarshal([]byte(response), &resp)
		if val, ok := resp["status"]; ok {
			switch val {
			case "error":
				t.Errorf(response)
			case "success":
				//all good
			default:
				t.Errorf("unexpected response")
			}
		} else {
			t.Errorf("error in the response")
		}
}

func TestAutenticated(t *testing.T) {
	//load the keys for the connection to crypto market.
	//apikey first, apisecret after,  every in its own line.
	file, err := os.Open("keys.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	scanner.Scan()
	apiKey := scanner.Text()
	scanner.Scan()
	apiSecret := scanner.Text()
	
	client, err := NewClient(apiKey, apiSecret)
	if err != nil {
		t.Errorf("error making the client")
	}
	//test every endpoint

	//test get methods
	t.Run("account", func(t *testing.T){
		response := client.Account()
		assertSuccess(response, t)	
	})
	t.Run("wallet", func(t *testing.T){
		response := client.Balance()
		assertSuccess(response, t)	
	})
	t.Run("transactions", func(t *testing.T){
		var args = map[string]string{
			"currency":"ETH",
		}
		response := client.Transactions(args)
		assertSuccess(response, t)	
	})
	t.Run("active orders=1", func(t *testing.T){
		var args = map[string]string{
			"market":"ETHCLP",
			"page":"0",
		}
		response := client.ActiveOrders(args)
		assertSuccess(response, t)	
	})
	t.Run("active orders=2", func(t *testing.T){
		var args = map[string]string{
			"market":"ETHARS",
			"page":"0",
		}
		response := client.ActiveOrders(args)
		assertSuccess(response, t)	
	})
	t.Run("order status", func(t *testing.T){
		var args = map[string]string{
			"id":"M103975",
		}
		response := client.OrderStatus(args)
		assertSuccess(response, t)	
	})
	t.Run("instant", func(t *testing.T){
		var args = map[string]string{
			"market":"ETHCLP",
			"type":"sell",
			"amount":"159",
		}
		response := client.Instant(args)
		assertSuccess(response, t)	
	})

	//test post methods
	t.Run("create order", func(t *testing.T){
		var args = map[string]string{
			"amount": "0.3",
			"market": "ETHCLP",
			"price": "10000",
			"type": "buy",
		}
		response := client.CreateOrder(args)
		assertSuccess(response, t)	
	})
	t.Run("cancel order", func(t *testing.T){
		var args = map[string]string{
			"id":"M103975",
		}
		response := client.CancelOrder(args)
		assertSuccess(response, t)	
	})
	t.Run("create instant", func(t *testing.T){
		var args = map[string]string{
			"market": "ETHCLP", 
			"type": "buy", 
			"amount": "10",
		}
		response := client.CreateInstant(args)
		assertSuccess(response, t)	
	})
	t.Run("request deposit", func(t *testing.T){
		var args = map[string]string{
			"bank_account": "213213", 
			"amount": "10234",
		}
		response := client.RequestDeposit(args)
		assertSuccess(response, t)	
	})
	t.Run("request withdrawal", func(t *testing.T){
		var args = map[string]string{
			"amount": "10234",
			"bank_account": "213213", 
		}
		response := client.RequestWithdrawal(args)
		assertSuccess(response, t)	
	})
	t.Run("transfer", func(t *testing.T){
		var args = map[string]string{
			"address": "GDMXNQBJMS3FYI4PFSYCCB4",
			"amount": "1200",
			"currency": "XLM",
			"memo": "162354",
		}
		response := client.Transfer(args)
		assertSuccess(response, t)	
	})
}

func TestCryptoCompra(t *testing.T) {
	//load the keys for the connection to crypto market.
	//apikey first, apisecret after,  every in its own line.
	file, err := os.Open("keys.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	apiKey := scanner.Text()
	scanner.Scan()
	apiSecret := scanner.Text()
	
	client, err := NewClient(apiKey, apiSecret)
	if err != nil {
		t.Errorf("error making the client")
	}
	//test every endpoint
	t.Run("new order", func(t *testing.T){
		var args = map[string]string{
			"callback_url": "",
			"error_url": "",
			"external_id": "ABC123",
			"payment_receiver": "user@email.com",
			"success_url": "",
			"to_receive": "3000",
			"to_receive_currency": "CLP",
		}
		response := client.NewOrder(args)
		assertSuccess(response, t)	
	})
	t.Run("create wallet", func(t *testing.T){
		var args = map[string]string{
			"id": "P2023132",
			"token": "xToY232aheSt8F",
			"wallet": "ETH",
		}
		response := client.CreateWallet(args)
		assertSuccess(response, t)	
	})
	t.Run("payment orders", func(t *testing.T){
		var args = map[string]string{
			"start_date": "01/03/2018", 
    		"end_date": "08/03/2018",
		}
		response := client.PaymentOrders(args)
		assertSuccess(response, t)	
	})
	t.Run("payment status", func(t *testing.T){
		var args = map[string]string{
			"id": "P13433",
		}
		response := client.PaymentStatus(args)
		assertSuccess(response, t)	
	})

}
