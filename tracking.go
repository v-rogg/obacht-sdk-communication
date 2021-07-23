package main

import (
	"context"
	"log"
	"strconv"
	"time"
	"xx_backend/pb"
)

var dbscanEps float32 = 250
var dbscanMinSamples int64 = 15

func runTracking() {
	//for {
	var zones []*pb.Zone
	for sensorAddress, sensorZone := range sensorZones {
		var coordinates []*pb.Coordinate
		for _, position := range sensorZone {
			coordinates = append(coordinates, &pb.Coordinate{
				X: position.X,
				Y: position.Y,
			})
		}
		zone := &pb.Zone{
			Address:               string(sensorAddress),
			ZoneObjectCoordinates: coordinates,
		}
		zones = append(zones, zone)
	}

	var scans []*pb.Scan
	for sensorAddress, sensorScan := range sensorLastScan {

		_, ok := sensors[sensorAddress]
		if ok {
			sensor := sensors[sensorAddress]

			if sensor.Connected && len(sensorScan) > 0 {
				var coordinates []*pb.Coordinate
				for _, position := range sensorScan {
					coordinates = append(coordinates, &pb.Coordinate{
						X: position.X,
						Y: position.Y,
					})
				}
				scan := &pb.Scan{
					Address:    string(sensorAddress),
					ScanPoints: coordinates,
				}
				scans = append(scans, scan)
			}
		}
	}

	resp, err := trackingClient.TrackPersons(context.Background(), &pb.TrackRequest{
		Zones:      zones,
		Scans:      scans,
		Eps:        dbscanEps,
		MinSamples: dbscanMinSamples,
	})
	if err != nil {
		log.Fatal(err)
	}

	persons := resp.GetPersons()

	message := "system;persons;"
	for _, person := range persons {
		//log.Println(person.X, person.Y)
		message += strconv.FormatFloat(float64(person.GetX())/threeScale, 'f', 5, 64) + ":" + strconv.FormatFloat(float64(person.GetY())/threeScale, 'f', 5, 64) + "!"
	}
	wsBroadcast <- message

	time.AfterFunc(60*time.Millisecond, runTracking)
}
