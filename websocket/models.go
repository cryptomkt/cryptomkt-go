package websocket

import (
	"fmt"
)

type wsNotification struct {
	ID     int64                  `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type wsResponse struct {
	ID     int64
	Method string
	Params struct {
		Symbol string
		Period string
	}
}

type withError struct {
	Error APIError
}

// APIError is an error from the exchange
type APIError map[string]interface{}

func (errMap APIError) String() string {
	errString := ""
	if code, ok := errMap["code"]; ok {
		errString += "(code=" + fmt.Sprint(code) + ")"
	}
	if message, ok := errMap["message"]; ok {
		errString += " " + message.(string)
	}
	if desc, ok := errMap["description"]; ok {
		errString += ". " + desc.(string)
	}
	return errString
}
