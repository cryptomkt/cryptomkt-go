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
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	apiURL     = "https://api.exchange.cryptomkt.com"
	apiVersion = "/api/2/"
)

// httpclient handles all the http logic, leaving public only whats needed.
// accepts Get, Post, Put and Delete functions, all with parameters and return
// the response bytes
type httpclient struct {
	client    *http.Client
	apiKey    string
	apiSecret string
}

// New creates a new httpclient
func newHTTPClient(apiKey, apiSecret string) httpclient {
	return httpclient{
		client:    &http.Client{},
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (hclient httpclient) doRequest(cxt context.Context, method, endpoint string, params map[string]interface{}, public bool) (result []byte, err error) {
	// build query
	rawQuery := buildQuery(params)
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
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msg := httpMethod + timestamp + apiVersion + method
	if len(query) != 0 {
		if httpMethod == methodGet {
			msg += "?"
		}
		msg += query
	}
	h := hmac.New(sha256.New, []byte(hclient.apiSecret))
	h.Write([]byte(msg))
	signature := hex.EncodeToString(h.Sum(nil))
	return "HS256 " + base64.StdEncoding.EncodeToString([]byte(hclient.apiKey+":"+timestamp+":"+signature))
}

func buildQuery(params map[string]interface{}) string {
	query := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case []string:
			strs := strings.Join(v, ",")
			query.Add(key, strs)
		case string:
			query.Add(key, v)
		case int:
			query.Add(key, strconv.Itoa(v))
		}
	}
	return query.Encode()
}
