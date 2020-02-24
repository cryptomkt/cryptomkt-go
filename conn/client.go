// Package client implements a client to connect with CryptoMarket using
// the endpoints given at https://developers.cryptomkt.com/.
package conn

import (
	"bytes"
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/requests"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

var (
	// DELAY is the amount to wait in seconds between requests to the server,
	// too many requests and the ip is blocked.
	DELAY float64 = 2.5

	// Version of the api used to connect with crypto market.
	apiVersion = "v1"

	// URI to connect with crypto market api.
	baseApiUri = "https://api.cryptomkt.com/"
)

// Client keep the needed data to connect with the asociated CryptoMarket account.
type Client struct {
	auth       *HMACAuth
	httpClient *http.Client
}

// New builds a new client and returns a pointer to it.
func NewClient(apiKey, apiSecret string) *Client {
	client := &Client{
		auth:       newAuth(apiKey, apiSecret),
		httpClient: &http.Client{},
	}
	return client
}

// runRequest makes the builded http request to cryptoMarket,
// and read the response
func (client *Client) runRequest(httpReq *http.Request) ([]byte, error) {
	resp, err := client.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}
	return respBody, nil
}

// getPublic makes an http request to a given enpoint, given a custom request that contains
// the needed arguments
func (client *Client) getPublic(endpoint string, request *requests.Request) ([]byte, error) {
	args := request.GetArguments()
	u, err := url.Parse(baseApiUri)
	if err != nil {
		return nil, fmt.Errorf("Error parsing url %s: %v", baseApiUri, err)
	}
	u.Path = path.Join(u.Path, apiVersion, endpoint)
	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building NewRequest struct: %v", err)
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

// get comunicates to Cryptomarket via the http get method, also set
// the needed headers of the request for an authenticated communication.
func (client *Client) get(endpoint string, request *requests.Request) ([]byte, error) {
	args := request.GetArguments()
	u, err := url.Parse(baseApiUri)
	if err != nil {
		return nil, fmt.Errorf("Error parsing url %s: %v", baseApiUri, err)
	}
	u.Path = path.Join(u.Path, apiVersion, endpoint)
	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error building NewRequest struct: %v", err)
	}
	// query the Arguments in the http request, if there are Arguments
	if len(args) != 0 {
		q := httpReq.URL.Query()
		for k, v := range args {
			q.Add(k, v)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	requestPath := "/" + apiVersion + "/" + endpoint
	client.auth.setHeaders(httpReq, requestPath, "")

	return client.runRequest(httpReq)
}

// post comunicates to Cryptomarket via the http post method, also set
// the needed headers of the request for an authenticated communication.
func (client *Client) post(endpoint string, request *requests.Request) ([]byte, error) {
	args := request.GetArguments()

	u, err := url.Parse(baseApiUri)
	if err != nil {
		return nil, fmt.Errorf("Error parsing url %s: %v", baseApiUri, err)
	}
	u.Path = path.Join(u.Path, apiVersion, endpoint)

	// builds a form from the Arguments
	form := url.Values{}
	for k, v := range args {
		form.Add(k, v)
	}
	httpReq, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Error building NewRequest struct: %v", err)
	}

	//sets the body for the header, arguments must be sorted
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var bb bytes.Buffer
	for _, k := range keys {
		bb.WriteString(args[k])
	}

	requestPath := "/" + apiVersion + "/" + endpoint
	client.auth.setHeaders(httpReq, requestPath, bb.String())

	//required header for the reciever to interpret the request as a http form post
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	return client.runRequest(httpReq)
}

// makeReq builds a request to ensure the presence of required arguments, stores the
// arguments in its string form.
func makeReq(required []string, args ...args.Argument) (*requests.Request, error) {
	req := requests.NewReq(required)
	for _, argument := range args {
		err := argument(req)
		if err != nil {
			return nil, fmt.Errorf("argument error: %s", err)
		}
	}
	err := req.AssertRequired()
	if err != nil {
		return nil, fmt.Errorf("required arguments not meeted:%s", err)
	}
	return req, nil
}

// postReq builds a post request and send it to CryptoMarket.
// Returns a []byte with the response
func (client *Client) postReq(endpoint string, caller string, required []string, args ...args.Argument) ([]byte, error) {
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in %s: %s", caller, err)
	}
	return client.post(endpoint, req)
}

// postReq builds a getReq request and send it to CryptoMarket.
// Returns a []byte with the response
func (client *Client) getReq(endpoint string, caller string, required []string, args ...args.Argument) ([]byte, error) {
	req, err := makeReq(required, args...)
	if err != nil {
		return nil, fmt.Errorf("Error in %s: %s", caller, err)
	}
	return client.get(endpoint, req)
}
