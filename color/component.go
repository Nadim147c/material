package color

import (
	"math"

	"github.com/Nadim147c/goyou/num"
)

// Linearized takes component (uint8) that represents R/G/B channel.
// Returns 0.0 <= output <= 100.0, color channel converted to linear RGB space
func Linearized(component uint8) float64 {
	normalized := float64(num.Clamp(0, 0xFF, component)) / 0xFF
	if normalized <= 0.040449936 {
		return normalized / 12.92 * 100.0
	} else {
		return math.Pow((normalized+0.055)/1.055, 2.4) * 100.0
	}
}

// Delinearized takes component (float64) that represents linear R/G/B channel.
// Component should be 0.0 < component < 100.0. Returns uint8 (0 <= n <= 255)
// representation of color component.
func Delinearized(component float64) uint8 {
	normalized := num.Clamp(0, 100, component) / 100
	delinearized := 0.0
	if normalized <= 0.0031308 {
		delinearized = normalized * 12.92
	} else {
		delinearized = 1.055*math.Pow(normalized, 1.0/2.4) - 0.055
	}
	return num.Clamp(0, 0xFF, uint8(math.Round(delinearized*255.0)))
}
