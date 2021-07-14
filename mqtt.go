package main

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"xx_backend/pb"
)

func mqttConnect() mqtt.Client {
	var broker = "192.168.178.48"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("communication")
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

func mqttListen(client mqtt.Client, topic string, handler func(mqtt.Client, mqtt.Message)) {
	client.Subscribe(topic, 0, handler)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func broadcastHandler(_ mqtt.Client, msg mqtt.Message) {

	message := string(msg.Payload())
	splitMessage := strings.Split(message, ";")

	if splitMessage[0] == "raw" {
		rawMessage := splitMessage[1:]

		var measures []string

		measuresSplit := strings.Split(rawMessage[0], "!")

		measures = append(measures, measuresSplit[:len(measuresSplit)-1]...)

		var coordinates []*pb.Coordinate

		wsMessage := msg.Topic() + ";pos;"
		for _, m := range measures {

			values := strings.Split(m, ":")
			angle, _ := strconv.ParseFloat(values[0], 64)
			distance, _ := strconv.ParseFloat(values[1], 64)

			x := float32(math.Sin(angle*math.Pi/180) * distance)
			y := float32(math.Cos(angle*math.Pi/180) * distance)

			coordinates = append(coordinates, &pb.Coordinate{X: x, Y: y})
		}

		resp, err := trackingClient.Transform(context.Background(), &pb.TransformRequest{
			RawCoordinates: coordinates,
			Radian:         sensors[SensorAddress(msg.Topic())].Position.Radian,
		})
		if err != nil {
			log.Fatal(err)
		}

		rotatedCoordinates := resp.GetTransformedCoordinates()

		scan := Scan{}

		for _, coordinate := range rotatedCoordinates {

			sensorPostionUpdatedX := coordinate.X + sensors[SensorAddress(msg.Topic())].Position.X
			sensorPostionUpdatedY := coordinate.Y + sensors[SensorAddress(msg.Topic())].Position.Y

			c := Coordinate{
				X: sensorPostionUpdatedX,
				Y: sensorPostionUpdatedY,
			}
			scan = append(scan, c)

			wsMessage += strconv.FormatFloat(float64(sensorPostionUpdatedX), 'f', 5, 32) + ":" + strconv.FormatFloat(float64(sensorPostionUpdatedY), 'f', 5, 64) + "!"
		}

		sensorLastScan[SensorAddress(msg.Topic())] = scan

		wsBroadcast <- wsMessage
	} else {
		wsBroadcast <- msg.Topic() + ";" + string(msg.Payload())
	}
}

func connectionHandler(_ mqtt.Client, msg mqtt.Message) {

	messageString := string(msg.Payload())

	messageParts := strings.Split(messageString, ";")

	if messageParts[1] == "+" {
		if len(messageParts) == 4 {
			sensorAddress := SensorAddress(messageParts[0])
			log.Println("connected", messageString)

			_, ok := sensors[sensorAddress]
			if ok {
				sensor := sensors[sensorAddress]
				sensor.Connected = true
				sensors[sensorAddress] = sensor
			} else {
				sensor := Sensor{
					Hostname: SensorHostname(messageParts[2]),
					Model:    SensorModel(messageParts[3]),
					Color:    SensorColor(sensorColors[sensorColorIndex%len(sensorColors)]),
					Position: &SensorPosition{
						X:      0,
						Y:      0,
						Radian: 0,
					},
					Connected: true,
				}
				sensors[sensorAddress] = sensor
				sensorColorIndex++
			}

			log.Println(string(sensorAddress))
			go mqttListen(mqttClient, string(sensorAddress), broadcastHandler)
		}
	} else if messageParts[1] == "-" {
		log.Println("disconnected", messageString)
		sensorAddress := SensorAddress(messageParts[0])
		_, ok := sensors[sensorAddress]
		if ok {
			sensor := sensors[sensorAddress]
			sensor.Connected = false
			sensors[sensorAddress] = sensor
			//delete(sensors, sensorAddress)
		}
		sensorLastScan[sensorAddress] = Scan{}
	}

	log.Println(sensors)

	wsBroadcast <- generateSensorsMessage()
}
