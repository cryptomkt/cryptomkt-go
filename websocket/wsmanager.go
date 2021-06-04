package websocket

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const addr string = "api.exchange.cryptomkt.com"
const publicPath string = "/api/2/ws/public"
const tradingPath string = "/api/2/ws/trading"

type wsManager struct {
	streamPath string
	conn       *websocket.Conn
	snd        chan []byte
	rcv        chan []byte
	isOpen     bool
}

func newPublicWSManager() *wsManager {
	return &wsManager{
		streamPath: publicPath,
		snd:        make(chan []byte, 1),
		rcv:        make(chan []byte, 1),
		isOpen:     false,
	}
}

func newTradingWSManager() *wsManager {
	return &wsManager{
		streamPath: tradingPath,
		snd:        make(chan []byte, 1),
		rcv:        make(chan []byte, 1),
		isOpen:     false,
	}
}

func (ws *wsManager) connect() error {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "wss", Host: addr, Path: ws.streamPath}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	ws.conn = c

	go ws.rcvLoop()
	go ws.sndLoop()

	ws.isOpen = true
	return nil
}

func (ws *wsManager) close() {
	ws.isOpen = false
	close(ws.snd)
}

func (ws *wsManager) sndLoop() {
	for msg := range ws.snd {
		err := ws.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("write:", err)
			return
		}
	}
	// send close msg to server
	err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("write close:", err)
		return
	}
}

func (ws *wsManager) rcvLoop() {
	defer close(ws.rcv)
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			ws.conn.Close()
			return
		}
		ws.rcv <- message
	}
}
