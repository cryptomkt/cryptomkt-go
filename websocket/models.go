package websocket

import "github.com/cryptomkt/go-api/models"

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
	Error map[string]interface{}
}

type requestsChans struct {
	ch    chan []byte
	errCh chan error
}

func makeChans() *requestsChans {
	return &requestsChans{
		ch:    make(chan []byte),
		errCh: make(chan error),
	}
}

func (chans *requestsChans) close() {
	close(chans.ch)
	close(chans.errCh)
}

type orderbookSnapshot struct {
	Symbol    string             `json:"symbol"`
	Sequence  int64              `json:"sequence"`
	Timestamp string             `json:"timestamp"`
	Ask       []models.BookLevel `json:"ask"`
	Bid       []models.BookLevel `json:"bid"`
}
