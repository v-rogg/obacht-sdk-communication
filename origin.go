package main

import "strconv"

func generateOriginMessage() string {
	x := strconv.FormatFloat(float64(originPosition[0]), 'f', 5, 64)
	y := strconv.FormatFloat(float64(originPosition[1]), 'f', 5, 64)

	return "system;origin;" + x + ";" + y
}
