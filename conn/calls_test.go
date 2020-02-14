package conn

import (
	"github.com/cryptomkt/cryptomkt-go/args"
	"strings"
	"testing"
	"time"
)

func TestTransactions(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("missing=Currency", func(t *testing.T) {
		_, err := client.Transactions()
		if err == nil {
			t.Errorf("no error rised, should rise missing Currency arg")
		}
		if !strings.Contains(err.Error(), "currency") {
			t.Errorf("should advise the lack of currency argument, got %s", err)
		}
	})
}

func TestActiveOrders(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("missing=market", func(t *testing.T) {
		_, err := client.ActiveOrders()
		if err == nil {
			t.Errorf("no error rised, should rise missing market arg")
		}
		if !strings.Contains(err.Error(), "market") {
			t.Errorf("should advise the lack of market argument, got %s", err)
		}
	})
}

func TestExecutedOrders(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("missing=market", func(t *testing.T) {
		_, err := client.ExecutedOrders()
		if err == nil {
			t.Errorf("no error rised, should rise missing market arg")
		}
		if !strings.Contains(err.Error(), "market") {
			t.Errorf("should advise the lack of market argument, got %s", err)
		}
	})
}

func TestOrderStatus(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("missing=id", func(t *testing.T) {
		_, err := client.OrderStatus()
		if err == nil {
			t.Errorf("no error rised, should rise missing id arg")
		}
		if !strings.Contains(err.Error(), "id") {
			t.Errorf("should advise the lack of id argument, got %s", err)
		}
	})
}

func TestInstant(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(6 * time.Second)
	t.Run("missing=market and type", func(t *testing.T) {
		_, err := client.Instant(
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
		_, err := client.Instant(
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
