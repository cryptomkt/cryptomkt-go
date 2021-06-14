package websocket

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

func TestTickerSubscription(t *testing.T) {
	client, _ := NewPublicClient()
	feedCh, err := client.SubscribeToTicker(args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	timeFlowCh := make(chan error)
	go func() {
		defer close(timeFlowCh)
		checker := newTimeFlowChecker()
		for ticker := range feedCh {
			err = checkTicker(&ticker)
			if err != nil {
				timeFlowCh <- err
				return
			}
			err := checker.checkNextTime(ticker.Timestamp)
			if err != nil {
				timeFlowCh <- err
				return
			}
		}
	}()
	select {
	case err := <-timeFlowCh:
		t.Fatal(err)
	case <-time.After(30 * time.Second): // wait 30 seconds
	}
	err = client.UnsubscribeToTicker(args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	client.Close()
	time.Sleep(5 * time.Second)
}

func TestMultipleTickerSubscriptions(t *testing.T) {
	client, _ := NewPublicClient()
	bg := context.Background()
	result, err := client.GetSymbols(bg)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(bg, 30*time.Second)
	defer cancelFunc()

	n := 5
	var wg sync.WaitGroup
	wg.Add(n)
	innerErrCh := make(chan error, n)

	for idx := 0; idx < n; idx++ {
		symbol := result[idx].ID
		feedCh, err := client.SubscribeToTicker(args.Symbol(symbol))
		if err != nil {
			t.Fatal(err)
		}
		go func(symbol string, feedCh chan models.Ticker, errorCh chan error) {
			for ticker := range feedCh {
				if err := checkTicker(&ticker); err != nil {
					errorCh <- fmt.Errorf("%s: %v", symbol, err)
					break
				}
			}
			wg.Done()
		}(symbol, feedCh, innerErrCh)
	}
	rcvngError := true
	for rcvngError {
		select {
		case err := <-innerErrCh:
			t.Fatal(err)
		case <-ctx.Done():
			rcvngError = false
		}
	}
	for idx := 0; idx < n; idx++ {
		symbol := result[idx].ID
		err := client.UnsubscribeToTicker(args.Symbol(symbol))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
	close(innerErrCh)
	wg.Wait()
}

func TestOrderbookSubscription(t *testing.T) {
	client, _ := NewPublicClient()
	feedCh, err := client.SubscribeToOrderbook(args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	innerErrCh := make(chan error)
	go func() {
		defer close(innerErrCh)
		checker := newTimeFlowChecker()
		for orderbook := range feedCh {
			fmt.Printf("bid:%v\task:%v\n", len(orderbook.Bid), len(orderbook.Ask))
			if err := checkOrderbook(&orderbook); err != nil {
				innerErrCh <- err
				break
			}
			if err = checker.checkNextTime(orderbook.Timestamp); err != nil {
				innerErrCh <- err
				break
			}
		}
	}()
	select {
	case err = <-innerErrCh:
		t.Fatal(err)
	case <-time.After(8 * time.Minute):
	}
	err = client.UnsubscribeToOrderbook(args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCandlesSubscription(t *testing.T) {
	client, _ := NewPublicClient()
	feedCh, err := client.SubscribeToCandles(args.Symbol("EOSETH"), args.Period(args.PeriodType15Minutes))
	if err != nil {
		t.Fatal(err)
	}
	innerErrCh := make(chan error)
	go func() {
		defer close(innerErrCh)
		for candles := range feedCh {
			for _, candle := range candles {
				err := checkCandle(&candle)
				if err != nil {
					innerErrCh <- err
					return
				}
			}
		}
	}()
	select {
	case err := <-innerErrCh:
		t.Fatal(err)
	case <-time.After(1 * time.Minute):
	}
	err = client.UnsubscribeToCandles(args.Symbol("EOSETH"), args.Period(args.PeriodType15Minutes))
	if err != nil {
		t.Fatal(err)
	}
}

func TestTradesSubscription(t *testing.T) {
	client, _ := NewPublicClient()
	feedCh, err := client.SubscribeToTrades(args.Symbol("EOSETH"))
	if err != nil {
		t.Fatal(err)
	}
	innerErrCh := make(chan error)
	go func() {
		defer close(innerErrCh)
		for trades := range feedCh {
			for _, trade := range trades {
				err := checkPublicTrade(&trade)
				if err != nil {
					innerErrCh <- err
					return
				}
			}
		}
	}()
	select {
	case err = <-innerErrCh:
		t.Fatal(err)
	case <-time.After(1 * time.Minute):
	}
	if err = client.UnsubscribeToTrades(args.Symbol("EOSETH")); err != nil {
		t.Fatal(err)
	}
}
