package conn

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
)

func TestPrices(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	rand.Seed(time.Now().UnixNano())
	optional := [2]args.Argument{args.Page(0), args.Limit(50)}
	for i := 0; i < 100; i++ {
		var numArgs int = rand.Intn(3)
		switch numArgs {
		case 0:
			if _, err := client.GetPrices(
				args.Market("ETHCLP"),
				args.Timeframe("60")); err != nil {
				t.Errorf("Prices with zero optional args failed, %s", err)
			}
		case 1:
			var randomIndex int = rand.Intn(2)
			if _, err := client.GetPrices(
				args.Market("ETHCLP"),
				args.Timeframe("60"),
				optional[randomIndex]); err != nil {
				t.Errorf("Prices with one optional argument failed, %s", err)
			}
		case 2:
			if _, err := client.GetPrices(
				args.Market("ETHCLP"),
				args.Timeframe("60"),
				optional[0],
				optional[1]); err != nil {
				t.Errorf("Prices with 2 optional args failed, %s", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func TestPrevious(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	prices, err := client.GetPrices(args.Market("ETHCLP"), args.Timeframe("60"), args.Page(1), args.Limit(40))
	if err != nil {
		t.Errorf("Error in previous prices: %s", err)
	}
	_, err = prices.GetPrevious()
	if err != nil {
		t.Errorf("Error in previous prices: %s", err)
	}
}

func TestNext(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	prices, err := client.GetPrices(args.Market("ETHCLP"), args.Timeframe("10080"), args.Page(0), args.Limit(40))
	if err != nil {
		t.Errorf("Error in next prices: %s", err)
	}
	_, err = prices.GetNext()
	if err != nil {
		t.Errorf("Error in next prices: %s", err)
	}
}

func TestGetAllPrices(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	time.Sleep(3 * time.Second)
	if _, err := client.GetPrices(args.Market("ETHCLP"), args.Timeframe("10080")); err != nil {
		t.Errorf("Prices with zero optional args failed: %s", err)
	}
}
