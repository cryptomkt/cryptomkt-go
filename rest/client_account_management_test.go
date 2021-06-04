package rest

import (
	"context"
	"testing"

	"github.com/cryptomkt/go-api/args"
)

func TestGetAccountBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getAccountBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("should have results")
	}
	for _, balance := range result {
		if err = checkBalance(&balance); err != nil {
			t.Error(err)
		}
	}
}

func TestGetDepositCryptoAddress(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getDepositCryptoAddress(context.Background(), args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Address == "" {
		t.Fatal("should have an address")
	}
}

func TestCreateDepositCryptoAddress(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.createDepositCryptoAddress(context.Background(), args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	if result.Address == "" {
		t.Fatal("should have an address")
	}
}

func TestGetLast10DepositCryptoAddress(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getLast10DepositCryptoAddresses(context.Background(), args.Currency("EOS"))
	if err != nil {
		t.Error(err)
	}
	// if len(result) == 0 is ok
	for _, addr := range result {
		if addr.Address == "" {
			t.Fatal("should have an address")
		}
	}
}

func TestGetLast10UsedCryptoAddresses(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getLast10UsedCryptoAddresses(context.Background(), args.Currency("EOS"))
	if err != nil {
		t.Error(err)
	}
	// if len(result) == 0 is ok
	for _, addr := range result {
		if addr.Address == "" {
			t.Fatal("should have an address")
		}
	}
}

func TestGetEstimatesWithdrawFee(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getEstimatesWithdrawFee(context.Background(), args.Currency("EOS"), args.Amount("199"))
	if err != nil {
		t.Fatal(err)
	}
	if result == "" {
		t.Fatal("should have a result")
	}
}

func TestTransferBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	bg := context.Background()
	symbol := "EOS"
	amount := "0.01"

	result, err := client.transferFromTradingToAccountBalance(bg, args.Currency(symbol), args.Amount(amount))
	if err != nil {
		t.Fatal(err)
	}
	if result.ID == "" {
		t.Fatal("should have only id defined")
	}

	result, err = client.transferFromAccountToTradingBalance(bg, args.Currency(symbol), args.Amount(amount))
	if err != nil {
		t.Error(err)
	}
	if result.ID == "" {
		t.Fatal("should have only id defined")
	}
}

func TestGetTransactionHistory(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.getTransactionHistory(context.Background(), args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("should have transactions")
	}
	for _, transaction := range result {
		if err = checkTransaction(&transaction); err != nil {
			t.Error(err)
		}
	}
}

func TestGetTransaction(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	id := "7bb7abb6-c748-40ab-8070-9a7e4460fbde"
	result, err := client.getTransaction(context.Background(), args.ID(id))
	if err != nil {
		t.Error(err)
	}
	if err = checkTransaction(result); err != nil {
		t.Error(err)
	}
}
