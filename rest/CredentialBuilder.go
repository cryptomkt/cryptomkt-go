package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

func (client *httpclient) getCredentialForRequest(request *http.Request, body string) string {
	timestamp := now()
	message := getMessage(request, body, timestamp, client.window)
	signature := getMessageSignature(client, message)
	credential := client.apiKey + ":" + signature + ":" + timestamp
	if client.window != 0 {
		credential += (":" + strconv.FormatInt(int64(client.window), 10))
	}
	return "HS256 " + base64.StdEncoding.EncodeToString([]byte(credential))
}

func now() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func getMessageSignature(client *httpclient, msg string) string {
	hash := hmac.New(sha256.New, []byte(client.apiSecret))
	hash.Write([]byte(msg))
	return hex.EncodeToString(hash.Sum(nil))
}

func getMessage(request *http.Request, body, timestamp string, window int) string {
	message := request.Method + request.URL.Path
	message = addBodyToMessageIfPresent(message, request.Method, body)
	message += timestamp
	if window != 0 {
		message += strconv.FormatInt(int64(window), 10)
	}
	return message
}

func addBodyToMessageIfPresent(msg, method, body string) string {
	if len(body) == 0 {
		return msg
	}
	if method == methodGet {
		return msg + "?" + body
	}
	return msg + body
}
