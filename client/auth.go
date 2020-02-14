package client

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type HMACAuth struct {
	apiKey    string
	apiSecret string
}

func NewAuth(apiKey, apiSecret string) (*HMACAuth, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("client: api key can't be empty")
	}
	if apiSecret == "" {
		return nil, fmt.Errorf("client: api secret can't be empty")
	}
	auth := &HMACAuth{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
	return auth, nil
}

func (auth *HMACAuth) SetHeaders(req *http.Request, endpoint string, body string) {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%v%s%s", timestamp, endpoint, body)
	h := hmac.New(sha512.New384, []byte(auth.apiSecret))
	h.Write([]byte(data))

	req.Header.Add("X-MKT-APIKEY", auth.apiKey)
	req.Header.Add("X-MKT-SIGNATURE", hex.EncodeToString(h.Sum(nil)))
	req.Header.Add("X-MKT-TIMESTAMP", fmt.Sprintf("%v", timestamp))
}
