package websocket

import (
	"context"
	"testing"

	"github.com/cryptomkt/go-api/args"
)

func TestGetCurrencies(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	result, err := client.getCurrencies(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, currency := range result {
		if currency.ID == "" {
			t.Fatal("currency should have id")
		}
		if currency.FullName == "" {
			t.Fatal("currency should have full name")
		}
	}
}


func TestGetChanCurrencies(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	resultCh, errCh := client.getChanCurrencies(ctx)
	select {
	case err:=<-errCh:
		t.Fatal(err)
	case result:=<-resultCh:
		for _, currency := range result {
			if currency.ID == "" {
				t.Fatal("currency should have id")
			}
			if currency.FullName == "" {
				t.Fatal("currency should have full name")
			}
		}
	}
}

func TestGetSymbols(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	result, err := client.getSymbols(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range result {
		if symbol.ID == "" {
			t.Fatal("symbol should have id")
		}
		if symbol.BaseCurrency == "" {
			t.Fatal("symbol should have Base Currency")
		}
		if symbol.QuoteCurrency == "" {
			t.Fatal("symbol should have Quote Currency")
		}
	}
}

func TestGetSymbol(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	symbol, err := client.getSymbol(ctx, args.Symbol("ETHBTC"))
	if err != nil {
		t.Fatal(err)
	}
	if symbol.ID == "" {
		t.Fatal("symbol should have id")
	}
	if symbol.BaseCurrency == "" {
		t.Fatal("symbol should have Base Currency")
	}
	if symbol.QuoteCurrency == "" {
		t.Fatal("symbol should have Quote Currency")
	}
}

func TestGetTrades(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	result, symbol, err := client.getTrades(ctx, args.Symbol("EOSETH"), args.Limit(2))
	if err != nil {
		t.Fatal(err)
	}
	if symbol != "EOSETH" {
		t.Fatal("should be same symbol")
	}
	for _, trade := range result {
		if trade.ID == 0 {
			t.Fatal("should not be default value for ID")
		}
		if trade.Price == "" {
			t.Fatal("should have price")
		}
		if trade.Quantity == "" {
			t.Fatal("should have quantity")
		}
		if !(trade.Side == "sell" || trade.Side == "buy") {
			t.Fatal("should be side sell or side buy")
		}
	}
}
