package conn

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"net/http"
	"time"
)

// A HMACAuth keeps the keys of a client
// making them usable along all request of the client to crpytomkt
type HMACAuth struct {
	apiKey    string
	apiSecret string
}

// newAuth creates a new HMACAuth
func newAuth(apiKey, apiSecret string) (*HMACAuth) {
	auth := &HMACAuth{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
	return auth
}

// setHeaders set the X-MKT-APIKEY, X-MKT-SIGNATURE and the X-MKT-TIMESTAMP headers.
// https://developers.cryptomkt.com/es/#api-key
func (auth *HMACAuth) setHeaders(req *http.Request, endpoint string, body string) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	data := timestamp + endpoint + body
	h := hmac.New(sha512.New384, []byte(auth.apiSecret))
	h.Write([]byte(data))

	req.Header.Add("X-MKT-APIKEY", auth.apiKey)
	req.Header.Add("X-MKT-SIGNATURE", hex.EncodeToString(h.Sum(nil)))
	req.Header.Add("X-MKT-TIMESTAMP", timestamp)
}
