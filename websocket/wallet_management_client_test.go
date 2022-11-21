package websocket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"

	"github.com/cryptomarket/cryptomarket-go/rest"
)

func loadWalletClient() (*WalletManagementClient, error) {
	apiKeys := LoadKeys()
	return NewWalletManagementClient(apiKeys.APIKey, apiKeys.APISecret, 15_000)
}

func TestWalletBalances(t *testing.T) {
	client, err := loadWalletClient()
	if err != nil {
		t.Fatal(err)
	}
	balances, err := client.GetWalletBalances(context.Background())
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

func TestWalletBalance(t *testing.T) {
	client, err := loadWalletClient()
	if err != nil {
		t.Fatal(err)
	}
	balance, err := client.GetWalletBalanceOfCurrency(
		context.Background(),
		args.Currency("CRO"),
	)
	if err != nil {
		t.Fatal(err)
	}
	balance.Currency = "CRO"
	if err = checkBalance(balance); err != nil {
		t.Fatal(balance)
	}
}

func TestGetTransactions(t *testing.T) {
	client, err := loadWalletClient()
	if err != nil {
		t.Fatal(err)
	}
	transactions, err := client.GetTransactions(
		context.Background(),
		args.TransactionType(args.TransactionTypeDeposit),
	)
	if err != nil {
		t.Fatal(err)
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
	client, _ := NewWalletManagementClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	restClient := rest.NewClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	notificationCh, err := client.SubscribeToTransactions()
	if err != nil {
		t.Fatal(err)
	}
	saver := runSaver()
	go func() {
		defer saver.close()
		for notification := range notificationCh {
			transaction := notification.Data
			saver.strSaveCh() <- fmt.Sprint(notification)
			if err := checkTransaction(&transaction); err != nil {
				saver.errSaveCh() <- err
				return
			}
		}
	}()
	<-time.After(5 * time.Second)
	restClient.TransferBetweenWalletAndExchange(
		bg,
		args.Amount("0.2"),
		args.Currency("EOS"),
		args.Source(args.AccountWallet),
		args.Destination(args.AccountSpot),
	)
	<-time.After(5 * time.Second)
	restClient.TransferBetweenWalletAndExchange(
		bg,
		args.Amount("0.2"),
		args.Currency("EOS"),
		args.Source(args.AccountSpot),
		args.Destination(args.AccountWallet),
	)
	<-time.After(5 * time.Second)
	err = client.UnsubscribeToTransactions()
	if err != nil {
		t.Fatal(err)
	}
	saver.printSavedStrings()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fatal()
	}
	<-time.After(5 * time.Second)
}

func TestBalanceSubscription(t *testing.T) {
	bg := context.Background()
	apiKeys := LoadKeys()
	client, _ := NewWalletManagementClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	restClient := rest.NewClient(apiKeys.APIKey, apiKeys.APISecret, 0)
	notificationCh, err := client.SubscribeToWalletBalance()
	if err != nil {
		t.Fatal(err)
	}
	saver := runSaver()
	go func() {
		defer saver.close()
		for notification := range notificationCh {
			saver.strSaveCh() <- fmt.Sprint(notification)
			for _, balance := range notification.Data {
				if err := checkBalance(&balance); err != nil {
					saver.errSaveCh() <- err
					return
				}
			}
		}
	}()
	<-time.After(5 * time.Second)
	_, err = restClient.TransferBetweenWalletAndExchange(
		bg,
		args.Amount("0.2"),
		args.Currency("EOS"),
		args.Source(args.AccountWallet),
		args.Destination(args.AccountSpot),
	)
	if err != nil {
		fmt.Println(err)
	}
	<-time.After(5 * time.Second)
	_, err = restClient.TransferBetweenWalletAndExchange(
		bg,
		args.Amount("0.2"),
		args.Currency("EOS"),
		args.Source(args.AccountWallet),
		args.Destination(args.AccountSpot),
	)
	if err != nil {
		fmt.Println(err)
	}
	<-time.After(5 * time.Second)
	err = client.UnsubscribeToWalletBalance()
	if err != nil {
		t.Fatal(err)
	}
	saver.printSavedStrings()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fatal()
	}
	<-time.After(5 * time.Second)
}
