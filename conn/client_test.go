package conn

import (
	"bufio"
	"encoding/json"
	"github.com/cryptomkt/cryptomkt-go/args"
	"log"
	"os"
	"testing"
	"time"
)

// newDebugClient initializes a client to run request,
// so its recomendable to not have money in the account for testing
// reads the first two lines of a file, the first one is the api key,
// the second one is the api secret
// **WARNING** DO NOT SHARE YOUR KEYS, KEEP THEM OUT OF THE REPOSITORY
// (with .gitignore for example)
func newDebugClient(filePath string) (*Client, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	apiKey := scanner.Text()
	scanner.Scan()
	apiSecret := scanner.Text()

	client := NewClient(apiKey, apiSecret)
	return client, nil
}

var keyfile = "../keys.txt"

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
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}

	//test get methods
	time.Sleep(6 * time.Second)
	t.Run("account", func(t *testing.T) {
		response, _ := client.Account()
		assertSuccess(response, "", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("wallet", func(t *testing.T) {
		response, _ := client.Balance()
		assertSuccess(response, "", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("transactions", func(t *testing.T) {
		response, _ := client.Transactions(
			args.Currency("ETH"))
		assertSuccess(response, "", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("active orders=1", func(t *testing.T) {
		response, _ := client.ActiveOrders(
			args.Market("ETHCLP"),
			args.Page(0))
		assertSuccess(response, "", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("active orders=2", func(t *testing.T) {
		response, _ := client.ActiveOrders(
			args.Market("ETHARS"),
			args.Page(1))
		assertSuccess(response, "", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("order status", func(t *testing.T) {
		response, _ := client.OrderStatus(
			args.Id("M103975"))
		assertSuccess(response, "invalid_scope", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("instant", func(t *testing.T) {
		response, _ := client.Instant(
			args.Market("ETHCLP"),
			args.Type("sell"),
			args.Amount("159"))
		assertSuccess(response, "", t)
	})

	time.Sleep(6 * time.Second)
	t.Run("executed orders", func(t *testing.T) {
		response, _ := client.ExecutedOrders(
			args.Market("ETHCLP"),
			args.Page(0))
		assertSuccess(response, "invalid_type", t)
	})

	//test post methods
	time.Sleep(6 * time.Second)
	t.Run("create order", func(t *testing.T) {
		response, _ := client.CreateOrder(
			args.Amount("0.3"),
			args.Market("ETHCLP"),
			args.Price("1000"),
			args.Type("buy"))
		assertSuccess(response, "not_enough_balance", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("cancel order", func(t *testing.T) {
		response, _ := client.CancelOrder(
			args.Id("M103975"),
		)
		assertSuccess(response, "invalid_request", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("create instant", func(t *testing.T) {
		response, _ := client.CreateInstant(
			args.Market("ETHCLP"),
			args.Type("buy"),
			args.Amount("10"),
		)
		assertSuccess(response, "not_enough_balance", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("request deposit", func(t *testing.T) {
		response, _ := client.RequestDeposit(
			args.BankAccount("213213"),
			args.Amount("10234"),
		)
		assertSuccess(response, "BankAccount matching query does not exist.", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("request withdrawal", func(t *testing.T) {
		response, _ := client.RequestWithdrawal(
			args.Amount("10234"),
			args.BankAccount("213213"),
		)
		assertSuccess(response, "BankAccount matching query does not exist.", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("transfer", func(t *testing.T) {
		response, _ := client.Transfer(
			args.Address("GDMXNQBJMS3FYI4PFSYCCB4"),
			args.Amount("1200"),
			args.Currency("XLM"),
			args.Memo("162354"),
		)
		assertSuccess(response, "max_limit_exceeded", t)
	})
}

func TestCryptoCompra(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}

	//test every endpoint
	time.Sleep(6 * time.Second)
	t.Run("new order", func(t *testing.T) {
		response, _ := client.NewOrder(
			args.CallbackUrl(""),
			args.ErrorUrl(""),
			args.ExternalId("ABC123"),
			args.PaymentReceiver("user@email.com"),
			args.SuccessUrl(""),
			args.ToReceive(3000),
			args.ToReceiveCurrency("CLP"),
			args.RefundEmail("refund@mail.com"),
		)
		assertSuccess(response, "invalid_request", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("create wallet", func(t *testing.T) {
		response, _ := client.CreateWallet(
			args.Id("P2023132"),
			args.Token("xToY232aheSt8F"),
			args.Wallet("ETH"),
		)
		assertSuccess(response, "payment_does_not_exist", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("payment orders", func(t *testing.T) {
		response, _ := client.PaymentOrders(
			args.StartDate("01/03/2018"),
			args.EndDate("08/03/2018"),
		)
		assertSuccess(response, "", t)
	})
	time.Sleep(6 * time.Second)
	t.Run("payment status", func(t *testing.T) {
		response, _ := client.PaymentStatus(
			args.Id("P13433"),
		)
		assertSuccess(response, "invalid_scope", t)
	})
}
