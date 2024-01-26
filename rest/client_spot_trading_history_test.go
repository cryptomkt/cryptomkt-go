package rest

import (
	"context"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestOrderHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	result, err := client.GetSpotOrdersHistory(context.Background(), args.Limit(200))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkOrder, result); err != nil {
		t.Fatal(err)
	}
}

func TestContextDeadline(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	deadline := time.Now().Add(1 * time.Millisecond)
	ctx, cancelFunc := context.WithDeadline(context.Background(), deadline)
	defer cancelFunc()
	result, err := client.GetSpotOrdersHistory(ctx, args.Symbol("EOSETH"))
	if err != nil {
		return
	}
	t.Error(result)
}

func TestGetTradesHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	result, err := client.GetSpotTradesHistory(context.Background(), args.Limit(199))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkTrade, result); err != nil {
		t.Fatal(err)
	}
}
