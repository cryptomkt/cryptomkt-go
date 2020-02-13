package client

import (
	"fmt"
	"testing"
	"time"
	"strings"
	"github.com/cryptomkt/cryptomkt-go/args"
)


func TestTransactions(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(10 * time.Second)
	t.Run("missing=Currency", func(t *testing.T) {
		_, err := client.Transactions()
		if err == nil {
			t.Errorf("no error rised, should rise missing Currency arg")
		}
		if !strings.Contains(fmt.Sprintf("%v",err), "currency") {
			t.Errorf("should advise the lack of currency argument, got%s", err)
		}
	})
	/*
		fails, currency has no restrictions, although it can be added
		to be CLP, ARS, BRL, MXN, EUR, ETH, XLM or BTC.
		t.Run("unsupported=Currency", func(t *testing.T) {
			_, err := client.Transactions(
				Currency("USD"),
			)
			if err == nil {
				t.Errorf("no error rised, should rise not supported currency")
			}
			//assert some kind of error
		})
	*/
}

func TestActiveOrders(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(10 * time.Second)
	t.Run("missing=market", func(t *testing.T) {
		_, err := client.ActiveOrders()
		if err == nil {
			t.Errorf("no error rised, shold rise missing market arg")
		}
		if !strings.Contains(fmt.Sprintf("%v",err), "market") {
			t.Errorf("should advise the lack of market argument, got%s", err)
		}
	})
}

func TestExecutedOrders(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(10 * time.Second)
	t.Run("missing=market", func(t *testing.T) {
		_, err := client.ExecutedOrders()
		if err == nil {
			t.Errorf("no error rised, shold rise missing market arg")
		}
		if !strings.Contains(fmt.Sprintf("%v",err), "market") {
			t.Errorf("should advise the lack of market argument, got%s", err)
		}
	})
}

func TestOrderStatus(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(10 * time.Second)
	t.Run("missing=id", func(t *testing.T) {
		_, err := client.OrderStatus()
		if err == nil {
			t.Errorf("no error rised, shold rise missing id arg")
		}
		if !strings.Contains(fmt.Sprintf("%v",err), "id") {
			t.Errorf("should advise the lack of id argument, got %s", err)
		}
	})
}

func TestInstant(t *testing.T) {
	client, err := newDebugClient(keyfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	time.Sleep(10 * time.Second)
	t.Run("missing=market and type", func(t *testing.T) {
		_, err := client.Instant(
			args.Amount("123"),
		)
		if err == nil {
			t.Errorf("no error rised, shold rise missing id arg")
		}
		if !strings.Contains(fmt.Sprintf("%v",err), "market") &&
		   !strings.Contains(fmt.Sprintf("%v",err), "type") {
			t.Errorf("should advise the lack of market and type arguments, got %s", err)
		}
	})
	time.Sleep(10 * time.Second)
	t.Run("unsupported=type", func(t *testing.T) {
		_, err := client.Instant(
			args.Price("2020"),
			args.Amount("123"),
			args.Type("see"),
		)
		if err == nil {
			t.Errorf("no error rised, shold rise missing id arg")
		}
		if !strings.Contains(fmt.Sprintf("%v",err), "market") &&
		   !strings.Contains(fmt.Sprintf("%v",err), "type") {
			t.Errorf("should advise the lack of market and type arguments, got %s", err)
		}
	})
}