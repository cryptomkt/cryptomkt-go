package client

import (
	"net/url"
	"fmt"
)

func checkUriSecurity(theUrl url.URL) (url.URL) {
	if theUrl.Scheme != "https" {
		fmt.Println(
			"WARNING: this client is sending a request to an insecure API ",
			"endpoint. Any API request you make may expose your API key and",
			"secret to third parties. Consider using the default endpoint:",
			"\n\n",
			theUrl,
		)
	}
	return theUrl
}

func parseMapss(data map[string]string) (string) {
	var args string = "?"
	for k, v := range data {
		args = fmt.Sprintf("%s%s=%v&",args, k, v)
	}
	return args
}

func parseMapsi(data map[string]interface{}) (string) {
	var args string = "?"
	for k, v := range data {
		args = fmt.Sprintf("%s%s=%v&",args, k, v)
		
	}
	return args[:len(args)-1]
}