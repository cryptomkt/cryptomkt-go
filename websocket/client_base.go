package websocket

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

type clientBase struct {
	wsManager *wsManager
	chanCache *chanCache
	window    int
}

// Close close all the channels related to the client as well as the websocket connection.
// trying to make requests over a closed client will result in error.
func (client *clientBase) Close() {
	client.chanCache.close()
	client.wsManager.close()
}

func (client *clientBase) handle(rcvCh chan []byte) {
	for data := range rcvCh {
		resp := wsResponse{}
		json.Unmarshal(data, &resp)
		if key, ok := getSubscriptionKey(resp); ok {
			client.chanCache.sendViaSubscriptionCh(key, data)
			continue
		}
		if ch, ok := client.chanCache.pop(resp.ID); ok {
			defer close(ch)
			ch <- data
			continue
		}
	}
}

func (client *clientBase) doRequest(
	ctx context.Context,
	method string,
	arguments []args.Argument,
	requiredArguments []string,
	model interface{},
) error {
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
			Error *models.APIError
		}
		json.Unmarshal(data, &resp)
		if resp.Error != nil {
			return fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
		}
		json.Unmarshal(data, model)
		return nil
	}
}

func (client *clientBase) doSubscription(
	method string,
	arguments []args.Argument,
	requiredArguments []string,
) (chan []byte, error) {
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
	key := subscriptionMapping[method]
	dataOut := make(chan []byte, 1)
	client.chanCache.storeSubscriptionCh(key, dataOut)
	client.wsManager.snd <- data
	data = <-ch
	var resp struct {
		Error *models.APIError
	}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		close(dataOut)
		return nil, fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
	}
	return dataOut, nil
}

func (client *clientBase) doUnsubscription(
	method string,
	arguments []args.Argument,
	requiredArguments []string,
) error {
	params, err := args.BuildParams(arguments, requiredArguments...)
	if err != nil {
		return err
	}
	if !client.wsManager.isOpen {
		return fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	key := subscriptionMapping[method]
	client.chanCache.deleteSubscriptionCh(key)
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
		Error *models.APIError
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
	intTimestamp := time.Now().Unix() * 1000
	timestamp := strconv.FormatInt(intTimestamp, 10)
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(timestamp))
	if client.window != 0 {
		h.Write([]byte(fmt.Sprint(client.window)))
	}
	signature := hex.EncodeToString(h.Sum(nil))
	params := map[string]interface{}{
		"type":      "HS256",
		"api_key":   apiKey,
		"timestamp": intTimestamp,
		"signature": signature,
	}
	if client.window != 0 {
		params["window"] = client.window
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
		Error *models.APIError
	}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		return fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
	}
	return nil
}
