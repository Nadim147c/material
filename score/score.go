package score

import (
	"math"
	"slices"
	"sort"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/num"
)

// GoogleBlue is Google's Blue color as ARGB
var GoogleBlue color.ARGB = 0xff4285f4 // #4285f4

// options provides configuration for ranking colors based on usage counts.
// Options:
//   - Desired: is the max count of the colors returned.
//   - FallbackColorARGB: Is the default color that should be used if no other
//     colors are suitable.
//   - Filter: controls if the resulting colors should be filtered to not
//     include hues that are not used often enough, and colors that are
//     effectively grayscale.
type options struct {
	Limit    int
	Fallback color.ARGB
	Filter   bool
}

var defaultOptions = options{Limit: 4, Fallback: GoogleBlue, Filter: false}

// Option is a function that modifies the underlying options
type Option func(*options)

// WithLimit sets the limit for number of desired color
func WithLimit(limit int) Option {
	return func(o *options) {
		o.Limit = limit
	}
}

// WithFallback sets a fallback color if no color has been found
func WithFallback(c color.ARGB) Option {
	return func(o *options) {
		o.Fallback = c
	}
}

// WithFilter sets if the resulting colors should be filtered to not include
// hues that are not used often enough, and colors that are effectively
// grayscale.
func WithFilter() Option {
	return func(o *options) {
		o.Filter = true
	}
}

// scoredColor holds a color and its calculated score
type scoredColor struct {
	hct   color.Hct
	score float64
}

const (
	targetChroma            = 48.0 // A1 Chroma
	weightProportion        = 0.7
	weightChromaAbove       = 0.3
	weightChromaBelow       = 0.1
	cutoffChroma            = 5.0
	cutoffExcitedProportion = 0.01
)

// Score ranks colors based on suitability for being used for a UI theme.
//
// colorsToPopulation: map with keys of colors and values of how often
// the color appears, usually from a source image.
// options: optional parameters for customizing scoring behavior.
//
// Returns Colors sorted by suitability for a UI theme. The most suitable color
// is the first item, the least suitable is the last. There will always be at
// least one color returned. If all the input colors were not suitable for a
// theme, a default fallback color will be provided, Google Blue.
func Score(
	colorsToPopulation map[color.ARGB]int,
	options ...Option,
) []color.ARGB {
	opts := defaultOptions
	for option := range slices.Values(options) {
		option(&opts)
	}

	// Get the HCT color for each Argb value, while finding the per hue count
	// and total count.
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
			neighborHue := num.NormalizeDegreeInt(i)
			hueExcitedProportions[neighborHue] += proportion
		}
	}

	// Scores each HCT color based on usage and chroma, while optionally
	// filtering out values that do not have enough chroma or usage.
	scoredHct := []scoredColor{}
	for hct := range slices.Values(colorsHct) {
		hue := num.NormalizeDegreeInt(int(math.Round(hct.Hue)))
		proportion := hueExcitedProportions[hue]

		if opts.Filter &&
			(hct.Chroma < cutoffChroma || proportion <= cutoffExcitedProportion) {
			continue
		}

		proportionScore := proportion * 100.0 * weightProportion

		var chromaWeight float64
		if hct.Chroma < targetChroma {
			chromaWeight = weightChromaBelow
		} else {
			chromaWeight = weightChromaAbove
		}

		chromaScore := (hct.Chroma - targetChroma) * chromaWeight
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
	for differenceDegrees := float64(90); differenceDegrees >= 15; differenceDegrees-- {
		chosenColors = []color.Hct{} // Clear the array

		for scored := range slices.Values(scoredHct) {
			duplicateHue := false

			for chosenHct := range slices.Values(chosenColors) {
				if num.DifferenceDegrees(
					scored.hct.Hue,
					chosenHct.Hue,
				) < differenceDegrees {
					duplicateHue = true
					break
				}
			}

			if !duplicateHue {
				chosenColors = append(chosenColors, scored.hct)
			}

			if len(chosenColors) >= opts.Limit {
				break
			}
		}

		if len(chosenColors) >= opts.Limit {
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
