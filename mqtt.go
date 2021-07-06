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

			x := math.Sin(angle*math.Pi/180) * distance
			y := math.Cos(angle*math.Pi/180) * distance

			coordinates = append(coordinates, &pb.Coordinate{X: float32(x), Y: float32(y)})
		}

		resp1, err := trackingClient.Transform(context.Background(), &pb.TransformRequest{
			RawCoordinates: coordinates,
			Radian:         float32(radian),
		})
		if err != nil {
			log.Fatal(err)
		}

		transformedCoordinates := resp1.GetTransformedCoordinates()

		for _, coordinate := range transformedCoordinates {
			wsMessage += strconv.FormatFloat(float64(coordinate.X), 'f', 5, 32) + ":" + strconv.FormatFloat(float64(coordinate.Y), 'f', 5, 64) + "!"
		}

		wsBroadcast <- wsMessage
	} else {
		wsBroadcast <- msg.Topic() + ";" + string(msg.Payload())
	}

	//req := pb.TransformRequest{
	//	RawCoordinate: nil,
	//	Radian:        0,
	//}

}

func pingHandler(_ mqtt.Client, msg mqtt.Message) {
	wsBroadcast <- string(msg.Payload())

	log.Println("end  ", time.Now().UnixNano())
	log.Println(string(msg.Payload()))
}

func connectionHandler(_ mqtt.Client, msg mqtt.Message) {

	messageString := string(msg.Payload())

	messageParts := strings.Split(messageString, ";")

	var sensor Sensor

	if messageParts[1] == "+" {
		if len(messageParts) == 4 {
			sensorAddress := SensorAddress(messageParts[0])
			sensor = Sensor{
				Hostname: SensorHostname(messageParts[2]),
				Model:    SensorModel(messageParts[3]),
			}

			sensors[sensorAddress] = sensor

			log.Println(string(sensorAddress))
			go mqttListen(mqttClient, string(sensorAddress), broadcastHandler)
		}
	} else if messageParts[1] == "-" {
		sensorAddress := SensorAddress(messageParts[2])
		_, ok := sensors[sensorAddress]
		if ok {
			delete(sensors, sensorAddress)
		}
	}

	log.Println(sensors)

	wsBroadcast <- messageString
	log.Println(messageString)
}
