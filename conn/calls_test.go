package conn

import (
	"fmt"
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
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
		account, err := client.GetAccount()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(account)
	})
}

func TestGetWallet(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
		_, err := client.GetBalance()
		if err != nil {
			t.Error(err)
		}
	})
}
func TestOrderFlow(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(delay)
	order, err := client.CreateOrder(
		args.Market("XLMCLP"),
		args.Type("sell"),
		args.Price("100"),
		args.Amount("1"),
	)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Println(order.String())
	time.Sleep(delay)
	order, _ = order.Refresh()
	fmt.Println(order.String())
	time.Sleep(delay)
	order, _ = order.Close()
	fmt.Println(order.String())
}

func TestOrderListFlow(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, v := range []string{"1", "2", "3", "4", "5"} {
		time.Sleep(delay)
		_, err := client.CreateOrder(
			args.Market("XLMCLP"),
			args.Type("sell"),
			args.Price("100"+v),
			args.Amount("1"),
		)
		if err != nil {
			t.Errorf("error: %v", err)
		}
	}
	oList, err := client.GetActiveOrders(args.Market("XLMCLP"))
	fmt.Println(oList)
	time.Sleep(delay)
	oList.Refresh()
	fmt.Println(oList)
	time.Sleep(delay)
	oList.Close()
	fmt.Println(oList)
}

func TestTransactions(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Run("xlm", func(t *testing.T) {
		time.Sleep(delay)
		transactions, err := client.GetTransactions(
			args.Currency("XLM"))
		if err != nil {
			t.Error(err)
		}
		fmt.Println(transactions)
	})
	t.Run("missing=Currency", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("xlm", func(t *testing.T) {
		time.Sleep(delay)
		activeOrders, err := client.GetActiveOrders(
			args.Market("XLMCLP"),
			args.Page(0))
		if err != nil {
			t.Error(err)
		}
		fmt.Println(activeOrders)
	})
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
		_, err := client.GetActiveOrders(
			args.Market("ETHARS"),
			args.Page(1))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("missing=market", func(t *testing.T) {
		time.Sleep(delay)
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

	t.Run("0", func(t *testing.T) {
		time.Sleep(delay)
		orderList, err := client.GetExecutedOrders(
			args.Market("XLMCLP"),
			args.Page(0))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_type") {
				t.Error(err)
			}
		}
		fmt.Println(orderList)
	})
	t.Run("missing=market", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("xlm", func(t *testing.T) {
		time.Sleep(delay)
		order, err := client.CreateOrder(
			args.Amount("1"),
			args.Market("XMLCLP"),
			args.Price("100"),
			args.Type("sell"))
		if err != nil {
			if !strings.Contains(err.Error(), "not_enough_balance") {
				t.Error(err)
			}
		}
		fmt.Println(order.String())
	})
}

func TestOrderStatus(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
		_, err := client.GetOrderStatus(
			args.Id("M103975"))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_scope") {
				t.Error(err)
			}
		}
	})
	t.Run("missing=id", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("cancel order", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("create instant", func(t *testing.T) {
		time.Sleep(delay)
		err := client.CreateInstant(
			args.Market("XLMCLP"),
			args.Type("sell"),
			args.Amount("1"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "not_enough_balance") {
				fmt.Println(err)
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
	t.Run("0", func(t *testing.T) {
		time.Sleep(delay)
		instant, err := client.GetInstant(
			args.Market("XLMCLP"),
			args.Type("buy"),
			args.Amount("100"))
		if err != nil {
			if !strings.Contains(err.Error(), "invalid_request") {
				t.Error(err)
			}
			t.Error(err)
		}
		fmt.Println(instant)
	})

	t.Run("missing=market and type", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("unsupported=type", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
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
	t.Run("", func(t *testing.T) {
		time.Sleep(delay)
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
