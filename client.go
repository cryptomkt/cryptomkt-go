package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
)

func main () {
	req, err := http.NewRequest("GET", "https://api.cryptomkt.com/v1/account", nil)
	if err!=nil {
		fmt.Println(err)
	}
	
	auth := New(os.Args[1], os.Args[2], "v1")
	err = auth.SetHeaders(req)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, err:= client.Do(req) 
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))

}