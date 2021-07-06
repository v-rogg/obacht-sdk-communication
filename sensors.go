package main

type SensorAddress string
type SensorHostname string
type SensorModel string

type Sensor struct {
	Hostname SensorHostname
	Model    SensorModel
	address  SensorAddress
}

var sensors = make(map[SensorAddress]Sensor)
