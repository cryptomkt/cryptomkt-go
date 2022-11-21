package rest

import (
	"context"
	"fmt"
	"testing"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func beforeEach() (*Client, context.Context) {
	apiKeys := LoadKeys()
	return NewClient(apiKeys.APIKey, apiKeys.APISecret, 1000), context.Background()
}

func TestGetWalletBalances(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetWalletBallances(bg)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("should have results")
	}
	if err = checkList(checkBalance, result); err != nil {
		t.Fatal(err)
	}
}

func TestGetWalletBalanceOfCurrency(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetWalletBalanceOfCurrency(
		bg,
		args.Currency("CRO"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkBalance(result); err != nil {
		t.Fatal(err)
	}
}

func TestGetDepositCryptoAddresses(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetDepositCryptoAddresses(bg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
	if err = checkList(checkCryptoAddress, result); err != nil {
		t.Fatal(err)
	}
}

func TestGetDepositCryptoAddressOfCurrency(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetDepositCryptoAddressOfCurrency(
		bg,
		args.Currency("USDT"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkCryptoAddress(result); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDepositCryptoAddress(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.CreateDepositCryptoAddress(bg, args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkCryptoAddress(result); err != nil {
		t.Fatal(err)
	}
}

func TestGetLast10DepositCryptoAddress(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetLast10DepositCryptoAddresses(bg, args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkCryptoAddress, result); err != nil {
		t.Fatal(err)
	}
}

func TestGetLast10WithdrawalCryptoAddresses(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetLast10WithdrawalCryptoAddresses(bg, args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkCryptoAddress, result); err != nil {
		t.Fatal(err)
	}
}

func TestWithdrawCrypto(t *testing.T) {
	client, bg := beforeEach()
	adaAddress, err := client.GetDepositCryptoAddressOfCurrency(bg, args.Currency("ADA"))
	if err != nil {
		t.Fatal(err)
	}
	transactionID, err := client.withdrawCrypto(bg,
		args.Currency("ADA"),
		args.Amount("0.1"),
		args.Address(adaAddress.Address),
	)
	if transactionID == "" {
		t.Fatal("no transaction id")
	}
}

func TestWithdrawCryptoCommit(t *testing.T) {
	client, bg := beforeEach()
	adaAddress, err := client.GetDepositCryptoAddressOfCurrency(bg, args.Currency("ADA"))
	if err != nil {
		t.Fatal(err)
	}
	transactionID, err := client.withdrawCrypto(bg,
		args.Currency("ADA"),
		args.Amount("0.1"),
		args.Address(adaAddress.Address),
		args.AutoCommit(false),
	)
	if transactionID == "" {
		t.Fatal("no transaction id")
	}
	success, err := client.WithdrawCryptoCommit(bg, args.ID(transactionID))
	if err != nil {
		t.Fatal(err)
	}
	if !success {
		t.Fatal("failed to commit withdrawal")
	}
}

func TestWithdrawCryptoRollback(t *testing.T) {
	client, bg := beforeEach()
	adaAddress, err := client.GetDepositCryptoAddressOfCurrency(bg, args.Currency("ADA"))
	if err != nil {
		t.Fatal(err)
	}
	transactionID, err := client.withdrawCrypto(bg,
		args.Currency("ADA"),
		args.Amount("0.1"),
		args.Address(adaAddress.Address),
		args.AutoCommit(false),
	)
	if transactionID == "" {
		t.Fatal("no transaction id")
	}
	success, err := client.WithdrawCryptoCommit(bg, args.ID(transactionID))
	if err != nil {
		t.Fatal(err)
	}
	if !success {
		t.Fatal("failed to commit withdrawal")
	}
}

func TestGetEstimateWithdrawFee(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetEstimateWithdrawFee(bg, args.Currency("EOS"), args.Amount("199"))
	if err != nil {
		t.Fatal(err)
	}
	if result == "" {
		t.Fatal("should have a result")
	}
}

func TestCryptoAddressBelongToCurrentAccount(t *testing.T) {
	client, bg := beforeEach()
	adaAddress, err := client.GetDepositCryptoAddressOfCurrency(bg, args.Currency("ADA"))
	if err != nil {
		t.Fatal(err)
	}
	belongs, err := client.CheckIfCryptoAddressBelongsToCurrentAccount(bg, args.Address(adaAddress.Address))
	if err != nil {
		t.Fatal(err)
	}
	if !belongs {
		t.Fatal("should belong")
	}
}

func TestTransferBetweenWalletAndExchange(t *testing.T) {
	client, bg := beforeEach()
	currency := "EOS"
	amount := "0.01"

	transactionID, err := client.TransferBetweenWalletAndExchange(bg,
		args.Currency(currency),
		args.Amount(amount),
		args.Source(args.AccountWallet),
		args.Destination(args.AccountSpot),
	)
	if err != nil {
		t.Fatal(err)
	}
	if transactionID == "" {
		t.Fatal("no transaction ID")
	}

	transactionID, err = client.TransferBetweenWalletAndExchange(bg,
		args.Currency(currency),
		args.Amount(amount),
		args.Source(args.AccountSpot),
		args.Destination(args.AccountWallet),
	)
	if err != nil {
		t.Fatal(err)
	}
	if transactionID == "" {
		t.Fatal("no transaction ID")
	}
}

func TestTransferMoneyToAnotherUser(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.TransferMoneyToAnotherUser(bg,
		args.Currency("CRO"),
		args.Amount("0.1"),
		args.IdentifyBy(args.IdentifyByEmail),
		args.Identifier("the_email"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result == "" {
		t.Fatal("no transaction ID")
	}
}

func TestGetTransactionHistory(t *testing.T) {
	client, bg := beforeEach()
	transactionsList := []args.TransactionTypeType{args.TransactionTypeDeposit}
	result, err := client.GetTransactionHistory(bg, args.TransactionTypes(transactionsList...))
	fmt.Println(result)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("should have transactions")
	}
	if err = checkList(checkTransaction, result); err != nil {
		t.Fatal(err)
	}
}

func TestGetTransaction(t *testing.T) {
	// see ruby sdk for further information
}

func TestOffchainAvailable(t *testing.T) {
	client, bg := beforeEach()
	eosAddress, err := client.GetDepositCryptoAddressOfCurrency(bg, args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CheckIfOffchainIsAvailable(bg, args.Currency("EOS"), args.Address(eosAddress.Address))
	if err != nil {
		t.Fatal(err)
	}
}
