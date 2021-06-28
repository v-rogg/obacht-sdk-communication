package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func main() {

	go runWebsocketHub()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
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
			messageType, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				//log.Println(string(message))
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	mqttClient := mqttConnect()
	go mqttListen(mqttClient, "test", subscriptionHandler)
	go mqttListen(mqttClient, "$connected", connectionHandler)
	go mqttListen(mqttClient, "pingcheck", pingHandler)

	//token := mqttClient.Publish("pingtest", 1, false, ".")
	//log.Println("start", time.Now().UnixNano())
	//token.Wait()

	log.Fatal(app.Listen(":3000"))
}
