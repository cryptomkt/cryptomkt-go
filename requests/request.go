package requests

import (
	"fmt"
)

type Request struct {
	arguments map[string]string
	required  []string
}

func (req *Request) AddArgument(key, value string) {
	req.arguments[key] = value
}

func (req *Request) AssertRequired() error {
	errMsg := ""
	needOptions := false
	for _, key := range req.required {
		if _, ok := req.arguments[key]; !ok {
			errMsg = fmt.Sprintf("%s %s,", errMsg, key)
			needOptions = true
		}
	}
	if needOptions {
		return fmt.Errorf("%s", errMsg[:len(errMsg)-1])
	}
	return nil
}

func (req *Request) GetArguments() map[string]string {
	return req.arguments
}

func NewReq(required []string) *Request {
	return &Request{
		arguments: make(map[string]string),
		required:  required,
	}
}

func NewEmptyReq() (*Request) {
	return &Request{
		arguments: make(map[string]string),
		required:  []string{},
	}
}
