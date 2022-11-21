package websocket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

func beforeEachMarketDataTest() (client *MarketDataClient, saver *saver) {
	client, _ = NewMarketDataClient()
	saver = runSaver()
	return client, saver
}

func afterEach(t *testing.T, client *MarketDataClient, saver *saver, notificationChannel string) {
	<-time.After(50 * time.Second)
	client.UnsubscribeTo(notificationChannel)
	<-time.After(5 * time.Second)
	client.Close()
	<-time.After(5 * time.Second)
	saver.printSavedStrings()
	saver.printSavedErrors()
	if saver.wereErrorsPrinted() {
		t.Fail()
	}
}

func dropData[ft models.FeedType](subscription *models.Subscription[ft]) {
	for range subscription.NotificationCh {
		continue
	}
}

func TestUnsubscription(t *testing.T) {
	client, _ := NewMarketDataClient()
	saver := runSaver()
	subscription, _ := client.SubscribeToFullOrderbook(args.Symbols([]string{"EOSETH"}))
	go func() {
		for notification := range subscription.NotificationCh {
			for key, orderbook := range notification.Data {
				saver.strSaveCh() <- fmt.Sprint(key, ":", orderbook.Timestamp)
			}
		}
	}()
	<-time.After(10 * time.Second)
	client.UnsubscribeTo(subscription.NotificationChannel)
	saver.strSaveCh() <- fmt.Sprint(time.Now().UnixMilli())
	<-time.After(10 * time.Second)
	client.Close()
	saver.close()
	saver.printSavedStrings()
}

func TestMultipleSubscriptionsSameChannel(t *testing.T) {
	client, _ := NewMarketDataClient()
	saver := runSaver()
	subscription, _ := client.SubscribeToMiniTicker(
		args.Symbols([]string{"EOSETH"}),
		args.TickerSpeed(args.TickerSpeedType(args.TickerSpeed1s)),
	)
	saver.strSaveCh() <- fmt.Sprint("subscription:", subscription.Symbols)
	go func() {
		for notification := range subscription.NotificationCh {
			saver.strSaveCh() <- fmt.Sprint(notification)
		}
	}()
	<-time.After(40 * time.Second)
	subscription, _ = client.SubscribeToMiniTicker(
		args.Symbols([]string{"ETHBTC"}),
		args.TickerSpeed(args.TickerSpeedType(args.TickerSpeed1s)),
	)
	saver.strSaveCh() <- fmt.Sprint("subscription:", subscription.Symbols)
	<-time.After(20 * time.Second)
	client.Close()
	saver.close()
	saver.printSavedStrings()
}

func TestClosingClientClosesAllChannels(t *testing.T) {
	client, err := NewMarketDataClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	subscription_1, err := client.SubscribeToCandles(
		args.Symbols([]string{"EOSETH"}),
		args.Period(args.Period1Minute),
	)
	if err != nil {
		t.Fatal(err)
	}
	go dropData(subscription_1)
	subscription_2, err := client.SubscribeToTrades(args.Symbols([]string{"EOSETH"}))
	if err != nil {
		t.Fatal(err)
	}
	go dropData(subscription_2)
	subscription_3, err := client.SubscribeToMiniTicker(args.TickerSpeed(args.TickerSpeed1s))
	if err != nil {
		t.Fatal(err)
	}
	go dropData(subscription_3)
	subscription_4, err := client.SubscribeToFullOrderbook(args.Symbols([]string{"ETHBTC"}))
	if err != nil {
		t.Fatal(err)
	}
	go dropData(subscription_4)
	subscription_5, err := client.SubscribeToOrderbookTop(args.OrderBookSpeed(args.OrderBookSpeed100ms))
	if err != nil {
		t.Fatal(err)
	}
	go dropData(subscription_5)
	<-time.After(10 * time.Second)
}

func TestAsteriscIfNoSymbols(t *testing.T) {
	client, _ := NewMarketDataClient()
	defer client.Close()
	subscription, _ := client.SubscribeToMiniTicker()
	if len(subscription.Symbols) == 0 {
		t.Fail()
	}
}
func TestTradesSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToTrades(args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(10))
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			saver.strSaveCh() <- fmt.Sprint(notification)
			for _, tradeList := range notification.Data {
				for _, trade := range tradeList {
					if err := checkWSTrade(&trade); err != nil {
						saver.errSaveCh() <- err
					}
				}
			}
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestCandlesSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToCandles(args.Symbols([]string{"EOSETH", "XLMETH"}), args.Period(args.Period1Minute))
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			saver.strSaveCh() <- fmt.Sprint(notification.NotificationType)
			checkCandleFeed(saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestMinitickerSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToMiniTicker(
		args.Symbols([]string{"EOSETH"}),
		args.TickerSpeed(args.TickerSpeed3s),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			checkMiniTickerFeed(saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestMiniTickerInBatchSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToMiniTickerInBatches(
		args.Symbols([]string{"EOSETH"}),
		args.TickerSpeed(args.TickerSpeed1s),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			checkMiniTickerFeed(saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestTickerSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToTicker(
		args.TickerSpeed(args.TickerSpeed1s),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notificationCh := range subscription.NotificationCh {
			checkWSTickerFeed(saver, &notificationCh.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestTickerInBatchesSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToTickerInBatches(
		args.Symbols([]string{"EOSETH"}),
		args.TickerSpeed(args.TickerSpeed1s),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			checkWSTickerFeed(saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestFullOrderbookSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToFullOrderbook(
		args.Symbols([]string{"EOSETH"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	obchecker := new(OBChecker)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			checkOrderbookFeed(obchecker, saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestPartialOrderbookSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToPartialOrderbook(
		args.Symbols([]string{"EOSETH"}),
		args.WSDepth(args.WSDepth10),
		args.OrderBookSpeed(args.OrderBookSpeed1000ms),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	obchecker := new(OBChecker)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			checkOrderbookFeed(obchecker, saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestPartialOrderbookInBatchesSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToPartialOrderbookInBatchers(
		args.Symbols([]string{"EOSETH"}),
		args.WSDepth(args.WSDepth10),
		args.OrderBookSpeed(args.OrderBookSpeed1000ms),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	obchecker := new(OBChecker)
	go func() {
		defer saver.close()
		for feed := range subscription.NotificationCh {
			checkOrderbookFeed(obchecker, saver, &feed.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestOrderbookTopSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToOrderbookTop(
		args.Symbols([]string{"EOSETH"}),
		args.OrderBookSpeed(args.OrderBookSpeed1000ms),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for feed := range subscription.NotificationCh {
			checkOrderbookTopFeed(saver, &feed.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestOrderbookTopInBatchesSubscription(t *testing.T) {
	client, saver := beforeEachMarketDataTest()
	subscription, err := client.SubscribeToOrderbookTopInBatchers(
		args.Symbols([]string{"EOSETH"}),
		args.OrderBookSpeed(args.OrderBookSpeed1000ms),
	)
	if err != nil {
		t.Fatal(err)
	}
	saver.strSaveCh() <- fmt.Sprint(subscription.Symbols)
	go func() {
		defer saver.close()
		for notification := range subscription.NotificationCh {
			checkOrderbookTopFeed(saver, &notification.Data)
		}
	}()
	afterEach(t, client, saver, subscription.NotificationChannel)
}

func TestGetActiveSubscriptions(t *testing.T) {
	client, _ := beforeEachMarketDataTest()
	_, err := client.SubscribeToOrderbookTopInBatchers(
		args.Symbols([]string{"EOSETH"}),
		args.OrderBookSpeed(args.OrderBookSpeed1000ms),
	)
	_, err = client.SubscribeToTrades(args.Symbols([]string{"ETHBTC"}))
	<-time.After(3 * time.Second)
	result, err := client.GetActiveSubscriptions(
		context.Background(),
		args.Subscription(args.SubscriptionTypeOrderbookTopInBatch(args.OrderBookSpeed1000ms)),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatal("wrong number of subscriptions")
	}
}
