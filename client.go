package client

import (
	"strings"

	//"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
)

type Client struct {
	apiVersion string
	baseApiUri string
	auth       *HMACAuth
	httpClient *http.Client
}

func NewClient(apiKey, apiSecret string) (*Client, error) {
	apiVersion := "v1"
	baseApiUri := "https://api.cryptomkt.com/"

	auth, err := NewAuth(apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	client := &Client{
		baseApiUri: baseApiUri,
		apiVersion: apiVersion,
		auth:       auth,
		httpClient: &http.Client{},
	}
	return client, nil
}

func (client *Client) get(endpoint string, argsmap map[string]interface{}) string {
	args, err := Mapss(argsmap)

	if err != nil {
		fmt.Println(err)
	}
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		fmt.Println("could not parse the base api uri", client.baseApiUri)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println(err)
	}

	if len(args) != 0 {
		q := req.URL.Query()
		for k, v := range args {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	requestPath := fmt.Sprintf("/%s/%s", client.apiVersion, endpoint)
	client.auth.SetHeaders(req, requestPath, "")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(respBody)
}

func (client *Client) post(endpoint string, argsmap map[string]interface{}) string {
	args, err := Mapss(argsmap)
	if err != nil {
		fmt.Println(err)
	}
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		fmt.Println("could not parse the base api uri", client.baseApiUri)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)
	if err != nil {
		fmt.Println(err)
	}

	form := url.Values{}
	for k, v := range args {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	if len(args) == 0 {
		fmt.Println("error, post with no information")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	var body string

	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		body = fmt.Sprintf("%s%v", body, args[k])
	}
	requestPath := fmt.Sprintf("/%s/%s", client.apiVersion, endpoint)
	client.auth.SetHeaders(req, requestPath, body)
	if err != nil {
		fmt.Println("cannot parse form")
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(respBody)
}
