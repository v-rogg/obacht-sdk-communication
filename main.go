package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"log"
)

var mqttClient mqtt.Client
var cubePosition [2]int64

func init() {
	cubePosition[0] = 0
	cubePosition[1] = 0
}

func main() {

	go runWebsocketHub()

	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
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
				fmt.Println(string(messagePayload))
				wsBroadcast <- string(messagePayload)
				// Do something with the message
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	app.Use(cors.New())
	app.Get("/init", func(c *fiber.Ctx) error {
		return c.JSON(sensors)
	})

	mqttClient = mqttConnect()
	go mqttListen(mqttClient, "$connected", connectionHandler)
	go mqttListen(mqttClient, "test", broadcastHandler)

	//go mqttListen(mqttClient, "pingcheck", pingHandler)
	//token := mqttClient.Publish("pingtest", 1, false, ".")
	//log.Println("start", time.Now().UnixNano())
	//token.Wait()

	log.Fatal(app.Listen(":3000"))
}
