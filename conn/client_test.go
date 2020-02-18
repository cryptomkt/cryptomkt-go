package conn

import (
	"bufio"
	"github.com/cryptomkt/cryptomkt-go/args"
	"log"
	"os"
	"strings"
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

var keysfile = "../keys.txt"

func TestAutenticated(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}

	//test get methods
	time.Sleep(6 * time.Second)
	t.Run("account", func(t *testing.T) {
		_, err := client.GetAccount()
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("wallet", func(t *testing.T) {
		_, err := client.GetBalance()
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("transactions", func(t *testing.T) {
		_, err := client.GetTransactions(
			args.Currency("ETH"))
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("active orders=1", func(t *testing.T) {
		_, err := client.GetActiveOrders(
			args.Market("ETHCLP"),
			args.Page(0))
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("active orders=2", func(t *testing.T) {
		_, err := client.GetActiveOrders(
			args.Market("ETHARS"),
			args.Page(1))
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("order status", func(t *testing.T) {
		_, err := client.GetOrderStatus(
			args.Id("M103975"))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_scope") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("instant", func(t *testing.T) {
		_, err := client.GetInstant(
			args.Market("ETHCLP"),
			args.Type("sell"),
			args.Amount("19000"))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_request") {
				t.Error(err)
			}
		}
	})

	time.Sleep(6 * time.Second)
	t.Run("executed orders", func(t *testing.T) {
		_, err := client.GetExecutedOrders(
			args.Market("ETHCLP"),
			args.Page(0))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_type") {
				t.Error(err)
			}
		}
	})

	//test post methods
	time.Sleep(6 * time.Second)
	t.Run("create order", func(t *testing.T) {
		_, err := client.CreateOrder(
			args.Amount("0.3"),
			args.Market("ETHCLP"),
			args.Price("1000"),
			args.Type("buy"))
		if err != nil {
			if !strings.Contains(err.Error(), "not_enough_balance") {
				t.Error(err)
			}
		}

	})
	time.Sleep(6 * time.Second)
	t.Run("cancel order", func(t *testing.T) {
		_, err := client.CancelOrder(
			args.Id("M103975"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_scope") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("create instant", func(t *testing.T) {
		err := client.CreateInstant(
			args.Market("ETHCLP"),
			args.Type("buy"),
			args.Amount("10000"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "not_enough_balance") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("request deposit", func(t *testing.T) {
		err := client.RequestDeposit(
			args.BankAccount("213213"),
			args.Amount("10234"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "Bank account does not exist") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("request withdrawal", func(t *testing.T) {
		err := client.RequestWithdrawal(
			args.Amount("10234"),
			args.BankAccount("213213"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "Bank account does not exist") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("transfer", func(t *testing.T) {
		err := client.Transfer(
			args.Address("GDMXNQBJMS3FYI4PFSYCCB4"),
			args.Amount("1200"),
			args.Currency("XLM"),
			args.Memo("162354"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "max_limit_exceeded") {
				t.Error(err)
			}
		}
	})
}

func TestCryptoCompra(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}

	//test every endpoint
	time.Sleep(6 * time.Second)
	t.Run("new order", func(t *testing.T) {
		_, err := client.NewOrder(
			args.CallbackUrl(""),
			args.ErrorUrl(""),
			args.ExternalId("ABC123"),
			args.PaymentReceiver("user@email.com"),
			args.SuccessUrl(""),
			args.ToReceive(3000),
			args.ToReceiveCurrency("CLP"),
			args.RefundEmail("refund@mail.com"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "temporarily disabled") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("create wallet", func(t *testing.T) {
		_, err := client.CreateWallet(
			args.Id("P2023132"),
			args.Token("xToY232aheSt8F"),
			args.Wallet("ETH"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "Petición inválida") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("payment orders", func(t *testing.T) {
		_, err := client.PaymentOrders(
			args.StartDate("01/03/2018"),
			args.EndDate("08/03/2018"),
		)
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("payment status", func(t *testing.T) {
		_, err := client.PaymentStatus(
			args.Id("P13433"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_request") {
				t.Error(err)
			}
		}
	})
}
