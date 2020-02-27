package conn

import (
	"strings"
	"testing"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
)

func TestMarket(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	if _, err := client.GetMarkets(); err != nil {
		t.Errorf("Market Test failed because %s", err)
	}
}
func TestTicker(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	if _, err := client.GetTicker(); err != nil {
		t.Errorf("Ticker with no optional arguments failed because %s", err)
	}
	if _, err := client.GetTicker(argus[0]); err != nil {
		t.Errorf("Ticker with market argument failed because %s", err)
	}
}
func TestGetAccount(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetAccount()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestGetWallet(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetBalance()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestTransactions(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetTransactions(
			args.Currency("ETH"))
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("missing=Currency", func(t *testing.T) {
		_, err := client.GetTransactions()
		if err == nil {
			t.Errorf("no error rised, should rise missing Currency arg")
		}
		if !strings.Contains(err.Error(), "currency") {
			t.Errorf("should advise the lack of currency argument, got %s", err)
		}
	})
}

func TestActiveOrders(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetActiveOrders(
			args.Market("ETHCLP"),
			args.Page(0))
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetActiveOrders(
			args.Market("ETHARS"),
			args.Page(1))
		if err != nil {
			t.Error(err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("missing=market", func(t *testing.T) {
		_, err := client.GetActiveOrders()
		if err == nil {
			t.Errorf("no error rised, should rise missing market arg")
		}
		if !strings.Contains(err.Error(), "market") {
			t.Errorf("should advise the lack of market argument, got %s", err)
		}
	})
}

func TestExecutedOrders(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}

	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetExecutedOrders(
			args.Market("ETHCLP"),
			args.Page(0))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_type") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("missing=market", func(t *testing.T) {
		_, err := client.GetExecutedOrders()
		if err == nil {
			t.Errorf("no error rised, should rise missing market arg")
		}
		if !strings.Contains(err.Error(), "market") {
			t.Errorf("should advise the lack of market argument, got %s", err)
		}
	})
}

func TestCreateOrder(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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
}

func TestOrderStatus(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetOrderStatus(
			args.Id("M103975"))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_scope") {
				t.Error(err)
			}
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("missing=id", func(t *testing.T) {
		_, err := client.GetOrderStatus()
		if err == nil {
			t.Errorf("no error rised, should rise missing id arg")
		}
		if !strings.Contains(err.Error(), "id") {
			t.Errorf("should advise the lack of id argument, got %s", err)
		}
	})
}

func TestCancelOrder(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
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
}
func TestCreateInstant(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
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
}

func TestGetInstant(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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
	t.Run("missing=market and type", func(t *testing.T) {
		_, err := client.GetInstant(
			args.Amount("123"),
		)
		if err == nil {
			t.Errorf("no error rised, should rise missing id arg")
		}
		if !strings.Contains(err.Error(), "market") &&
			!strings.Contains(err.Error(), "type") {
			t.Errorf("should advise the lack of market and type arguments, got %s", err)
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("unsupported=type", func(t *testing.T) {
		_, err := client.GetInstant(
			args.Price("2020"),
			args.Amount("123"),
			args.Type("see"),
		)
		if err == nil {
			t.Errorf("no error rised, should rise missing id arg")
		}
		if !strings.Contains(err.Error(), "type") &&
			!strings.Contains(err.Error(), "buy") &&
			!strings.Contains(err.Error(), "sell") {
			t.Errorf("should advise about accepted types, buy and sell, got %s", err)
		}
	})
}

func TestRequestDeposit(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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
}

func TestRequestWithdrawal(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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
}

func TestTransfer(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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

func TestNewOrder(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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
}

func TestCreateWallet(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
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
}

func TestPaymentOrders(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetPaymentOrders(
			args.StartDate("01/03/2018"),
			args.EndDate("08/03/2018"),
		)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestGetPaymentStatus(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("", func(t *testing.T) {
		_, err := client.GetPaymentStatus(
			args.Id("P13433"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_request") {
				t.Error(err)
			}
		}
	})
}
