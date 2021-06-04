package websocket

import (
	"context"
	"testing"
)

func TestAuth(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := New(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	err = client.authenticate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTradingBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := New(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	err = client.authenticate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.getTradingBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
