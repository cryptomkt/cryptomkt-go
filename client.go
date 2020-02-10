package main

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

	auth, err := NewAuth(apiKey, apiSecret, apiVersion)
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

func (client *Client) get(endpoint string) (string) {
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		fmt.Println("could not parse the base api uri", client.baseApiUri)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)	

	req, err := http.NewRequest("GET", u.String(), nil)
	if err!=nil {
		fmt.Println(err)
	}
	
	err = client.auth.SetHeaders(req)
	if err != nil {
		fmt.Println(err)
	}
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
	return client.get("account")
}
/*
func (client *Client) getActiveOrders(market, page) (string) {
	requestBody, err := json.Marshal(map[string]string{
		"market":market,
		"page":page,
	})
	if err != nil {
		fmt.Println(err)
	}

	return client.get()
}
*/

func main () {
	client, err := NewClient(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("error making the client")
	}
	fmt.Println(client.get("account"))
	fmt.Println(client.getAccount())

}