package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

type client struct{}

var clients = make(map[*websocket.Conn]client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{}
			log.Println("connection registered")
		case message := <-broadcast:
			for connection := range clients {
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
					delete(clients, connection)
				}
			}
		case connection := <-unregister:
			delete(clients, connection)
			log.Println("connection unregistered")
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func subscriptionHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	broadcast <- string(msg.Payload())
}

func connect() mqtt.Client {
	var broker = "192.168.178.48"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("backend")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func listen(topic string) {
	client := connect()
	client.Subscribe(topic, 0, subscriptionHandler)
}

func main() {
	app := fiber.New()
	//
	//app.Get("/", func (c *fiber.Ctx) error {
	//	return c.SendString("Hello, World!")
	//})

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go runHub()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		register <- c

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

	go listen("test")

	//defer func(client mqtt.Client) {client.Disconnect(250)}(client)

	log.Fatal(app.Listen(":3000"))
}
