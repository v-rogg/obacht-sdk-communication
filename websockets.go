package main

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

type wsClient struct{}

var wsClients = make(map[*websocket.Conn]wsClient)
var wsRegister = make(chan *websocket.Conn)
var wsBroadcast = make(chan string)
var wsUnregister = make(chan *websocket.Conn)

func runWebsocketHub() {
	for {
		select {
		case connection := <-wsRegister:
			wsClients[connection] = wsClient{}
			log.Println("connection registered")
		case message := <-wsBroadcast:
			for connection := range wsClients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					err = connection.WriteMessage(websocket.CloseMessage, []byte{})
					if err != nil {
						log.Println("Couldn't write CloseMessage", err)
					}
					err = connection.Close()
					if err != nil {
						log.Println("Couldn't close connection", err)
					}
					delete(wsClients, connection)
				}
			}
		case connection := <-wsUnregister:
			delete(wsClients, connection)
			log.Println("connection unregistered")
		}
	}
}
