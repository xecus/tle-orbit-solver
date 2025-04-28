package util

import "math"

// EarthRadius is the Earth's radius in kilometers
const EarthRadius = 6356.752

// Deg2Rad converts degrees to radians
func Deg2Rad(deg float64) float64 {
	return deg / 180.0 * math.Pi
}

// Rad2Deg converts radians to degrees
func Rad2Deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}