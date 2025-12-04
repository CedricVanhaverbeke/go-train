package angle

import "math"

func ToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func ToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}
