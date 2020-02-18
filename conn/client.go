// Package client implements a client to connect with CryptoMarket,
// using the endpoints given at https://developers.cryptomkt.com/
package conn

import (
	"bytes"
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/requests"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

// Client keep the needed information to connect with the asociated CryptoMarket account.
type Client struct {
	apiVersion string
	baseApiUri string
	auth       *HMACAuth
	httpClient *http.Client
}

// New builds a new client and returns a pointer to it.
// It can fail if the api key or the api secret are empty
func NewClient(apiKey, apiSecret string) *Client {
	apiVersion := "v1"
	baseApiUri := "https://api.cryptomkt.com/"
	auth := newAuth(apiKey, apiSecret)

	client := &Client{
		baseApiUri: baseApiUri,
		apiVersion: apiVersion,
		auth:       auth,
		httpClient: &http.Client{},
	}
	return client
}

func (client *Client) runRequest(httpReq *http.Request) (string, error) {
	resp, err := client.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}
	return string(respBody), nil
}

func (client *Client) getPublic(endpoint string, request *requests.Request) (string, error) {
	args := request.GetArguments()
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		return "", fmt.Errorf("Error parsing url %s: %v", client.baseApiUri, err)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)
	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("Error building NewRequest struct: %v", err)
	}
	if len(args) != 0 {
		q := httpReq.URL.Query()
		for k, v := range args {
			q.Add(k, v)
		}
		httpReq.URL.RawQuery = q.Encode()
	}
	return client.runRequest(httpReq)
}

// get comunicates to Cryptomarket via the http get method
// Its the base implementation which the public methods use.
func (client *Client) get(endpoint string, request *requests.Request) (string, error) {
	args := request.GetArguments()
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		return "", fmt.Errorf("Error parsing url %s: %v", client.baseApiUri, err)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)
	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("Error building NewRequest struct: %v", err)
	}
	// query the Arguments in the http request, if there are Arguments
	if len(args) != 0 {
		q := httpReq.URL.Query()
		for k, v := range args {
			q.Add(k, v)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	requestPath := "/" + client.apiVersion + "/" + endpoint
	client.auth.setHeaders(httpReq, requestPath, "")

	return client.runRequest(httpReq)
}

// post comunicates to Cryptomarket via the http post method.
// Its the base implementation which the public methods use.
// Arguments are required.
func (client *Client) post(endpoint string, request *requests.Request) (string, error) {
	args := request.GetArguments()

	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		return "", fmt.Errorf("Error parsing url %s: %v", client.baseApiUri, err)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)

	// builds a form from the Arguments
	form := url.Values{}
	for k, v := range args {
		form.Add(k, v)
	}
	httpReq, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("Error building NewRequest struct: %v", err)
	}

	//sets the body for the header
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var bb bytes.Buffer
	for _, k := range keys {
		bb.WriteString(args[k])
	}

	requestPath := "/" + client.apiVersion + "/" + endpoint
	client.auth.setHeaders(httpReq, requestPath, bb.String())

	//required header for the reciever to interpret the request as a http form post
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	return client.runRequest(httpReq)
}
