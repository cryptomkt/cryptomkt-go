package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cryptomarket/cryptomarket-go/args"
	"github.com/cryptomarket/cryptomarket-go/models"
)

type parseNotificationFnType func([]byte)

func (client *clientBase) doRequest(
	ctx context.Context,
	method string,
	arguments []args.Argument,
	requiredArguments []string,
	model interface{},
) error {
	return client.doRequestOfNNotifications(ctx, method, arguments, requiredArguments, unmarshalDataToModel(model), "")
}

func unmarshalDataToModel(model interface{}) parseNotificationFnType {
	return func(data []byte) { json.Unmarshal(data, model) }
}

func (client *clientBase) doRequestOfNNotifications(
	ctx context.Context,
	method string,
	arguments []args.Argument,
	requiredArguments []string,
	parseNotificationFn parseNotificationFnType,
	nNotificationsParam string,
) error {
	params, err := args.BuildParams(arguments, requiredArguments...)
	if err != nil {
		return err
	}
	nNotifications := getNOfNotifications(params, nNotificationsParam)
	if !client.wsManager.isOpen {
		return fmt.Errorf("CryptomarketSDKError: websocket connection is closed")
	}
	ch := make(chan []byte, 1)
	id := client.chanCache.saveCh(ch, nNotifications)
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
	client.sendRequest(data)
	return client.handleResponse(ctx, id, ch, nNotifications, parseNotificationFn)
}

func (client *clientBase) sendRequest(data []byte) {
	client.wsManager.snd <- data
}

func getNOfNotifications(params map[string]interface{}, nNotificationsParam string) int {
	if len(nNotificationsParam) == 0 {
		return 1
	}
	param, ok := params[nNotificationsParam]
	if !ok {
		return 1
	}
	if reflect.TypeOf(param).Kind() != reflect.Slice {
		return 1
	}
	return reflect.ValueOf(param).Len()
}

func (client *clientBase) handleResponse(ctx context.Context, id int64, ch chan []byte, nNotifications int, updateModelAggFn parseNotificationFnType) error {
	for i := 0; i < nNotifications; i++ {
		select {
		case <-ctx.Done():
			client.chanCache.closeAndRemoveCh(id)
			return ctx.Err()
		case data := <-ch:
			var resp struct {
				Error *models.APIError
			}
			json.Unmarshal(data, &resp)
			if resp.Error != nil {
				return fmt.Errorf("CryptomarketAPIError: %v", resp.Error)
			}
			updateModelAggFn(data)
		}
	}
	return nil
}
