package client

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func assertSuccess(response string, acceptableError string, t *testing.T) {
	var resp map[string]interface{}
	json.Unmarshal([]byte(response), &resp)
	if val, ok := resp["status"]; ok {
		switch val {
		case "error":
			if val, ok := resp["message"]; !(ok && val == acceptableError) {
				t.Errorf(response)
			}
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

	client, err := New(apiKey, apiSecret)
	if err != nil {
		t.Errorf("error making the client")
	}
	//test every endpoint

	//test get methods
	t.Run("account", func(t *testing.T) {
		response, _ := client.Account()
		assertSuccess(response, "", t)
	})
	t.Run("wallet", func(t *testing.T) {
		response, _ := client.Balance()
		assertSuccess(response, "", t)
	})
	t.Run("transactions", func(t *testing.T) {
		var args = map[string]interface{}{
			"currency": "ETH",
		}
		response, _ := client.Transactions(args)
		assertSuccess(response, "", t)
	})
	t.Run("active orders=1", func(t *testing.T) {
		var args = map[string]interface{}{
			"market": "ETHCLP",
			"page":   "0",
		}
		response, _ := client.ActiveOrders(args)
		assertSuccess(response, "", t)
	})
	t.Run("active orders=2", func(t *testing.T) {
		var args = map[string]interface{}{
			"market": "ETHARS",
			"page":   "0",
		}
		response, _ := client.ActiveOrders(args)
		assertSuccess(response, "", t)
	})
	t.Run("order status", func(t *testing.T) {
		var args = map[string]interface{}{
			"id": "M103975",
		}
		response, _ := client.OrderStatus(args)
		assertSuccess(response, "invalid_scope", t)
	})
	t.Run("instant", func(t *testing.T) {
		var args = map[string]interface{}{
			"market": "ETHCLP",
			"type":   "sell",
			"amount": "159",
		}
		response, _ := client.Instant(args)
		assertSuccess(response, "", t)
	})

	t.Run("executed orders", func(t *testing.T) {
		var args = map[string]interface{}{
			"market": "ETHCLP",
			"page": "0",
		}
		response, _ := client.Instant(args)
		assertSuccess(response, "invalid_type", t)
	})

	//test post methods
	t.Run("create order", func(t *testing.T) {
		var args = map[string]interface{}{
			"amount": "0.3",
			"market": "ETHCLP",
			"price":  "10000",
			"type":   "buy",
		}
		response, _ := client.CreateOrder(args)
		assertSuccess(response, "not_enough_balance", t)
	})
	t.Run("cancel order", func(t *testing.T) {
		var args = map[string]interface{}{
			"id": "M103975",
		}
		response, _ := client.CancelOrder(args)
		assertSuccess(response, "invalid_request", t)
	})
	t.Run("create instant", func(t *testing.T) {
		var args = map[string]interface{}{
			"market": "ETHCLP",
			"type":   "buy",
			"amount": "10",
		}
		response, _ := client.CreateInstant(args)
		assertSuccess(response, "not_enough_balance", t)
	})
	t.Run("request deposit", func(t *testing.T) {
		var args = map[string]interface{}{
			"bank_account": "213213",
			"amount":       "10234",
		}
		response, _ := client.RequestDeposit(args)
		assertSuccess(response, "BankAccount matching query does not exist.", t)
	})
	t.Run("request withdrawal", func(t *testing.T) {
		var args = map[string]interface{}{
			"amount":       "10234",
			"bank_account": "213213",
		}
		response, _ := client.RequestWithdrawal(args)
		assertSuccess(response, "BankAccount matching query does not exist.", t)
	})
	t.Run("transfer", func(t *testing.T) {
		var args = map[string]interface{}{
			"address":  "GDMXNQBJMS3FYI4PFSYCCB4",
			"amount":   "1200",
			"currency": "XLM",
			"memo":     "162354",
		}
		response, _ := client.Transfer(args)
		assertSuccess(response, "max_limit_exceeded", t)
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

	client, err := New(apiKey, apiSecret)
	if err != nil {
		t.Errorf("error making the client")
	}
	//test every endpoint
	t.Run("new order", func(t *testing.T) {
		var args = map[string]interface{}{
			"callback_url":        "",
			"error_url":           "",
			"external_id":         "ABC123",
			"payment_receiver":    "user@email.com",
			"success_url":         "",
			"to_receive":          "3000",
			"to_receive_currency": "CLP",
		}
		response, _ := client.NewOrder(args)
		assertSuccess(response, "invalid_request", t)
	})
	t.Run("create wallet", func(t *testing.T) {
		var args = map[string]interface{}{
			"id":     "P2023132",
			"token":  "xToY232aheSt8F",
			"wallet": "ETH",
		}
		response, _ := client.CreateWallet(args)
		assertSuccess(response, "payment_does_not_exist", t)
	})
	t.Run("payment orders", func(t *testing.T) {
		var args = map[string]interface{}{
			"start_date": "01/03/2018",
			"end_date":   "08/03/2018",
		}
		response, _ := client.PaymentOrders(args)
		assertSuccess(response, "", t)
	})
	t.Run("payment status", func(t *testing.T) {
		var args = map[string]interface{}{
			"id": "P13433",
		}
		response, _ := client.PaymentStatus(args)
		assertSuccess(response, "invalid_scope", t)
	})

}
