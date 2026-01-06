package color

import (
	"image/color"
)

// Zwift Hex Colors mapped to RGBA
var (
	// Zone 1: 7F7F7F (Grey)
	zwiftGrey = color.RGBA{127, 127, 127, 255}
	// Zone 2: 388AF5 (Blue)
	zwiftBlue = color.RGBA{56, 138, 245, 255}
	// Zone 3: 59BD59 (Green)
	zwiftGreen = color.RGBA{89, 189, 89, 255}
	// Zone 4: F8CC44 (Yellow)
	zwiftYellow = color.RGBA{248, 204, 68, 255}
	// Zone 5: ED6334 (Orange)
	zwiftOrange = color.RGBA{237, 99, 52, 255}
	// Zone 6: EC3123 (Red)
	zwiftRed = color.RGBA{236, 49, 35, 255}
)

func PowerToColor(power, ftp float64) color.RGBA {
	if ftp <= 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	ratio := power / ftp

	switch {
	case ratio < 0.60: // Zone 1: < 60% (Grey -> Blue)
		return zwiftGrey
	case ratio <= 0.75: // Zone 2: 60-75% (Blue -> Green)
		return zwiftBlue
	case ratio <= 0.89: // Zone 3: 76-89% (Green -> Yellow)
		return zwiftGreen
	case ratio <= 1.04: // Zone 4: 90-104% (Yellow -> Orange)
		return zwiftYellow
	case ratio <= 1.18: // Zone 5: 105-118% (Orange -> Red)
		return zwiftOrange
	default: // Zone 6: > 118% (Red)
		return zwiftRed
	}
}
