package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"time"
)

func mqttConnect() mqtt.Client {
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
	wsBroadcast <- "~" + msg.Topic() + ";" + string(msg.Payload())
}

func pingHandler(_ mqtt.Client, msg mqtt.Message) {
	wsBroadcast <- string(msg.Payload())

	log.Println("end  ", time.Now().UnixNano())
	log.Println(string(msg.Payload()))
}

func connectionHandler(_ mqtt.Client, msg mqtt.Message) {

	messageString := string(msg.Payload())

	messageParts := strings.Split(messageString, ":")

	var sensor Sensor

	if messageParts[0][0:1] == "+" {
		if len(messageParts) == 4 {
			sensorAddress := SensorAddress(messageParts[2])
			sensor = Sensor{
				Hostname: SensorHostname(messageParts[1]),
				Model:    SensorModel(messageParts[3]),
			}

			sensors[sensorAddress] = sensor

			log.Println(string(sensorAddress))
			go mqttListen(mqttClient, string(sensorAddress), broadcastHandler)
		}
	} else if messageParts[0][0:1] == "-" {
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
