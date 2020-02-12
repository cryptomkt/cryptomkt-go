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
		t.Errorf("error in the response: %s", resp)
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
		response, _ := client.Transactions(
			Currency("ETH"))
		assertSuccess(response, "", t)
	})
	t.Run("active orders=1", func(t *testing.T) {
		response, _ := client.ActiveOrders(
			Market("ETHCLP"),
			Page("0"))
		assertSuccess(response, "", t)
	})
	t.Run("active orders=2", func(t *testing.T) {
		response, _ := client.ActiveOrders(
			Market("ETHARS"), 
			Page("1"))
		assertSuccess(response, "", t)
	})
	t.Run("order status", func(t *testing.T) {
		response, _ := client.OrderStatus(
			Id("M103975"))
		assertSuccess(response, "invalid_scope", t)
	})
	t.Run("instant", func(t *testing.T) {
		response, _ := client.Instant(
			Market("ETHCLP"), 
			Type("sell"), 
			Amount("159"))
		assertSuccess(response, "", t)
	})

	t.Run("executed orders", func(t *testing.T) {
		response, _ := client.ExecutedOrders(
			Market("ETHCLP"), 
			Page("0"))
		assertSuccess(response, "invalid_type", t)
	})

	//test post methods
	t.Run("create order", func(t *testing.T) {
		response, _ := client.CreateOrder(
			Amount("0.3"),
			Market("ETHCLP"),
			Price("1000"),
			Type("buy"))
		assertSuccess(response, "not_enough_balance", t)
	})
	t.Run("cancel order", func(t *testing.T) {
		response, _ := client.CancelOrder(
			Id("M103975"),
		)
		assertSuccess(response, "invalid_request", t)
	})
	t.Run("create instant", func(t *testing.T) {
		response, _ := client.CreateInstant(
			Market("ETHCLP"),
			Type("buy"),
			Amount("10"),
		)
		assertSuccess(response, "not_enough_balance", t)
	})
	t.Run("request deposit", func(t *testing.T) {
		response, _ := client.RequestDeposit(
			BankAccount("213213"),
			Amount("10234"),
		)
		assertSuccess(response, "BankAccount matching query does not exist.", t)
	})
	t.Run("request withdrawal", func(t *testing.T) {
		response, _ := client.RequestWithdrawal(
			Amount("10234"),
			BankAccount("213213"),
		)
		assertSuccess(response, "BankAccount matching query does not exist.", t)
	})
	t.Run("transfer", func(t *testing.T) {
		response, _ := client.Transfer(
			Address("GDMXNQBJMS3FYI4PFSYCCB4"),
			Amount("1200"),
			Currency("XLM"),
			Memo("162354"),
		)
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
		response, _ := client.NewOrder(
			CallbackUrl(""),
			ErrorUrl(""),
			ExternalId("ABC123"),
			PaymentReceiver("user@email.com"),
			SuccessUrl(""),
			ToReceive("3000"),
			ToReceiveCurrency("CLP"),
			RefundMail("refund@mail.com"),
		)
		assertSuccess(response, "invalid_request", t)
	})
	t.Run("create wallet", func(t *testing.T) {
		response, _ := client.CreateWallet(
			Id("P2023132"),
			Token("xToY232aheSt8F"),
			Wallet("ETH"),
		)
		assertSuccess(response, "payment_does_not_exist", t)
	})
	t.Run("payment orders", func(t *testing.T) {
		response, _ := client.PaymentOrders(
			StartDate("01/03/2018"),
			EndDate("08/03/2018"),
		)
		assertSuccess(response, "", t)
	})
	t.Run("payment status", func(t *testing.T) {
		response, _ := client.PaymentStatus(
			Id("P13433"),
		)
		assertSuccess(response, "invalid_scope", t)
	})
}

