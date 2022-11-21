package rest

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cryptomarket/cryptomarket-go/args"
)

var (
	apiURL     = "https://api.exchange.cryptomkt.com"
	apiVersion = "/api/3/"
)

// httpclient handles all the http logic, leaving public only whats needed.
// accepts Get, Post, Put and Delete functions, all with parameters and return
// the response bytes
type httpclient struct {
	client    *http.Client
	apiKey    string
	apiSecret string
	window    int
}

// New creates a new httpclient
func newHTTPClient(apiKey, apiSecret string, window int) httpclient {
	return httpclient{
		client:    &http.Client{},
		apiKey:    apiKey,
		apiSecret: apiSecret,
		window:    window,
	}
}

func (hclient httpclient) doRequest(cxt context.Context, method, endpoint string, params map[string]interface{}, public bool) (result []byte, err error) {
	// build query
	rawQuery := args.BuildQuery(params)
	// build request
	var req *http.Request
	if method == methodGet {
		req, err = http.NewRequestWithContext(cxt, method, apiURL+apiVersion+endpoint, nil)
		req.URL.RawQuery = rawQuery
	} else {
		req, err = http.NewRequestWithContext(cxt, method, apiURL+apiVersion+endpoint, strings.NewReader(rawQuery))
	}
	if err != nil {
		return nil, errors.New("CryptomarketSDKError: Can't build the request: " + err.Error())
	}
	req.Header.Add("User-Agent", "cryptomarket/go")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	// add auth header if is not a public call
	if !public {
		req.Header.Add("Authorization", hclient.buildCredential(method, endpoint, rawQuery))
	}

	// make request
	resp, err := hclient.client.Do(req)
	if err != nil {
		return nil, errors.New("CryptomarketSDKError: Can't make the request: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("CryptomarketSDKError: Can't read the response body: " + err.Error())
	}
	return body, nil
}

func (hclient httpclient) buildCredential(httpMethod, method, query string) string {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	msg := httpMethod + apiVersion + method
	if len(query) != 0 {
		if httpMethod == methodGet {
			msg += "?"
		}
		msg += query
	}
	msg += timestamp
	if hclient.window != 0 {
		msg += strconv.FormatInt(int64(hclient.window), 10)
	}
	h := hmac.New(sha256.New, []byte(hclient.apiSecret))
	h.Write([]byte(msg))
	signature := hex.EncodeToString(h.Sum(nil))
	str := hclient.apiKey + ":" + signature + ":" + timestamp
	if hclient.window != 0 {
		str += (":" + strconv.FormatInt(int64(hclient.window), 10))
	}
	return "HS256 " + base64.StdEncoding.EncodeToString([]byte(str))
}
