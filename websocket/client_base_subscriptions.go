package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/models"
)

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
	id := client.chanCache.saveCh(ch, 1)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}

	data, err := json.Marshal(notification)
	if err != nil {
		client.chanCache.closeAndRemoveCh(id)
		return nil, fmt.Errorf("CryptomarketSDKError: invalid notification: %v", err)
	}
	key := subscriptionMapping[method]
	dataOut := make(chan []byte, 1)
	client.chanCache.saveSubscriptionCh(key, dataOut)
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
	id := client.chanCache.saveCh(ch, 1)
	notification := wsNotification{
		ID:     id,
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		client.chanCache.closeAndRemoveCh(id)
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
