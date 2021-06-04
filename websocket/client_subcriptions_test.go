package websocket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptomkt/go-api/args"
)

func TestTickerSubscription(t *testing.T) {
	client, _ := New("", "")
	bg := context.Background()

	checker := newTimeFlowChecker()
	time.Sleep(5) // let the time flow a little
	feedCh, err := client.subscribeToTicker(bg, args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	errCh := make(chan error)
	go func() {
		for ticker := range feedCh {
			// check the right ordering of timestamps
			err := checker.checkNextTime(ticker.Timestamp)
			if err != nil {
				errCh <- err
				return
			}
		}
	}()
	select {
	case err:=<-errCh: // terminate if error
		t.Fatal(err)
	case <-time.After(30 * time.Second): // wait 10 seconds
	}
	err = client.unsubscribeToTicker(bg, args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	client.Close()
	time.Sleep(5 * time.Second)
}

func TestMultipleTickerSubscriptions(t *testing.T) {
	client, _ := New("", "")
	bg := context.Background()
	result, _ := client.getSymbols(bg)
	n := 7
	symbolList := make([]string, n)
	for idx := 0; idx < n; idx++ {
		symbolList[idx] = result[idx].ID
	}
	for _, symbol := range symbolList {
		feedCh, err := client.subscribeToTicker(bg, args.Symbol(symbol))
		if err != nil {
			t.Fatal(err)
		}
		go printSymbolOfTicker(feedCh, symbol)

	}
	time.Sleep(30 * time.Second)
	for _, symbol := range symbolList {
		err := client.unsubscribeToTicker(bg, args.Symbol(symbol))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func TestDeadlineFromContext(t *testing.T) {
	client, _ := New("", "")
	bg := context.Background()
	ctx, cancelFunc := context.WithDeadline(bg, time.Now().Add(50*time.Millisecond))
	defer cancelFunc()

	symbol := "EOSETH"
	_, err := client.subscribeToTicker(ctx, args.Symbol(symbol))
	if err == nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
}

func TestOrderbookSubscription(t *testing.T) {
	client, _ := New("", "")
	ctx := context.Background()
	feedCh, err := client.subscribeToOrderbooks(ctx, args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for orderbook := range feedCh {
			fmt.Printf("ask len:%v\tbid len:%v\n", len(orderbook.Ask), len(orderbook.Bid))
			// for _, side := range [][]models.BookLevel{orderbook.Ask, orderbook.Bid} {
			// 	for _, level := range side {
			// 		size, _ := new(big.Float).SetString(level.Size)
			// 		zero := new(big.Float)
			// 		if size.Cmp(zero) == 0 {
			// 			t.Fatal("invalid size")
			// 		}
			// 	}

			// }
		}
		fmt.Println("channel closed")
	}()
	time.Sleep(30 * time.Second)
	// err = client.unsubscribeToTicker(ctx, args.Symbol("EOSETH"))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	time.Sleep(5 * time.Second)
}
