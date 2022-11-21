package websocket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func beforeEachTradingClientTest() (
	client *SpotTradingClient,
	saver *saver,
	bg context.Context,
	err error,
) {
	apiKeys := LoadKeys()
	client, err = NewSpotTradingClient(apiKeys.APIKey, apiKeys.APISecret, 20_000)
	bg = context.Background()
	saver = runSaver()
	return
}

func TestAuth(t *testing.T) {
	apiKeys := LoadKeys()
	if _, err := NewSpotTradingClient(apiKeys.APIKey, apiKeys.APISecret, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := NewSpotTradingClient("no apikey", "no apisecret", 0); err == nil {
		t.Fatal(err)
	}

}

func TestGetTradingBalances(t *testing.T) {
	client, _, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	balances, err := client.GetSpotTradingBalances(bg)
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

func TestOrderFlow(t *testing.T) {
	client, saver, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		clientOrderID := fmt.Sprint(time.Now().Unix())
		report, err := client.CreateSpotOrder(
			bg,
			args.Symbol("EOSETH"),
			args.Side(args.SideSell),
			args.Price("1000"),
			args.Quantity("0.01"),
			args.ClientOrderID(clientOrderID),
		)
		if err != nil {
			saver.errSaveCh() <- err
			err = nil
		}
		if err := checkReport(report); err != nil {
			saver.errSaveCh() <- err
		}
		saver.strSaveCh() <- fmt.Sprint(report)
		defer saver.close()
		newClientOrderID := clientOrderID + "new"
		report, err = client.ReplaceSpotOrder(
			bg,
			args.ClientOrderID(clientOrderID),
			args.NewClientOrderID(newClientOrderID),
			args.Quantity("0.02"),
			args.Price("2000"),
		)
		saver.strSaveCh() <- fmt.Sprint(report)
		if err != nil {
			saver.errSaveCh() <- err
			err = nil
		}
		if err := checkReport(report); err != nil {
			saver.errSaveCh() <- err
		}
		report, err = client.CancelSpotOrder(
			bg,
			args.ClientOrderID(newClientOrderID),
		)
		if err != nil {
			saver.errSaveCh() <- err
			err = nil
		}
		if err := checkReport(report); err != nil {
			saver.errSaveCh() <- err
		}
		saver.strSaveCh() <- fmt.Sprint(report)
	}()
	saver.printSavedStrings()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fail()
	}
}

func TestReportsSubscription(t *testing.T) {
	client, saver, bg, err := beforeEachTradingClientTest()
	notificationCh, err := client.SubscribeToReports()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer saver.close()
		for notification := range notificationCh {
			saver.strSaveCh() <- fmt.Sprint(notification)
			for _, report := range notification.Data {
				if err := checkReport(&report); err != nil {
					saver.errSaveCh() <- err
				}
			}
		}
	}()
	<-time.After(5 * time.Second)
	clientOrderID := fmt.Sprint(time.Now().Unix())
	client.CreateSpotOrder(
		bg,
		args.Symbol("EOSETH"),
		args.Side(args.SideSell),
		args.Price("1000"),
		args.Quantity("0.01"),
		args.ClientOrderID(clientOrderID),
	)
	<-time.After(5 * time.Second)
	client.CancelSpotOrder(
		bg,
		args.ClientOrderID(clientOrderID),
	)
	<-time.After(5 * time.Second)
	client.UnsubscribeToReports()

	saver.printSavedStrings()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fail()
	}
}

func TestGetActiveSpotOrdersAndCancelAllSpotOrders(t *testing.T) {
	client, saver, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.CancelAllSpotOrders(bg)
	if err != nil {
		saver.errSaveCh() <- fmt.Errorf("fail on preparation (empty) of order list: %v", err)
	}
	for i := 0; i < 4; i++ {
		_, err := client.CreateSpotOrder(
			bg,
			args.Symbol("EOSETH"),
			args.Side(args.SideSell),
			args.Price("1000"),
			args.Quantity("0.01"),
		)
		if err != nil {
			saver.errSaveCh() <- fmt.Errorf("fail on creation of order: %v", err)
		}
	}
	<-time.After(2 * time.Second)
	reports, err := client.GetActiveSpotOrders(bg)
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range reports {
		if err = checkReport(&report); err != nil {
			saver.errSaveCh() <- err
		}
	}
	if len(reports) != 4 {
		saver.errSaveCh() <- fmt.Errorf("wrong number of reports:%v", len(reports))
	}
	client.CancelAllSpotOrders(bg)
	saver.close()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fail()
	}
}

func TestCreateSpotOrderList(t *testing.T) {}

func TestGetSpotTradingBalances(t *testing.T) {
	client, saver, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	balances, err := client.GetSpotTradingBalances(bg)
	if err != nil {
		t.Fatal(err)
	}
	for _, balance := range balances {
		if err = checkBalance(&balance); err != nil {
			saver.errSaveCh() <- err
		}
	}
	saver.close()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fail()
	}
}

func TestGetSpotTradingBalance(t *testing.T) {
	client, _, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	balance, err := client.GetSpotTradingBalanceOfCurrency(bg, args.Currency("ADA"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkBalance(balance); err != nil {
		t.Fatal(err)
	}
}

func TestSpotFees(t *testing.T) {
	client, saver, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	commissions, err := client.GetTradingCommissions(bg)
	if err != nil {
		t.Fatal(err)
	}
	for _, commission := range commissions {
		if err = checkCommission(&commission); err != nil {
			saver.errSaveCh() <- err
		}
	}
	saver.close()
	saver.printSavedErrors()
	if saver.errorsPrinted {
		t.Fail()
	}
}

func TestSpotFee(t *testing.T) {
	client, _, bg, err := beforeEachTradingClientTest()
	if err != nil {
		t.Fatal(err)
	}
	commission, err := client.GetSpotFee(bg, args.Symbol("ADAETH"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkCommission(commission); err != nil {
		t.Fatal(err)
	}
}
