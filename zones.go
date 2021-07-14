package main

import (
	"log"
	"strconv"
	"strings"
)

type ZonePosition struct {
	X float32
	Y float32
}

type Zone []ZonePosition

var sensorZones = make(map[SensorAddress]Zone)

func generateZonesMessage() string {
	message := "system;zones;"

	for sensorAddress, sensorZone := range sensorZones {
		message += string(sensorAddress) + ","

		for _, position := range sensorZone {

			x := strconv.FormatFloat(float64(position.X)/threeScale, 'f', 5, 64)
			y := strconv.FormatFloat(float64(position.Y)/threeScale, 'f', 5, 64)

			message += x + ":" + y + "?"
		}

		message += "!"
	}

	log.Println(message)

	return message
}

func analyseZonesMessage(message string) {
	sensorMessageParts := strings.Split(message, "!")

	for i := 0; i < len(sensorMessageParts)-1; i++ {
		sensorMessage := strings.Split(sensorMessageParts[i], ",")

		sensorAddress := SensorAddress(sensorMessage[0])
		var zone Zone
		zone = []ZonePosition{}

		sensorPositions := strings.Split(sensorMessage[1], "?")

		for j := 0; j < len(sensorPositions)-1; j++ {
			xy := strings.Split(sensorPositions[j], ":")

			x, _ := strconv.ParseFloat(xy[0], 64)
			y, _ := strconv.ParseFloat(xy[1], 64)

			position := ZonePosition{
				X: float32(x * threeScale),
				Y: float32(y * threeScale),
			}
			log.Println(position)

			zone = append(zone, position)
		}
		sensorZones[sensorAddress] = zone
	}

	log.Println(sensorLastScan)
}
