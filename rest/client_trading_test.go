package rest

import (
	"context"
	"testing"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

func TestGetTradingBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetTradingBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, balance := range result {
		if err = checkBalance(&balance); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetActiveOrders(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetActiveOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreateOrder(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.CreateOrder(context.Background(), args.Symbol("EOSETH"), args.Side(args.SideTypeSell), args.Quantity("0.01"), args.Price("8999"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(result); err != nil {
		t.Fatal(err)
	}
}

func TestCancelAllOrders(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.CancelAllOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOrderFlow(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	order, err := client.CreateOrder(context.Background(), args.Symbol("EOSETH"), args.Side("sell"), args.Quantity("0.01"), args.Price("9999"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
	order, err = client.GetActiveOrder(context.Background(), args.ClientOrderID(order.ClientOrderID))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
	order, err = client.CancelOrder(context.Background(), args.ClientOrderID(order.ClientOrderID))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
	if order.Status != models.OrderStatusCanceled {
		t.Fatal("order not Cancelled")
	}
}

func TestGetTradingFee(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetTradingFee(context.Background(), args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	if result.ProvideLiquidityRate == "" || result.TakeLiquidityRate == "" {
		t.Fatal("fee should be defined")
	}
}
