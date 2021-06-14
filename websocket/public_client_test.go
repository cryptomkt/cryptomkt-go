package websocket

import (
	"context"
	"fmt"
	"testing"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestGetCurrencies(t *testing.T) {
	client, err := NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	bg := context.Background()
	result, err := client.GetCurrencies(bg)
	if err != nil {
		t.Fatal(err)
	}
	for _, currency := range result {
		err := checkCurrency(&currency)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetSymbols(t *testing.T) {
	client, err := NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	errCh := make(chan error)
	done := make(chan struct{})
	bg := context.Background()
	go func() {
		defer close(errCh)
		defer close(done)
		symbols, err := client.GetSymbols(bg)
		if err != nil {
			errCh <- err
			return
		}
		if len(symbols) == 0 {
			errCh <- fmt.Errorf("no symbols")
		}
		for _, symbol := range symbols {
			err := checkSymbol(&symbol)
			if err != nil {
				errCh <- err
			}
		}
		done <- struct{}{}
	}()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-done:
	}
}

func TestGetSymbol(t *testing.T) {
	client, err := NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	bg := context.Background()
	symbol, err := client.GetSymbol(bg, args.Symbol("ETHBTC"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkSymbol(symbol); err != nil {
		t.Fatal(err)
	}
}

func TestGetTrades(t *testing.T) {
	client, err := NewPublicClient()
	if err != nil {
		t.Fatal(err)
	}
	bg := context.Background()
	result, err := client.GetTrades(bg, args.Symbol("EOSETH"), args.Limit(2))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) > 2 {
		t.Fatal("to many trades")
	}
	for _, trade := range result {
		if err = checkPublicTrade(&trade); err != nil {
			t.Fatal(err)
		}
	}
}
