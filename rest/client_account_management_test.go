package rest

import (
	"context"
	"testing"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestGetAccountBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client := NewClient(apiKeys.APIKey, apiKeys.APISecret)
	result, err := client.GetAccountBalance(context.Background())
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
	result, err := client.GetDepositCryptoAddress(context.Background(), args.Currency("EOS"))
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
	result, err := client.CreateDepositCryptoAddress(context.Background(), args.Currency("EOS"))
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
	result, err := client.GetLast10DepositCryptoAddresses(context.Background(), args.Currency("EOS"))
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
	result, err := client.GetLast10UsedCryptoAddresses(context.Background(), args.Currency("EOS"))
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
	result, err := client.GetEstimatesWithdrawFee(context.Background(), args.Currency("EOS"), args.Amount("199"))
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

	result, err := client.TransferMoneyFromTradingToAccountBalance(bg, args.Currency(symbol), args.Amount(amount))
	if err != nil {
		t.Fatal(err)
	}
	if result.ID == "" {
		t.Fatal("should have only id defined")
	}

	result, err = client.TransferMoneyFromAccountToTradingBalance(bg, args.Currency(symbol), args.Amount(amount))
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
	result, err := client.GetTransactionHistory(context.Background(), args.Currency("EOS"))
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
	result, err := client.GetTransaction(context.Background(), args.ID(id))
	if err != nil {
		t.Error(err)
	}
	if err = checkTransaction(result); err != nil {
		t.Error(err)
	}
}
