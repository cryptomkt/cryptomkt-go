package requests

import (
	"fmt"
)

// A Request holds the arguments for an http request to the server.
type Request struct {
	arguments map[string]string
	required  []string
}

// AddArgument adds an argument to a request, guiven a pair key value.
func (req *Request) AddArgument(key, value string) {
	req.arguments[key] = value
}

// AssertRequired asserts that the required arguments of the request
// are meeted, if they arent, an error describing the missings required
// arguments is returned.
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

// GetArguments returns a map with both keys and values as strings.
// This function is intended to be used after the request has been
// proven valid by assertRequired.
func (req *Request) GetArguments() map[string]string {
	return req.arguments
}

// NewReq creates a new request with the required arguments to be meet
// as a string array slice.
func NewReq(required []string) *Request {
	return &Request{
		arguments: make(map[string]string),
		required:  required,
	}
}

// NewEmptyReq creates a new request that needs no arguments,
// and hence no validation to be used.
func NewEmptyReq() *Request {
	return &Request{
		arguments: make(map[string]string),
		required:  []string{},
	}
}
