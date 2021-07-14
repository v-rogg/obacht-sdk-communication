package main

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"log"
	"strconv"
	"strings"
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
			message := generateSensorsMessage()
			if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("write error initial sensor message:", err)
			}
			message = generateOriginMessage()
			if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("write error initial origin message:", err)
			}
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

func runWebSocket(c *websocket.Conn) {
	defer func() {
		wsUnregister <- c
		_ = c.Close()
	}()

	wsRegister <- c

	for {
		messageType, messagePayload, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			return // Calls the deferred function, i.e. closes the connection on error
		}
		if messageType == websocket.TextMessage {

			message := string(messagePayload)
			fmt.Println(message)

			splitMessage := strings.Split(message, ";")

			if splitMessage[0] == "system" {

				messageProtocolType := splitMessage[1]

				switch messageProtocolType {
				case "rotate":
					objectType := splitMessage[2]

					switch objectType {
					case "sensors":
						f, _ := strconv.ParseFloat(splitMessage[4], 32)
						sensors[SensorAddress(splitMessage[3])].Position.Radian = float32(f)
						break
					case "zones":
						break
					default:
						break
					}
					break
				case "move":
					objectType := splitMessage[2]

					switch objectType {
					case "origin":
						x, _ := strconv.ParseFloat(splitMessage[4], 32)
						y, _ := strconv.ParseFloat(splitMessage[5], 32)
						originPosition[0] = float32(x * threeScale)
						originPosition[1] = float32(y * threeScale)
						break
					case "sensors":
						_, ok := sensors[SensorAddress(splitMessage[3])]
						if ok {
							x, _ := strconv.ParseFloat(splitMessage[4], 32)
							y, _ := strconv.ParseFloat(splitMessage[5], 32)
							sensors[SensorAddress(splitMessage[3])].Position.X = float32(x * threeScale)
							sensors[SensorAddress(splitMessage[3])].Position.Y = float32(y * threeScale)
						}
						break
					case "zones":
						break
					default:
						break
					}
					break
				case "zones":
					analyseZonesMessage(splitMessage[2])
					log.Println(string(messagePayload))
					break
				default:
					break
				}
			}

			wsBroadcast <- string(messagePayload)
			// Do something with the message
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
