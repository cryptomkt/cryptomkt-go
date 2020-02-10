package main

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