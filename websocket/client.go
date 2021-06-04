package websocket

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cryptomkt/go-api/args"

	"github.com/cryptomkt/go-api/models"
)

const (
	publicCall  bool = true
	tradingCall bool = false
)

type Client struct {
	tradingWSManager *wsManager
	publicWSManager  *wsManager
	chanCache        *chanCache
	subsCache        subscriptionCache
	apiKey           string
	apiSecret        string
}

func New(apiKey, apiSecret string) (*Client, error) {
	client := &Client{
		publicWSManager:  newPublicWSManager(),
		tradingWSManager: newTradingWSManager(),
		chanCache:        newChanCache(),
		subsCache:        newSubscriptionCache(),
		apiKey:           apiKey,
		apiSecret:        apiSecret,
	}

	// connect to streamings
	err := client.publicWSManager.connect()
	if err != nil {
		return nil, fmt.Errorf("Error in Client connection: %s", err)
	}
	err = client.tradingWSManager.connect()
	if err != nil {
		return nil, fmt.Errorf("Error in Client connection: %s", err)
	}
	// handle incomming data
	go client.handle(client.publicWSManager.rcv)
	go client.handle(client.tradingWSManager.rcv)
	return client, nil
}

// Close close all the channels related to the client as well as the websocket connection.
// trying to make requests over a closed client will result in error.
func (client *Client) Close() {
	client.publicWSManager.close()
	client.tradingWSManager.close()
	client.subsCache.close()
	client.chanCache.close()
}


func (client *Client) doPublicRequest(ctx context.Context, method string, arguments []args.Argument, requiredArguments []string, model interface{}) error {
	return client.doRequest(ctx, method, arguments, requiredArguments, model, publicCall)
}

func (client *Client) doTradingRequest(ctx context.Context, method string, arguments []args.Argument, requiredArguments []string, model interface{}) error {
	return client.doRequest(ctx, method, arguments, requiredArguments, model, tradingCall)
}

func (client *Client) doRequest(ctx context.Context, method string, arguments []args.Argument, requiredArguments []string, model interface{}, public bool) error {
	params, err := args.BuildParams(arguments, requiredArguments...)
	if err != nil {
		return err
	}
	if !client.publicWSManager.isOpen && public {
		return fmt.Errorf("CryptomarketSDKError: websocket for public methods closed")
	}

	if !client.tradingWSManager.isOpen && !public {
		return fmt.Errorf("CryptomarketSDKError: websocket for trading methods closed")
	}
	chans := makeChans()
	id := client.chanCache.store(chans)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}

	data, err := json.Marshal(notification)
	if err != nil {
		chans.close()
		return fmt.Errorf("CryptomarketSDKError: invalid notification: %v", notification)
	}
	if public {
		client.publicWSManager.snd <- data
	} else { // trading
		client.tradingWSManager.snd <- data
	}
	select {
	case <-ctx.Done():
		if chans, ok := client.chanCache.pop(id); ok {
			chans.close()
		}
		return ctx.Err()
	case err := <-chans.errCh:
		return err
	case data := <-chans.ch:
		json.Unmarshal(data, model)
		return nil
	}
}

func (client *Client) sendByID(method string, params map[string]interface{}, public bool) (*requestsChans, error) {
	if !client.publicWSManager.isOpen && public {
		return nil, fmt.Errorf("Websocket for public methods closed")
	}

	if !client.tradingWSManager.isOpen && !public {
		return nil, fmt.Errorf("Websocket for trading methods closed")
	}
	chans := makeChans()
	id := client.chanCache.store(chans)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}
	fmt.Println("sending notification", notification)
	data, _ := json.Marshal(notification)
	if public {
		client.publicWSManager.snd <- data
	} else { // trading
		client.tradingWSManager.snd <- data
	}
	return chans, nil
}

func (client *Client) handle(rcvCh chan []byte) {
	for data := range rcvCh {
		resp := wsResponse{}
		json.Unmarshal(data, &resp)
		if resp.ID != 0 {
			if chans, ok := client.chanCache.pop(resp.ID); ok {
				defer chans.close()
				if resp.Error != nil {
					chans.errCh <- fmt.Errorf("Error:%s", resp.Error)
					continue
				}
				chans.ch <- data
			}
		} else if resp.Method != "" {
			fmt.Println("method ", resp.Method)
			key := buildKey(resp.Method, resp.Params.Symbol, resp.Params.Period)
			// if resp.Method == methodUpdateOrderbook && client.subsCache.waitingOrderbookSnapshot(key) {
			// 	return
			// }
			client.subsCache.update(resp.Method, key, data)
			// if resp.Method == methodUpdateOrderbook && client.subsCache.requireOrderbookSnapshot(key) {
			// 	go func() {
			// 		chans := makeChans()
			// 		client.sendByID(methodSubscribeOrderbook, map[string]interface{}{"symbol": resp.Params.Symbol}, chans, publicCall)
			// 		var err error
			// 		select {
			// 		case err = <-chans.errCh:
			// 		case <-chans.ch:
			// 		}
			// 		if err != nil {
			// 			fmt.Println("Filed to rebase orderbook, freeing orderbook asociated memory")
			// 			client.subsCache.closeFeedOrderbook(key)
			// 		}
			// 	}()
			// }
		}
	}
}

func buildKey(method, symbol, period string) string {
	methodKey := methodMapping[method]
	var key string
	if method == methodReport || method == methodActiveOrders {
		key = methodKey + "::"
	} else {
		key = methodKey + ":" + symbol + ":" + period
	}
	return strings.ToUpper(key)
}


func (client *Client) getCurrencies(ctx context.Context) ([]models.Currency, error) {
	var resp struct {
		Result []models.Currency
	}
	err := client.doPublicRequest(ctx, methodGetCurrencies, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}


func (client *Client) getChanCurrencies(ctx context.Context) (<-chan []models.Currency, <-chan error) {
	chans, err := client.sendByID(methodGetCurrencies, nil, publicCall)
	resultCh := make(chan []models.Currency)
	errCh := make(chan error)
	if err != nil {
		go func() {
			defer close(errCh)
			defer close(resultCh)
			errCh <- err
		}()
	} else {
		go func() {
			defer close(errCh)
			select {
			case err:=<-chans.errCh:
				close(errCh)
				errCh<-err
			case <-ctx.Done():
					close(errCh)
					close(resultCh)
					errCh<-ctx.Err()

			}
		}

	}
	var resp struct {
		Result []models.Currency
	}
	err := client.doPublicRequest(ctx, methodGetCurrencies, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}



func (client *Client) getCurrency(ctx context.Context, arguments ...args.Argument) (*models.Currency, error) {
	var resp struct {
		Result models.Currency
	}
	err := client.doPublicRequest(ctx, methodGetCurrency, arguments, []string{"currency"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

func (client *Client) getSymbols(ctx context.Context) ([]models.Symbol, error) {
	var resp struct {
		Result []models.Symbol
	}
	err := client.doPublicRequest(ctx, methodGetSymbols, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (client *Client) getSymbol(ctx context.Context, arguments ...args.Argument) (*models.Symbol, error) {
	var resp struct {
		Result models.Symbol
	}
	err := client.doPublicRequest(ctx, methodGetSymbol, arguments, []string{"symbol"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

func (client *Client) getTrades(ctx context.Context, arguments ...args.Argument) (trades []models.PublicTrade, symbol string, err error) {
	var resp struct {
		Result struct {
			Symbol string
			Data   []models.PublicTrade
		}
	}
	err = client.doPublicRequest(ctx, methodGetTrades, arguments, []string{"symbol"}, &resp)
	if err != nil {
		return nil, "", err
	}
	return resp.Result.Data, resp.Result.Symbol, nil
}

func (client *Client) authenticate(ctx context.Context) (err error) {
	nonce := makeNonce(30)
	h := hmac.New(sha256.New, []byte(client.apiSecret))
	h.Write([]byte(nonce))
	signature := hex.EncodeToString(h.Sum(nil))
	params := map[string]interface{}{
		"algo":      "HS256",
		"pKey":      client.apiKey,
		"nonce":     nonce,
		"signature": signature,
	}
	chans := makeChans()
	err = client.sendByID("login", params, chans, tradingCall)
	if err != nil {
		chans.close()
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chans.errCh:
		return err
	case <-chans.ch:
		return nil
	}
}

func (client *Client) getTradingBalance(ctx context.Context) ([]models.Balance, error) {
	var resp struct {
		Result []models.Balance
	}
	err := client.doTradingRequest(ctx, "getTradingBalance", nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (client *Client) getOrders(ctx context.Context) ([]models.Report, error) {
	var resp struct {
		Result []models.Report
	}
	err := client.doTradingRequest(ctx, "getOrders", nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (client *Client) createOrder(ctx context.Context, arguments ...args.Argument) (*models.Report, error) {
	var resp struct {
		Result models.Report
	}
	err := client.doTradingRequest(ctx, "newOrder", arguments, []string{"clientOrderId", "symbol", "side", "quantity"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

func (client *Client) cancelOrder(ctx context.Context, arguments ...args.Argument) (*models.Report, error) {
	var resp struct {
		Result models.Report
	}
	err := client.doTradingRequest(ctx, "cancelOrder", arguments, []string{"clientOrderId"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

func (client *Client) replaceOrder(ctx context.Context, arguments ...args.Argument) (*models.Report, error) {
	var resp struct {
		Result models.Report
	}
	err := client.doTradingRequest(ctx, "cancelReplaceOrder", arguments, []string{"clientOrderId", "requestClientId", "price", "quantity"}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

func (client *Client) subscribeToTicker(ctx context.Context, arguments ...args.Argument) (chan models.Ticker, error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return nil, err
	}
	key := buildKey(methodSubscribeTicker, params["symbol"].(string), "")
	buffSize := 0
	if limit, ok := params["limit"]; ok {
		buffSize = limit.(int)
	}
	feedCh := client.subsCache.makeFeedTicker(key, buffSize)
	chans := makeChans()
	err = client.sendByID(methodSubscribeTicker, params, chans, publicCall)
	if err != nil {
		chans.close()
		return nil, err
	}
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-chans.errCh:
	case <-chans.ch:
	}
	if err != nil {
		client.subsCache.closeFeedTicker(key)
		return nil, err
	}
	return feedCh, nil
}

func (client *Client) unsubscribeToTicker(ctx context.Context, arguments ...args.Argument) error {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return err
	}
	key := buildKey(methodUnsubcribeTicker, params["symbol"].(string), "")
	client.subsCache.closeFeedTicker(key)
	chans := makeChans()
	err = client.sendByID(methodUnsubcribeTicker, params, chans, publicCall)
	if err != nil {
		chans.close()
		return err
	}
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-chans.errCh:
	case <-chans.ch:
	}
	return err
}

func (client *Client) subscribeToOrderbooks(ctx context.Context, arguments ...args.Argument) (chan models.OrderBook, error) {
	params, err := args.BuildParams(arguments, "symbol")
	if err != nil {
		return nil, err
	}
	key := buildKey(methodSubscribeOrderbook, params["symbol"].(string), "")

	buffSize := 0
	if limit, ok := params["limit"]; ok {
		buffSize = limit.(int)
	}
	feedCh := client.subsCache.makeFeedOrderbook(key, buffSize)
	chans := makeChans()
	err = client.sendByID(methodSubscribeOrderbook, params, chans, publicCall)
	if err != nil {
		chans.close()
		return nil, err
	}
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-chans.errCh:
	case <-chans.ch:
	}
	if err != nil {
		client.subsCache.closeFeedOrderbook(key)
		return nil, err
	}
	return feedCh, nil
}
