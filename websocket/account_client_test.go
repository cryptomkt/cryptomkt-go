package websocket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"

	"github.com/cryptomarket/cryptomarket-go/rest"
)

func TestGetAccountBalance(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := NewAccountClient(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	balances, err := client.GetAccountBalance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(balances) == 0 {
		t.Fatal("no balances")
	}
	for _, balance := range balances {
		if err = checkBalance(&balance); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFindTransactions(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := NewAccountClient(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	transactions, err := client.FindTransactions(context.Background(), args.Limit(5))
	if err != nil {
		t.Fatal(err)
	}
	if len(transactions) != 5 {
		t.Fatal("not the right numbre of transactions")
	}
	for _, transaction := range transactions {
		if err = checkTransaction(&transaction); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadTransactions(t *testing.T) {
	apiKeys := LoadKeys()
	client, err := NewAccountClient(apiKeys.APIKey, apiKeys.APISecret)
	if err != nil {
		t.Fatal(err)
	}
	transactions, err := client.LoadTransactions(context.Background(), args.Limit(3), args.Sort(args.SortTypeASC))
	if err != nil {
		t.Fatal(err)
	}
	if len(transactions) != 3 {
		t.Fatal("not the right numbre of transactions")
	}
	for _, transaction := range transactions {
		if err = checkTransaction(&transaction); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTransactionsSubscription(t *testing.T) {
	bg := context.Background()
	apiKeys := LoadKeys()
	client, _ := NewAccountClient(apiKeys.APIKey, apiKeys.APISecret)
	restClient := rest.NewClient(apiKeys.APIKey, apiKeys.APISecret)
	feedCh, err := client.SubscribeToTransactions()
	if err != nil {
		t.Fatal(err)
	}
	innerErrCh := make(chan error)
	go func() {
		defer close(innerErrCh)
		for transaction := range feedCh {
			if err := checkTransaction(&transaction); err != nil {
				innerErrCh <- err
				return
			}
		}
	}()
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
	restClient.TransferMoneyFromTradingToAccountBalance(bg, args.Amount("0.2"), args.Currency("EOS"))
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
	restClient.TransferMoneyFromAccountToTradingBalance(bg, args.Amount("0.2"), args.Currency("EOS"))
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
	client.UnsubscribeToTransactions()
	<-time.After(5 * time.Second)
}

func TestBalanceSubscription(t *testing.T) {
	bg := context.Background()
	apiKeys := LoadKeys()
	client, _ := NewAccountClient(apiKeys.APIKey, apiKeys.APISecret)
	restClient := rest.NewClient(apiKeys.APIKey, apiKeys.APISecret)
	feedCh, err := client.SubscribeToBalance()
	if err != nil {
		t.Fatal(err)
	}
	innerErrCh := make(chan error)
	go func() {
		defer close(innerErrCh)
		for balances := range feedCh {
			fmt.Println(balances)
			for _, balance := range balances {
				if err := checkBalance(&balance); err != nil {
					innerErrCh <- err
					return
				}

			}
		}
	}()
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
	restClient.TransferMoneyFromTradingToAccountBalance(bg, args.Amount("0.2"), args.Currency("EOS"))
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
	restClient.TransferMoneyFromAccountToTradingBalance(bg, args.Amount("0.2"), args.Currency("EOS"))
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
	client.UnsubscribeToBalance()
	<-time.After(5 * time.Second)
}
