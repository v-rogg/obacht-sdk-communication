package main

import (
	"log"
	"strconv"
)

type SensorAddress string
type SensorHostname string
type SensorModel string
type SensorColor string
type SensorConnectionStatus bool

var sensorColors = []string{
	"18A0FB",
	"93C700",
	"FFD500",
	"FF0000",
	"CC00CC",
	"960064",
}
var sensorColorIndex = 0

type SensorPosition struct {
	X      float32
	Y      float32
	Radian float32
}

type Sensor struct {
	Hostname  SensorHostname
	Model     SensorModel
	address   SensorAddress
	Color     SensorColor
	Position  *SensorPosition
	Connected SensorConnectionStatus
}

type Coordinate struct {
	X float32
	Y float32
}

type Scan []Coordinate

var sensors = make(map[SensorAddress]Sensor)
var sensorLastScan = make(map[SensorAddress]Scan)

func generateSensorsMessage() string {
	message := "system;sensors;"

	for address, sensor := range sensors {
		if sensor.Connected {
			x := strconv.FormatFloat(float64(sensor.Position.X), 'f', 5, 64)
			y := strconv.FormatFloat(float64(sensor.Position.Y), 'f', 5, 64)
			radian := strconv.FormatFloat(float64(sensor.Position.Radian), 'f', 5, 64)
			message += string(address) + ":" + string(sensor.Hostname) + ":" + string(sensor.Model) + ":" + string(sensor.Color) + ":" + x + ":" + y + ":" + radian + "!"
		}
	}

	log.Println(message)

	return message
}
