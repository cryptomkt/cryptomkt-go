package rest

import (
	"context"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestOrderHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetOrderHistory(context.Background(), args.Limit(200))
	if err != nil {
		t.Fatal(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetOldOrder(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetOrders(context.Background(), args.ClientOrderID("0a027a4ae8f44934519b211cf0e8e52e"))
	if err != nil {
		t.Fatal(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Fatal(err)
		}
	}
}

func TestContextDeadline(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	d := time.Now().Add(1 * time.Millisecond)
	ctx, cancelFunc := context.WithDeadline(context.Background(), d)
	defer cancelFunc()
	result, err := client.GetOrders(ctx, args.ClientOrderID("0a027a4ae8f44934519b211cf0e8e52e"))
	if err != nil {
		return
	}
	t.Error(result)
}

func TestGetTradesHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetTradeHistory(context.Background(), args.Limit(199))
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range result {
		if err = checkTrade(&trade); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetTradesByOrderID(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetTradesByOrderID(context.Background(), args.OrderID(337789478188))
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range result {
		if err = checkTrade(&trade); err != nil {
			t.Fatal(err)
		}
	}
}
