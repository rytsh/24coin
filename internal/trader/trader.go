package trader

import (
	"crypto/tls"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/rytsh/24coin/internal/api"
)

var socketConn *websocket.Conn
var done = make(chan struct{})

// WebSocketConnect to get stream data
func WebSocketConnect() {
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	u := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws"}

	var err error
	socketConn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("[Error] Dial: %v", err)
	}
	log.Println("Connected to binance websocket", u.String())

	go func() {
		defer func() {
			close(done)
		}()
		for {
			if _, message, err := socketConn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					log.Printf("[Warning] ReadMessage: %v", err)
				}
				return
			} else {
				api.BroadcastMessage(message)
			}
		}
	}()

	log.Println("Adding btcusdt symbol")
	err = socketConn.WriteMessage(websocket.TextMessage, []byte(
		`{
			"method": "SUBSCRIBE",
			"params": [
			  "btcusdt@bookTicker"
			],
			"id": 1
		}`,
	))
	if err != nil {
		log.Printf("[Error] WriteMessage %v", err)
		return
	}
}

func Close() {
	if socketConn != nil {
		socketConn.Close()
	}
	<-done
}
