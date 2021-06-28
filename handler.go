package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func subscriptionHandler(_ mqtt.Client, msg mqtt.Message) {
	//if string(msg.Payload())[0] == '@' {
	//messageDate := strings.Split(string(msg.Payload()), ";")
	//scanDate, err := time.Parse("2006-01-02T15:04:05.000-0700", messageDate[1])
	//if err != nil {
	//	log.Println(err)
	//}
	//delta := time.Since(scanDate)
	//fmt.Println(delta)
	//}
	wsBroadcast <- string(msg.Payload())
}

func pingHandler(_ mqtt.Client, msg mqtt.Message) {
	wsBroadcast <- string(msg.Payload())

	log.Println("end  ", time.Now().UnixNano())
	log.Println(string(msg.Payload()))
}

func connectionHandler(_ mqtt.Client, msg mqtt.Message) {
	wsBroadcast <- string(msg.Payload())
	log.Println(string(msg.Payload()))
}
