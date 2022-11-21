package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/internal"
	"github.com/cryptomarket/cryptomarket-go/models"
)

// PublicClient connects via websocket to cryptomarket to get market information of the exchange.
type MarketDataClient struct {
	clientBase
}

// NewPublicClient returns a new chan client if the connection with the
// cryptomarket server is successful, and error otherwise.
//
// Subscriptions reuse channels between symbols if there are any already.
// different speeds and depths means diferent channels.
// e.g. a subscription for the tickers of ETHBTC at speed '3s' uses the same channel that a new subscription for the tickers of XLMETH at speed '3s', while a new subscription is used for a new subscription for the tickers of USDTBTC at speed '1s'
//
// Unsubscriptions only closes the relevant channel, and does not make the server stop recieving the subscription data
// it is safe to unsubscribe to an unsubscribed channel
func NewMarketDataClient() (*MarketDataClient, error) {
	client := &MarketDataClient{
		clientBase: clientBase{
			wsManager: newWSManager("/api/3/ws/public"),
			chanCache: newChanCache(),
			window:    0,
		},
	}

	// connect to streaming
	err := client.wsManager.connect()
	if err != nil {
		return nil, fmt.Errorf("Error in websocket client connection: %s", err)
	}
	// handle incomming data
	go client.handle(client.wsManager.rcv)
	return client, nil
}

type channelSubscriptionResult struct {
	Subscriptions []string `json:"subscriptions"`
}
type channelSubscriptionResponse struct {
	Error  *models.APIError
	Result channelSubscriptionResult `json:"result"`
}

func (client *MarketDataClient) doChannelSubscription(
	method string,
	channel string,
	params map[string]interface{},
) (*struct {
	ch      chan []byte
	symbols []string
}, error) {
	if !client.wsManager.isOpen {
		return nil, fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.store(ch)
	notification := wsSubscription{
		ID:      id,
		Method:  method,
		Channel: channel,
		Params:  params,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return nil, fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	key := channel
	var dataOut chan []byte
	if ch, ok := client.chanCache.getCh(key); ok {
		dataOut = ch
	} else {
		dataOut = make(chan []byte, 1)
	}
	client.chanCache.storeSubscriptionCh(key, dataOut)
	client.wsManager.snd <- data
	data = <-ch
	var resp channelSubscriptionResponse
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		close(dataOut)
		return nil, fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
	}
	return &(struct {
		ch      chan []byte
		symbols []string
	}{ch: dataOut, symbols: resp.Result.Subscriptions}), nil
}

// UnsubscribeTo closes the receiving channel of a subscription, given his NorificationChannel name.
// further messages recieved from the server on the corresponding channel will be droped.
func (client *MarketDataClient) UnsubscribeTo(notificationChannel string) {
	client.chanCache.deleteSubscriptionCh(notificationChannel)
}

func convertChan[ft models.FeedType](dataCh chan []byte) chan models.Notification[ft] {
	notificationCh := make(chan models.Notification[ft], 1)
	go func() {
		defer close(notificationCh)
		var resp struct {
			Snapshot *ft
			Update   *ft
			Data     *ft
		}
		for data := range dataCh {
			resp.Snapshot = nil
			resp.Update = nil
			resp.Data = nil
			json.Unmarshal(data, &resp)
			if resp.Update != nil {
				notificationCh <- models.Notification[ft]{Data: *resp.Update, NotificationType: args.NotificationUpdate}
				continue
			}
			if resp.Snapshot != nil {
				notificationCh <- models.Notification[ft]{Data: *resp.Snapshot, NotificationType: args.NotificationSnapshot}
				continue
			}
			if resp.Data != nil {
				notificationCh <- models.Notification[ft]{Data: *resp.Data, NotificationType: args.NotificationData}
				continue
			}
		}
	}()
	return notificationCh
}

func (client *MarketDataClient) GetActiveSubscriptions(
	ctx context.Context,
	arguments ...args.Argument,
) ([]string, error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSubscription)
	if err != nil {
		return nil, err
	}
	if !client.wsManager.isOpen {
		return nil, fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.store(ch)
	var notificationWithChannel struct {
		ID      int64  `json:"id"`
		Method  string `json:"method"`
		Channel string `json:"ch"`
	}
	notificationWithChannel.ID = id
	notificationWithChannel.Method = methodSubscriptions
	notificationWithChannel.Channel = string(
		params[internal.ArgNameSubscription].(args.SubscriptionType),
	)
	data, err := json.Marshal(notificationWithChannel)
	if err != nil {
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return nil, fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	client.wsManager.snd <- data
	select {
	case <-ctx.Done():
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return nil, ctx.Err()
	case data := <-ch:
		var resp struct {
			Error *models.APIError
		}
		json.Unmarshal(data, &resp)
		if resp.Error != nil {
			return nil, fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
		}
		var response struct {
			Result struct {
				Subscriptions []string
			}
		}
		json.Unmarshal(data, &response)
		return response.Result.Subscriptions, nil
	}
}

func addAsteriscIfNoSymbols(params *map[string]interface{}) {
	if _, ok := (*params)[internal.ArgNameSymbols]; !ok {
		(*params)[internal.ArgNameSymbols] = []string{"*"}
	}
}

// SubscribeToTrades subscribe to a feed of trades
//
// subscription is only for the specified symbols
//
// normal subscriptions have one update message per symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-trades
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
//  Limit(int64)  // Number of historical entries returned in the first feed. Min is 0. Max is 1000. Default is 0
func (client *MarketDataClient) SubscribeToTrades(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSTradeFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbols)
	if err != nil {
		return nil, err
	}
	channel := internal.ChannelTrades
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSTradeFeed]{
		NotificationCh:      convertChan[models.WSTradeFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToCandles subscribe to a feed of candles
//
// subscription is only for the specified symbols
//
// normal subscriptions have one update message per symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-candles
//
// Arguments:
//  Period(PeriodType)  // Optional. A valid tick interval. 'M1' (one minute), 'M3', 'M5', 'M15', 'M30', 'H1' (one hour), 'H4', 'D1' (one day), 'D7', '1M' (one month). Default is 'M30'
//  Symbols([]string)  // Optional. A list of symbol ids
//  Limit(int64)  // Number of historical entries returned in the first feed. Min is 0. Max is 1000. Default is 0
func (client *MarketDataClient) SubscribeToCandles(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSCandleFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbols, internal.ArgNamePeriod)
	if err != nil {
		return nil, err
	}
	channel := fmt.Sprintf(internal.ChannelCandles, params[internal.ArgNamePeriod])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSCandleFeed]{
		NotificationCh:      convertChan[models.WSCandleFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToMiniTicker subscribe to a feed of mini tickers
//
// subscription is for all symbols or for the specified symbols
//
// normal subscriptions have one update message per symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-mini-ticker
//
// Arguments:
//  TickerSpeed(TickerSpeedType)  // The speed of the feed. TickerSpeed1s or TickerSpeed3s
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToMiniTicker(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.MiniTickerFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(internal.ChannelMiniTicker, params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.MiniTickerFeed]{
		NotificationCh:      convertChan[models.MiniTickerFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToMiniTickerInBatches subscribe to a feed of mini tickers
//
// subscription is for all symbols or for the specified symbols
//
// batch subscriptions have a joined update for all symbols
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-mini-ticker-in-batches
//
// Arguments:
//  TickerSpeed(TickerSpeedType)  // The speed of the feed. TickerSpeed1s or TickerSpeed3s
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToMiniTickerInBatches(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.MiniTickerFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(internal.ChannelMiniTickerInBatch, params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.MiniTickerFeed]{
		NotificationCh:      convertChan[models.MiniTickerFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToTicker subscribe to a feed of tickers
//
// subscription is for all symbols or for the specified symbols
//
// normal subscriptions have one update message per symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-ticker
//
// Arguments:
//  TickerSpeed(TickerSpeedType)  // The speed of the feed. TickerSpeed1s or TickerSpeed3s
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToTicker(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSTickerFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	fmt.Println(params)
	channel := fmt.Sprintf(internal.ChannelTicker, params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSTickerFeed]{
		NotificationCh:      convertChan[models.WSTickerFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToTickerInBatches subscribe to a feed of tickers
//
// subscription is for all symbols or for the specified symbols
//
// batch subscriptions have a joined update for all symbols
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-ticker-in-batches
//
// Arguments:
//  TickerSpeed(TickerSpeedType)  // The speed of the feed. TickerSpeed1s or TickerSpeed3s
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToTickerInBatches(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSTickerFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(internal.ChannelTickerInBatch, params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSTickerFeed]{
		NotificationCh:      convertChan[models.WSTickerFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToFullOrderbook subscribe to a feed of a full orderbook
//
// subscription is only for the specified symbols
//
// normal subscriptions have one update message per symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-full-order-book
//
// Arguments:
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToFullOrderbook(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSOrderbookFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSymbols)
	if err != nil {
		return nil, err
	}
	channel := internal.ChannelOrderBookFull
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSOrderbookFeed]{
		NotificationCh:      convertChan[models.WSOrderbookFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToPartialOrderbook subscribe to a feed of a partial orderbook
//
// subscription is for all symbols or for the specified symbols
//
// normal subscriptions have one update message per symbol
//
// Requires no API key Access Rights
//
// https://api.exchange.cryptomkt.com/#subscribe-to-partial-order-book
//
// Arguments:
//  OrderBookSpeed(OrderBookSpeedType)  // The speed of the feed. OrderBookSpeedType100ms, OrderBookSpeedType500ms or OrderBookSpeedType1000ms
//  WSDepth(WSDepthType)  // The depth of the partial orderbook, WSDepth5, WSDepth10 or WSDepth20
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToPartialOrderbook(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSOrderbookFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(internal.ChannelOrderbookPartial, params[internal.ArgNameDepth], params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSOrderbookFeed]{
		NotificationCh:      convertChan[models.WSOrderbookFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToPartialOrderbookInBatches subscribe to a feed of a partial orderbook in batches
//
// subscription is for all symbols or for the specified symbols
//
// batch subscriptions have a joined update for all symbols
//
//
//
// https://api.exchange.cryptomkt.com/#subscribe-to-partial-order-book-in-batches
//
// Arguments:
//  OrderBookSpeed(OrderBookSpeedType)  // The speed of the feed. OrderBookSpeedType100ms, OrderBookSpeedType500ms or OrderBookSpeedType1000ms
//  WSDepth(WSDepthType)  // The depth of the partial orderbook, WSDepth5, WSDepth10 or WSDepth20
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToPartialOrderbookInBatchers(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.WSOrderbookFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(
		internal.ChannelOrderbookPartialInBatch,
		params[internal.ArgNameDepth],
		params[internal.ArgNameSpeed],
	)
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.WSOrderbookFeed]{
		NotificationCh:      convertChan[models.WSOrderbookFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToTopOfOrderbook subscribe to a feed of the top of the orderbook
//
// subscription is for all symbols or for the specified symbols
//
// normal subscriptions have one update message per symbol
//
//
//
// https://api.exchange.cryptomkt.com/#subscribe-to-top-of-book
//
// Arguments:
//  OrderBookSpeed(OrderBookSpeedType)  // The speed of the feed. OrderBookSpeedType100ms, OrderBookSpeedType500ms or OrderBookSpeedType1000ms
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToOrderbookTop(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.OrderbookTopFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(internal.ChannelOrderbookTop, params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.OrderbookTopFeed]{
		NotificationCh:      convertChan[models.OrderbookTopFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}

// SubscribeToTopOfOrderbookInBatches subscribe to a feed of the top of the orderbook
//
// subscription is for all symbols or for the specified symbols
//
// batch subscriptions have a joined update for all symbols
//
//
//
// https://api.exchange.cryptomkt.com/#subscribe-to-top-of-book-in-batches
//
// Arguments:
//  OrderBookSpeed(OrderBookSpeedType)  // The speed of the feed. OrderBookSpeedType100ms, OrderBookSpeedType500ms or OrderBookSpeedType1000ms
//  Symbols([]string)  // Optional. A list of symbol ids
func (client *MarketDataClient) SubscribeToOrderbookTopInBatchers(
	arguments ...args.Argument,
) (subscription *models.Subscription[models.OrderbookTopFeed], err error) {
	params, err := args.BuildParams(arguments, internal.ArgNameSpeed)
	if err != nil {
		return nil, err
	}
	addAsteriscIfNoSymbols(&params)
	channel := fmt.Sprintf(internal.ChannelOrderbookTopInBatch, params[internal.ArgNameSpeed])
	response, err := client.doChannelSubscription(
		methodSubscribe,
		channel,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &models.Subscription[models.OrderbookTopFeed]{
		NotificationCh:      convertChan[models.OrderbookTopFeed](response.ch),
		Symbols:             response.symbols,
		NotificationChannel: channel,
	}, nil
}
