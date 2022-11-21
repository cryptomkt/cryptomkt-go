package rest

import (
	"strconv"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestGetSpotTradingBalance(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetSpotTradingBalances(bg)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkBalance, result); err != nil {
		t.Fatal(err)
	}
}

func TestGetSpotTradingBalanceOfCurrency(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetSpotTradingBalanceOfCurrency(bg, args.Currency("CRO"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkBalance(result); err != nil {
		t.Fatal(err)
	}
}

func TestGetAllActiveSpotOrders(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetAllActiveSpotOrders(bg)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkOrder, result); err != nil {
		t.Fatal(err)
	}
}

func TestCreateOrder(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.CreateSpotOrder(bg, args.Symbol("EOSETH"), args.Side(args.SideSell), args.Quantity("0.01"), args.Price("8999"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(result); err != nil {
		t.Fatal(err)
	}
}

func TestCancelAllOrders(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.CancelAllSpotOrders(bg)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkOrder, result); err != nil {
		t.Fatal(err)
	}
}

func TestOrderFlow(t *testing.T) {
	client, bg := beforeEach()
	order, err := client.CreateSpotOrder(
		bg,
		args.Symbol("EOSETH"),
		args.Side("sell"),
		args.Quantity("0.01"),
		args.Price("9999"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
	order, err = client.GetActiveSpotOrder(bg, args.ClientOrderID(order.ClientOrderID))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
	newClientOrderID := strconv.FormatInt(time.Now().UnixMicro(), 10)
	order, err = client.ReplaceSpotOrder(bg,
		args.ClientOrderID(order.ClientOrderID),
		args.NewClientOrderID(newClientOrderID),
		args.Quantity("0.02"),
		args.Price("9999"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
	if order.ClientOrderID != newClientOrderID {
		t.Fatal("should be the same client order id")
	}
	order, err = client.GetActiveSpotOrder(bg, args.ClientOrderID(newClientOrderID))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}

	order, err = client.CancelSpotOrder(bg, args.ClientOrderID(order.ClientOrderID))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkOrder(order); err != nil {
		t.Fatal(err)
	}
}

func TestGetAllTradingCommisions(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetAllTradingCommissions(bg)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkList(checkTradingCommission, result); err != nil {
		t.Fatal(err)
	}
}

func TestGetTradingCommision(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetTradingCommissionOfSymbol(bg, args.Symbol("eoseth"))
	if err != nil {
		t.Fatal(err)
	}
	if err = checkTradingCommission(result); err != nil {
		t.Fatal(err)
	}
}
