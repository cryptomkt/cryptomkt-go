package rest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	apiURL                  = "https://api.exchange.cryptomkt.com"
	apiVersion              = "/api/3/"
	headerContentType       = "Content-type"
	headerUserAgent         = "User-Agent"
	userAgentCryptomarketGo = "cryptomarket/go"
	applicationJson         = "application/json"
	applicationUrlEncoded   = "application/x-www-form-urlencoded"
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

type RequestData struct {
	cxt                context.Context
	method             string
	endpoint           string
	urlEncodedPayload  string
	jsonEncodedPayload string
	public             bool
}

func (hclient httpclient) makeRequest(requestData *RequestData) (result []byte, err error) {
	request, err := hclient.buildRequest(requestData)
	if err != nil {
		return nil, err
	}
	response, err := hclient.client.Do(request)
	if err != nil {
		return nil, errors.New("CryptomarketSDKError: Can't make the request: " + err.Error())
	}
	defer response.Body.Close()
	return readResponse(response)
}

func (hclient httpclient) buildRequest(requestData *RequestData) (*http.Request, error) {
	if requestData.method == methodPost {
		return hclient.buildPostRequest(requestData)
	}
	if requestData.method == methodGet {
		return hclient.buildGetRequest(requestData, requestData.urlEncodedPayload)
	}
	return hclient.buildOtherRequest(requestData, requestData.urlEncodedPayload)
}

func (hclient httpclient) buildGetRequest(requestData *RequestData, rawQuery string) (*http.Request, error) {
	request, err := hclient.buildRequestHelper(requestData, nil, rawQuery)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = rawQuery
	return request, nil
}

func (hclient httpclient) buildPostRequest(requestData *RequestData) (*http.Request, error) {
	request, err := hclient.buildRequestHelper(requestData, strings.NewReader(requestData.jsonEncodedPayload), requestData.jsonEncodedPayload)
	if err != nil {
		return nil, err
	}
	request.Header.Add(headerContentType, applicationJson)
	return request, err
}

func (hclient httpclient) buildOtherRequest(requestData *RequestData, urlEncodedQuery string) (*http.Request, error) {
	request, err := hclient.buildRequestHelper(requestData, strings.NewReader(urlEncodedQuery), urlEncodedQuery)
	if err != nil {
		return nil, err
	}
	request.Header.Add(headerContentType, applicationUrlEncoded)
	return request, err
}

func (hclient httpclient) buildRequestHelper(requestData *RequestData, body io.Reader, query string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(requestData.cxt, requestData.method, apiURL+apiVersion+requestData.endpoint, body)
	if err != nil {
		return nil, errors.New("CryptomarketSDKError: Can't build the request: " + err.Error())
	}
	request.Header.Add(headerUserAgent, userAgentCryptomarketGo)
	if !requestData.public {
		request.Header.Add("Authorization", hclient.getCredentialForRequest(request, query))
	}
	return request, nil
}

func readResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("CryptomarketSDKError: Can't read the response body: " + err.Error())
	}
	return body, nil
}
