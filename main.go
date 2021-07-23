package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"google.golang.org/grpc"
	"log"
	"os"
	"xx_backend/pb"
)

var mqttClient mqtt.Client
var trackingClient pb.RawDataClient
var originPosition [2]float32
var threeScale float64 = 1000
var recording = false
var mqttOutput = true
var logOutput = false
var logfile *os.File

func main() {

	//------------
	// gRPC
	//------------

	conn, err := grpc.Dial("localhost:3010", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	trackingClient = pb.NewRawDataClient(conn)

	//------------
	// Fiber
	//------------

	go runWebsocketHub()

	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		runWebSocket(c)
	}))

	mqttClient = mqttConnect()
	go mqttListen(mqttClient, "connection", connectionHandler)
	go runTracking()
	//go mqttListen(mqttClient, "test", broadcastHandler)

	//go mqttListen(mqttClient, "pingcheck", pingHandler)
	//token := mqttClient.Publish("pingtest", 1, false, ".")
	//log.Println("start", time.Now().UnixNano())
	//token.Wait()

	log.Fatal(app.Listen(":3000"))
}
