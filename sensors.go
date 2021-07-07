package main

type SensorAddress string
type SensorHostname string
type SensorModel string

type SensorPosition struct {
	X      float32
	Y      float32
	Radian float32
}

type Sensor struct {
	Hostname SensorHostname
	Model    SensorModel
	address  SensorAddress
	Position *SensorPosition
}

var sensors = make(map[SensorAddress]Sensor)
