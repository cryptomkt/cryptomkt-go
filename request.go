package client

import (
	"fmt"
)

type Request struct {
	arguments map[string]string
	required []string
}

func (req *Request) addArgument(key, value string){
	req.arguments[key] = value
}

func (req *Request) assertRequired() error{
	errMsg := "Required options not meeted:"
	needOptions := false
	for _, key := range req.required {
		if _, ok := req.arguments[key]; !ok {
			errMsg = fmt.Sprintf("%s %s,", errMsg, key)
			needOptions = true
		}
	}
	if needOptions {
		return fmt.Errorf("client: %s", errMsg[:len(errMsg)-1])
	}
	return nil
}

