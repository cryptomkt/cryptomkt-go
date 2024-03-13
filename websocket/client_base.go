package websocket

import (
	"encoding/json"
)

type clientBase struct {
	wsManager *wsManager
	chanCache *chanCache
	window    int
}

// Close closes all the channels related to the client as well as the websocket connection.
// trying to make requests over a closed client may result in error.
func (client *clientBase) Close() {
	client.chanCache.close()
	client.wsManager.close()
}

func (client *clientBase) handle(rcvCh chan []byte) {
	for data := range rcvCh {
		resp := wsResponse{}
		json.Unmarshal(data, &resp)
		if reusableCh, ok := client.chanCache.getChan(resp.ID); ok {
			reusableCh.send(data)
			if reusableCh.isDone() {
				client.chanCache.closeAndRemoveCh(resp.ID)
			}
			continue
		}
		if key, ok := getSubscriptionKey(resp); ok {
			client.chanCache.sendViaSubscriptionCh(key, data)
			continue
		}
	}
}
