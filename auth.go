package	client

import (
    "crypto/hmac"
    "crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"
	"fmt"
	"errors"
)
type HMACAuth struct {
	apiKey string
	apiSecret string
}

func NewAuth(apiKey, apiSecret string) (*HMACAuth, error) {
	if apiKey == "" {
		return nil, errors.New("api key can't be empty")

	}
	if apiSecret == ""{
		return nil, errors.New("api secret can't be empty")
	}
	auth := &HMACAuth{
		apiKey:apiKey,
		apiSecret:apiSecret,
	}
	return auth, nil
}

func (auth *HMACAuth) SetHeaders(req *http.Request, endpoint string) {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%v%s",timestamp, endpoint)
    h := hmac.New(sha512.New384, []byte(auth.apiSecret))
    h.Write([]byte(data))

	req.Header.Add("X-MKT-APIKEY", auth.apiKey)
	req.Header.Add("X-MKT-SIGNATURE",  hex.EncodeToString(h.Sum(nil)))
	req.Header.Add("X-MKT-TIMESTAMP", fmt.Sprintf("%v",timestamp))
}