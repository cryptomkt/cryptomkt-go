package	main

import (
    "crypto/hmac"
    "crypto/sha512"
	"encoding/hex"
	"net/http"
	"time"
	"fmt"
)
type HMACAuth struct {
	apiKey string
	apiSecret string
	apiVersion string
}

func New(apiKey, apiSecret, apiVersion string) (HMACAuth) {
	return HMACAuth{
		apiKey:apiKey,
		apiSecret:apiSecret,
		apiVersion:apiVersion,
	}
}

func (auth *HMACAuth) SetHeaders(req *http.Request) (error) {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%v",timestamp) + "/"+ auth.apiVersion +"/account"
    h := hmac.New(sha512.New384, []byte(auth.apiSecret))
    h.Write([]byte(data))

	req.Header.Add("X-MKT-APIKEY", auth.apiKey)
	req.Header.Add("X-MKT-SIGNATURE",  hex.EncodeToString(h.Sum(nil)))
	req.Header.Add("X-MKT-TIMESTAMP", fmt.Sprintf("%v",timestamp))
	return nil
}