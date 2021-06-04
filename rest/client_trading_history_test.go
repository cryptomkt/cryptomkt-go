package rest

import (
	"context"
	"testing"
	"time"

	"github.com/cryptomkt/go-api/args"
)

func TestOrderHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getOrderHistory(context.Background(), args.Limit(200))
	if err != nil {
		t.Error(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Error(err)
		}
	}
}

func TestGetOldOrder(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getOrders(context.Background(), args.ClientOrderID("0a027a4ae8f44934519b211cf0e8e52e"))
	if err != nil {
		t.Error(err)
	}
	for _, order := range result {
		if err = checkOrder(&order); err != nil {
			t.Error(err)
		}
	}
}

func TestContextDeadline(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	d := time.Now().Add(1 * time.Millisecond)
	ctx, cancelFunc := context.WithDeadline(context.Background(), d)
	defer cancelFunc()
	result, err := client.getOrders(ctx, args.ClientOrderID("0a027a4ae8f44934519b211cf0e8e52e"))
	if err != nil {
		return
	}
	t.Error(result)
}

func TestGetTradesHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getTradeHistory(context.Background(), args.Limit(199))
	if err != nil {
		t.Error(err)
	}
	for _, trade := range result {
		if err = checkTrade(&trade); err != nil {
			t.Error(err)
		}
	}
}

func TestGetTradesByOrder(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getTradesByOrder(context.Background(), args.OrderID(337789478188))
	if err != nil {
		t.Error(err)
	}
	for _, trade := range result {
		if err = checkTrade(&trade); err != nil {
			t.Error(err)
		}
	}
}
