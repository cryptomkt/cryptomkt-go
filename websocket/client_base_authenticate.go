package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cryptomkt/cryptomkt-go/models"
)

func (client *clientBase) authenticate(apiKey, apiSecret string) (err error) {
	if !client.wsManager.isOpen {
		return fmt.Errorf("CryptomarketSDKError: websocket connection closed")
	}
	intTimestamp := time.Now().UnixMilli()
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
	id := client.chanCache.saveCh(ch, 1)
	notification := wsNotification{
		ID:     id,
		Method: "login",
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
