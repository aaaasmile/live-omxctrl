package live

import (
	"log"
	"net/http"

	"github.com/aaaasmile/live-omxctrl/web/live/ws"
	"github.com/gorilla/websocket"
)

var (
	upgrader  websocket.Upgrader
	wsHanlder *ws.WsHandler
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS error", err)
		return
	}

	wsHanlder.AddConn(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error ", err)
			wsHanlder.CloseConn(conn)
			return
		}
		if messageType == websocket.TextMessage {
			log.Println("Message rec: ", string(p))
		}
	}
}

func init() {
	wsHanlder = ws.NewWsHandler()
	wsHanlder.StartWS()
}
