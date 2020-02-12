package client

import (
	"strings"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
)

type Client struct {
	apiVersion string
	baseApiUri string
	auth       *HMACAuth
	httpClient *http.Client
}

func New(apiKey, apiSecret string) (*Client, error) {
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

func (client *Client) get(endpoint string, argsmap map[string]interface{}) (string, error) {
	args, err := Mapss(argsmap)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		return "", fmt.Errorf("client: Error parsing url %s: %v", client.baseApiUri, err)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("client: Error building NewRequest struct: %v", err)
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
		return "", fmt.Errorf("client: Error making request: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("client: Error reading response: %v", err)
	}
	return string(respBody), nil
}

func (client *Client) post(endpoint string, argsmap map[string]interface{}) (string, error) {
	if len(argsmap) == 0 {
		return "", fmt.Errorf("client: Must call with arguments")
	}
	args, err := Mapss(argsmap)
	if err != nil {
		return "", fmt.Errorf("client: Error parsing args as map[string]interface{} to map[string]string: %v", err)
	}
	u, err := url.Parse(client.baseApiUri)
	if err != nil {
		return "", fmt.Errorf("client: Error parsing url %s: %v", client.baseApiUri, err)
	}
	u.Path = path.Join(u.Path, client.apiVersion, endpoint)

	form := url.Values{}
	for k, v := range args {
		form.Add(k, v)
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("client: Error building NewRequest struct: %v", err)
	}
	
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var body string
	for _, k := range keys {
		body = fmt.Sprintf("%s%v", body, args[k])
	}
	requestPath := fmt.Sprintf("/%s/%s", client.apiVersion, endpoint)
	client.auth.SetHeaders(req, requestPath, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("client: Error making Request: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("client: Error reading response: %v", err)
	}
	return string(respBody), nil
}
