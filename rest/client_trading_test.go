package rest

import (
	"context"
	"testing"

	"github.com/cryptomkt/go-api/args"
	"github.com/cryptomkt/go-api/models"
)

func TestGetTradingBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getTradingBalance(context.Background())
	if err != nil {
		t.Error(err)
	}
	for _, balance := range result {
		if err = checkBalance(&balance); err != nil {
			t.Error(err)
		}
	}
}

func TestGetActiveOrders(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getActiveOrders(context.Background())
	if err != nil {
		t.Error(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Error(err)
		}
	}
}

func TestCreateOrder(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.createOrder(context.Background(), args.Symbol("EOSETH"), args.Side("sell"), args.Quantity("0.01"), args.Price("9999"))
	if err != nil {
		t.Error(err)
	}
	if err = checkOrder(result); err != nil {
		t.Error(err)
	}
}

func TestCancelAllOrders(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.cancelAllOrders(context.Background(), args.Symbol("ETHBTC"))
	if err != nil {
		t.Error(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Error(err)
		}
	}
}

func TestOrderFlow(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	order, err := client.createOrder(context.Background(), args.Symbol("EOSETH"), args.Side("sell"), args.Quantity("0.01"), args.Price("9999"))
	if err != nil {
		t.Error(err)
	}
	if err = checkOrder(order); err != nil {
		t.Error(err)
	}
	order, err = client.getActiveOrder(context.Background(), args.ClientOrderID(order.ClientOrderID))
	if err != nil {
		t.Error(err)
	}
	if err = checkOrder(order); err != nil {
		t.Error(err)
	}
	order, err = client.cancelOrder(context.Background(), args.ClientOrderID(order.ClientOrderID))
	if err != nil {
		t.Error(err)
	}
	if err = checkOrder(order); err != nil {
		t.Error(err)
	}
	if order.Status != models.OrderStatusCanceled {
		t.Error("order not cancelled")
	}
}

func TestGetTradingFee(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getTradingFee(context.Background(), args.Symbol("EOSETH"))
	if err != nil {
		t.Error(err)
	}
	if result.ProvideLiquidityRate == "" || result.TakeLiquidityRate == "" {
		t.Error("fee should be defined")
	}
}
