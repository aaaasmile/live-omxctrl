package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WsHandler struct {
	clients     map[*websocket.Conn]bool
	broadcastCh chan string
	mux         sync.Mutex
}

func NewWsHandler() *WsHandler {
	res := WsHandler{
		clients:     make(map[*websocket.Conn]bool),
		broadcastCh: make(chan string),
	}
	return &res
}

func (wh *WsHandler) AddConn(conn *websocket.Conn) {
	wh.mux.Lock()
	wh.clients[conn] = true
	wh.mux.Unlock()
	log.Println("New connection. Connected clients", len(wh.clients))
}

func (wh *WsHandler) CloseConn(conn *websocket.Conn) {
	wh.RemoveConn(conn)
	conn.Close()
}

func (wh *WsHandler) RemoveConn(conn *websocket.Conn) {
	wh.mux.Lock()
	delete(wh.clients, conn)
	wh.mux.Unlock()
	log.Println("Clients still connected ", len(wh.clients))
}

func (wh *WsHandler) closeAllConn() {
	for conn := range wh.clients {
		wh.CloseConn(conn)
	}
}

func (wh *WsHandler) broadcastMsg() {
	log.Println("WS Waiting for broadcast")
	for {
		msg := <-wh.broadcastCh

		for conn := range wh.clients {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				wh.CloseConn(conn)
				log.Println("Socket error: ", err)
			}
		}
	}
}

func (wh *WsHandler) StartWS() {
	go wh.broadcastMsg()
}

func (wh *WsHandler) EndWS() {
	log.Println("End od websocket service")
	wh.closeAllConn()
}
