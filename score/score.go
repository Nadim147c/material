package score

import (
	"math"
	"slices"
	"sort"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
)

// Google Blue
var FallbackColor color.ARGB = 0xff4285f4

// ScoreOptions provides configuration for ranking colors based on usage counts.
//
// Desired: is the max count of the colors returned.
// FallbackColorARGB: Is the default color that should be used if no other
// colors are suitable.
//
// Filter: controls if the resulting colors should be filtered to not include
// hues that are not used often enough, and colors that are effectively
// grayscale.
type ScoreOptions struct {
	Desired  int
	Fallback color.ARGB
	Filter   bool
}

// scoredColor holds a color and its calculated score
type scoredColor struct {
	hct   color.Hct
	score float64
}

// score provides functions for ranking colors based on suitability for UI
// themes.
type score struct {
	// These would be constants in Go
	targetChroma            float64
	weightProportion        float64
	weightChromaAbove       float64
	weightChromaBelow       float64
	cutoffChroma            float64
	cutoffExcitedProportion float64
}

// NewScore creates a new Score instance with default constants
func NewScore() *score {
	return &score{
		targetChroma:            48.0, // A1 Chroma
		weightProportion:        0.7,
		weightChromaAbove:       0.3,
		weightChromaBelow:       0.1,
		cutoffChroma:            5.0,
		cutoffExcitedProportion: 0.01,
	}
}

// SanitizeDegreesInt ensures a degree measure is within the range [0, 360).
func SanitizeDegreesInt(degrees int) int {
	degrees = degrees % 360
	if degrees < 0 {
		degrees += 360
	}
	return degrees
}

// DifferenceDegrees returns the shortest angular difference between two angles
// in degrees.
func DifferenceDegrees(a, b float64) float64 {
	return num.NormalizeDegree(a - b)
}

// ScoreColors ranks colors based on suitability for being used for a UI theme.
//
// colorsToPopulation: map with keys of colors and values of how often
// the color appears, usually from a source image.
// options: optional parameters for customizing scoring behavior.
//
// Returns: Colors sorted by suitability for a UI theme. The most suitable
// color is the first item, the least suitable is the last. There will
// always be at least one color returned. If all the input colors
// were not suitable for a theme, a default fallback color will be
// provided, Google Blue.
func (s *score) ScoreColors(
	colorsToPopulation map[color.ARGB]int,
	opts ScoreOptions,
) []color.ARGB {
	if opts.Desired == 0 {
		opts.Desired = 4
	}
	if opts.Fallback == 0 {
		opts.Fallback = FallbackColor
	}

	// Get the HCT color for each Argb value, while finding the per hue count
	// and
	// total count.
	colorsHct := []color.Hct{}
	huePopulation := make([]int, 360)
	populationSum := 0

	for argb, population := range colorsToPopulation {
		hct := argb.ToHct()
		colorsHct = append(colorsHct, hct)
		hue := int(math.Floor(hct.Hue))
		huePopulation[hue] += population
		populationSum += population
	}

	// Hues with more usage in neighboring 30 degree slice get a larger number.
	hueExcitedProportions := make([]float64, 360)
	for hue := range 360 {
		proportion := float64(huePopulation[hue]) / float64(populationSum)
		for i := hue - 14; i < hue+16; i++ {
			neighborHue := SanitizeDegreesInt(i)
			hueExcitedProportions[neighborHue] += proportion
		}
	}

	// Scores each HCT color based on usage and chroma, while optionally
	// filtering out values that do not have enough chroma or usage.
	scoredHct := []scoredColor{}
	for hct := range slices.Values(colorsHct) {
		hue := SanitizeDegreesInt(int(math.Round(hct.Hue)))
		proportion := hueExcitedProportions[hue]

		if opts.Filter &&
			(hct.Chroma < s.cutoffChroma || proportion <= s.cutoffExcitedProportion) {
			continue
		}

		proportionScore := proportion * 100.0 * s.weightProportion

		var chromaWeight float64
		if hct.Chroma < s.targetChroma {
			chromaWeight = s.weightChromaBelow
		} else {
			chromaWeight = s.weightChromaAbove
		}

		chromaScore := (hct.Chroma - s.targetChroma) * chromaWeight
		score := proportionScore + chromaScore

		scoredHct = append(scoredHct, scoredColor{hct: hct, score: score})
	}

	// Sort so that colors with higher scores come first
	sort.Slice(scoredHct, func(i, j int) bool {
		return scoredHct[i].score > scoredHct[j].score
	})

	// Iterates through potential hue differences in degrees in order to select
	// the colors with the largest distribution of hues possible. Starting at
	// 90 degrees(maximum difference for 4 colors) then decreasing down to a
	// 15 degree minimum.
	chosenColors := []color.Hct{}
	for differenceDegrees := 90; differenceDegrees >= 15; differenceDegrees-- {
		chosenColors = []color.Hct{} // Clear the array

		for scored := range slices.Values(scoredHct) {
			duplicateHue := false

			for chosenHct := range slices.Values(chosenColors) {
				if DifferenceDegrees(
					scored.hct.Hue,
					chosenHct.Hue,
				) < float64(
					differenceDegrees,
				) {
					duplicateHue = true
					break
				}
			}

			if !duplicateHue {
				chosenColors = append(chosenColors, scored.hct)
			}

			if len(chosenColors) >= opts.Desired {
				break
			}
		}

		if len(chosenColors) >= opts.Desired {
			break
		}
	}

	// Convert chosen HCT colors back to ARGB integers
	colors := []color.ARGB{}
	if len(chosenColors) == 0 {
		colors = append(colors, opts.Fallback)
	}

	for chosenHct := range slices.Values(chosenColors) {
		colors = append(colors, chosenHct.ToARGB())
	}

	return colors
}

// Score is a package-level convenience function that creates a Score instance
// and returns scored colors
func Score(
	colorsToPopulation map[color.ARGB]int,
	opts ScoreOptions,
) []color.ARGB {
	scorer := NewScore()
	return scorer.ScoreColors(colorsToPopulation, opts)
}
