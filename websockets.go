package main

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"log"
	"os"
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
			if err := connection.WriteMessage(websocket.TextMessage, []byte("system;dbscan;eps;"+strconv.FormatFloat(float64(dbscanEps), 'f', 2, 32))); err != nil {
				log.Println("write error initial dbscan eps message:", err)
			}
			if err := connection.WriteMessage(websocket.TextMessage, []byte("system;dbscan;minSamples;"+strconv.FormatInt(dbscanMinSamples, 10))); err != nil {
				log.Println("write error initial dbscan minSamples message:", err)
			}
			if err := connection.WriteMessage(websocket.TextMessage, []byte("system;output;mqtt;"+strconv.FormatBool(mqttOutput))); err != nil {
				log.Println("write error initial mqttOutput message:", err)
			}
			if err := connection.WriteMessage(websocket.TextMessage, []byte("system;output;log;"+strconv.FormatBool(logOutput))); err != nil {
				log.Println("write error initial logOutput message:", err)
			}
			if err := connection.WriteMessage(websocket.TextMessage, []byte("system;output;recording;"+strconv.FormatBool(recording))); err != nil {
				log.Println("write error initial recording message:", err)
			}
			if err := connection.WriteMessage(websocket.TextMessage, []byte(generateZonesMessage())); err != nil {
				log.Println("write error initial zones message:", err)
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

			if recording {
				if logOutput {
					if _, err := logfile.WriteString(message + "\n"); err != nil {
						log.Println(err)
					}
				}
				if mqttOutput {
					mqttClient.Publish("output", 0, false, message)
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
						_, ok := sensorZones[SensorAddress(splitMessage[3])]
						if ok {
							x, _ := strconv.ParseFloat(splitMessage[5], 32)
							y, _ := strconv.ParseFloat(splitMessage[6], 32)
							index, _ := strconv.ParseInt(splitMessage[4], 10, 32)

							sensorZones[SensorAddress(splitMessage[3])][index].X = float32(x * threeScale)
							sensorZones[SensorAddress(splitMessage[3])][index].Y = float32(y * threeScale)
						}
						break
					default:
						break
					}
					break
				case "add":
					objectType := splitMessage[2]

					if objectType == "zones" {
						address := SensorAddress(splitMessage[3])
						//index := splitMessage[4]
						x, _ := strconv.ParseFloat(splitMessage[5], 32)
						y, _ := strconv.ParseFloat(splitMessage[6], 32)

						position := ZonePosition{
							X: float32(x * threeScale),
							Y: float32(y * threeScale),
						}

						_, ok := sensorZones[address]
						if ok {
							sensorZones[address] = append(sensorZones[address], position)
						} else {
							var zone Zone
							zone = []ZonePosition{}
							zone = append(zone, position)
							sensorZones[address] = zone
						}

						log.Println(sensorZones)
					}
				case "remove":
					objectType := splitMessage[2]

					if objectType == "zones" {
						address := SensorAddress(splitMessage[3])
						delete(sensorZones, address)
					}
				case "zones":
					analyseZonesMessage(splitMessage[2])
					//log.Println(string(messagePayload))
					break
				case "dbscan":
					switch splitMessage[2] {
					case "eps":
						messageEps, _ := strconv.ParseFloat(splitMessage[3], 32)
						dbscanEps = float32(messageEps)
						//log.Println(dbscanEps)
						break
					case "minSamples":
						messageMinSamples, _ := strconv.ParseInt(splitMessage[3], 10, 32)
						dbscanMinSamples = messageMinSamples
						//log.Println(dbscanMinSamples)
						break
					}
					break
				case "output":
					switch splitMessage[2] {
					case "mqtt":
						messageMqtt, _ := strconv.ParseBool(splitMessage[3])
						mqttOutput = messageMqtt
						break
					case "log":
						messageLog, _ := strconv.ParseBool(splitMessage[3])
						logOutput = messageLog
						break
					case "recording":
						messageRecording, _ := strconv.ParseBool(splitMessage[3])
						recording = messageRecording

						if recording {
							logfile, err = os.OpenFile("D:/output.log",
								os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
							if err != nil {
								log.Println(err)
							}
						} else {
							err = logfile.Close()
							if err != nil {
								log.Println(err)
							}
						}
						break
					}
					break
				default:
					break
				}
			}

			wsBroadcast <- string(messagePayload)
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
