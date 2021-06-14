package websocket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestAuth(t *testing.T) {
	apiKeys := LoadKeys()
	if _, err := NewTradingClient(apiKeys.APIKey, apiKeys.APISecret); err != nil {
		t.Fatal(err)
	}
	if _, err := NewTradingClient("no apikey", "no apisecret"); err == nil {
		t.Fatal(err)
	}

}

func TestGetTradingBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := NewTradingClient(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	balances, err := client.GetTradingBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(balances) == 0 {
		t.Fatal("no balances")
	}
	for _, balance := range balances {
		if err = checkBalance(&balance); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOrderFlow(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := NewTradingClient(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	bg := context.Background()
	clientOrderID := fmt.Sprint(time.Now().Unix())
	_, err = client.CreateOrder(bg, args.Symbol("EOSETH"), args.Side(args.SideTypeSell), args.Price("1000"), args.Quantity("0.01"), args.ClientOrderID(clientOrderID))
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{}, 1)
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		defer close(done)
		report, err := client.ReplaceOrder(bg, args.ClientOrderID(clientOrderID), args.RequestClientID(clientOrderID+"new"), args.Quantity("0.02"), args.Price("2000"))
		if err != nil {
			errCh <- err
		}
		if err := checkReport(report); err != nil {
			errCh <- err
		}
		report, err = client.CancelOrder(bg, args.ClientOrderID(clientOrderID+"new"))
		if err != nil {
			errCh <- err
		}
		done <- struct{}{}
	}()
	if err = <-errCh; err != nil {
		t.Error(err)
	}
	<-done
}

func TestReportsSubscription(t *testing.T) {
	apiKeys := LoadKeys()
	client, _ := NewTradingClient(apiKeys.APIKey, apiKeys.APISecret)
	feedCh, err := client.SubscribeToReports()
	if err != nil {
		t.Fatal(err)
	}
	innerErrCh := make(chan error)
	go func() {
		defer close(innerErrCh)
		for report := range feedCh {
			if err := checkReport(&report); err != nil {
				innerErrCh <- err
				return
			}
		}
	}()
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(8 * time.Minute):
	}
}
