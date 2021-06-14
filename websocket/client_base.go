package websocket

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cryptomarket/cryptomarket-go/args"
)

const (
	publicCall  bool = true
	tradingCall bool = false
)

type clientBase struct {
	wsManager            *wsManager
	chanCache            *chanCache
	subscriptionKeysFunc func(string, map[string]interface{}) (string, bool)
	keyFromResponse      func(wsResponse) string
}

// Close close all the channels related to the client as well as the websocket connection.
// trying to make requests over a closed client will result in error.
func (client *clientBase) Close() {
	client.wsManager.close()
	client.chanCache.close()
}

func (client *clientBase) handle(rcvCh chan []byte) {
	for data := range rcvCh {
		resp := wsResponse{}
		json.Unmarshal(data, &resp)
		if resp.ID != 0 {
			if ch, ok := client.chanCache.pop(resp.ID); ok {
				defer close(ch)
				ch <- data
			}
		} else if resp.Method != "" {
			key := client.keyFromResponse(resp)
			if feedCh, ok := client.chanCache.getSubcriptionCh(key); ok {
				feedCh <- data
			}
		}
	}
}

func (client *clientBase) buildKeyFromResponse(response wsResponse) string {
	methodKey := methodMapping[response.Method]
	var key string
	if response.Method == "report" || response.Method == "ActiveOrders" {
		key = methodKey + "::"
	} else {
		period := response.Params.Period
		if methodKey == "candles" && period == "" { // default period
			period = string(args.PeriodType30Minutes)
		}
		key = methodKey + ":" + response.Params.Symbol + ":" + period
	}
	return strings.ToUpper(key)
}

func (client *clientBase) buildKey(method string, params map[string]interface{}) string {
	if key, ok := client.subscriptionKeysFunc(method, params); ok {
		return key
	}
	return "subscription"
}

func (client *clientBase) doRequest(ctx context.Context, method string, arguments []args.Argument, requiredArguments []string, model interface{}) error {
	params, err := args.BuildParams(arguments, requiredArguments...)
	if err != nil {
		return err
	}
	if !client.wsManager.isOpen {
		return fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.store(ch)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	client.wsManager.snd <- data
	select {
	case <-ctx.Done():
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return ctx.Err()
	case data := <-ch:
		var resp struct {
			Error APIError
		}
		json.Unmarshal(data, &resp)
		if resp.Error != nil {
			return fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
		}
		json.Unmarshal(data, model)
		return nil
	}
}

func (client *clientBase) doSubscription(method string, arguments []args.Argument, requiredArguments []string) (chan []byte, error) {
	params, err := args.BuildParams(arguments, requiredArguments...)
	if err != nil {
		return nil, err
	}
	if !client.wsManager.isOpen {
		return nil, fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.store(ch)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return nil, fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	key := client.buildKey(method, params)
	dataOut := make(chan []byte, 1)
	client.chanCache.storeSubscriptionCh(key, dataOut)
	client.wsManager.snd <- data
	data = <-ch
	var resp struct {
		Error APIError
	}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		close(dataOut)
		return nil, fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
	}
	return dataOut, nil
}

func (client *clientBase) doUnsubscription(method string, arguments []args.Argument, requiredArguments []string) error {
	params, err := args.BuildParams(arguments, requiredArguments...)
	if err != nil {
		return err
	}
	if !client.wsManager.isOpen {
		return fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	key := client.buildKey(method, params)
	if ch, ok := client.chanCache.getSubcriptionCh(key); ok {
		client.chanCache.deleteSubscriptionCh(key)
		close(ch)
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.store(ch)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	client.wsManager.snd <- data
	data = <-ch
	var resp struct {
		Error APIError
	}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		return fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
	}
	return nil
}

func (client *clientBase) authenticate(apiKey, apiSecret string) (err error) {
	if !client.wsManager.isOpen {
		return fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	nonce := makeNonce(30)
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(nonce))
	signature := hex.EncodeToString(h.Sum(nil))
	params := map[string]interface{}{
		"algo":      "HS256",
		"pKey":      apiKey,
		"nonce":     nonce,
		"signature": signature,
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.store(ch)
	notification := wsNotification{
		ID:     id,
		Method: "login",
		Params: params,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		if ch, ok := client.chanCache.pop(id); ok {
			close(ch)
		}
		return fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	client.wsManager.snd <- data
	data = <-ch
	var resp struct {
		Error APIError
	}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		return fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
	}
	return nil
}
