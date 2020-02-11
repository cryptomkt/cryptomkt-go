package client

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"net/url"
	"path"
)

type Client struct {
	apiVersion string
	baseApiUri string
	auth *HMACAuth
	httpClient *http.Client
	
}

func NewClient(apiKey, apiSecret string) (*Client, error) {
	apiVersion := "v1"
	baseApiUri := "https://api.cryptomkt.com/"

	auth, err := NewAuth(apiKey, apiSecret)
	if err != nil {
		fmt.Println("error with the api key or the api secert")
		os.Exit(1)
	}

	client := &Client{
		baseApiUri: baseApiUri,
		apiVersion: apiVersion,
		auth: auth,
		httpClient: &http.Client{},
	}
	return client, nil
}

func (client *Client) get(endpoint string, args map[string]string) (string) {
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		fmt.Println("could not parse the base api uri", client.baseApiUri)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)	
	req, err := http.NewRequest("GET", u.String(), nil)
	if err!=nil {
		fmt.Println(err)
	}
	q := req.URL.Query()
	for k, v := range args {
		q.Add(k,v)
	}
	req.URL.RawQuery = q.Encode()
	requestPath := fmt.Sprintf("/%s/%s", client.apiVersion, endpoint)
	client.auth.SetHeaders(req, requestPath)
	
	resp, err:= client.httpClient.Do(req) 
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(body)
}

func (client *Client) getAccount() (string) {
	return client.get("account", map[string]string{})
}

func (client *Client) getBalance() (string) {
	return client.get("balance", map[string]string{})
}

func (client *Client) getWallet() (string) {
	return client.get("balance", map[string]string{})
}

func (client *Client) getTransactions() (string) {
	return client.get("transactions", map[string]string{})
}

func (client *Client) getActiveOrders(args map[string]string) (string) {
	return client.get("orders/active", args)
}


func mockTrades() {
	args := map[string]string {
		"market": "ETHCLP",
		"end":"2018-06-06",
		"page":"2",
		"limit":"10",
	}
	urlandargs := "https://api.cryptomkt.com/v1/trades" + parseMapss(args)
	fmt.Println(urlandargs)
	req, err := http.NewRequest("GET", urlandargs, nil)
	if err!=nil {
		fmt.Println(err)
	}
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}
