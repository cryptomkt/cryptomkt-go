package websocket

type wsNotification struct {
	ID     int64                  `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type wsSubscription struct {
	ID      int64                  `json:"id"`
	Method  string                 `json:"method"`
	Channel string                 `json:"ch"`
	Params  map[string]interface{} `json:"params"`
}

type wsResponse struct {
	ID     int64
	Method string
	Ch     string
}
