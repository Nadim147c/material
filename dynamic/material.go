package dynamic

import (
	"math"

	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/goyou/dislike"
	"github.com/Nadim147c/goyou/palettes"
)

// IsFidelity returns whether the scheme is a fidelity scheme
func IsFidelity(scheme DynamicScheme) bool {
	return scheme.Variant == Fidelity || scheme.Variant == Content
}

// IsMonochrome returns whether the scheme is monochrome
func IsMonochrome(scheme DynamicScheme) bool {
	return scheme.Variant == Monochrome
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

// FindDesiredChromaByTone finds a tone where the chroma is as close as possible to the requested value
func FindDesiredChromaByTone(hue, chroma, tone float64, byDecreasingTone bool) float64 {
	answer := tone

	closestToChroma := color.NewHct(hue, chroma, tone)
	if closestToChroma.Chroma < chroma {
		chromaPeak := closestToChroma.Chroma
		for closestToChroma.Chroma < chroma {
			if byDecreasingTone {
				answer -= 1.0
			} else {
				answer += 1.0
			}
			potentialSolution := color.NewHct(hue, chroma, answer)
			if chromaPeak > potentialSolution.Chroma {
				break
			}
			if math.Abs(potentialSolution.Chroma-chroma) < 0.4 {
				break
			}

			potentialDelta := math.Abs(potentialSolution.Chroma - chroma)
			currentDelta := math.Abs(closestToChroma.Chroma - chroma)
			if potentialDelta < currentDelta {
				closestToChroma = potentialSolution
			}
			chromaPeak = math.Max(chromaPeak, potentialSolution.Chroma)
		}
	}

	return answer
}

const contentAccentToneDelta = 15.0

// HighestSurface returns the highest surface color based on dark mode
func HighestSurface(s DynamicScheme) DynamicColor {
	if s.IsDark {
		return DynamicSchemeProvider.SurfaceBright()
	}
	return DynamicSchemeProvider.SurfaceDim()
}

// DynamicSchemeProvider provides predefined dynamic color schemes
type MaterialColor struct{}

// DynamicSchemeProvider is the singleton instance of DynamicSchemeProviderProvider
var DynamicSchemeProvider = MaterialColor{}

func (m MaterialColor) PrimaryPaletteKeyColor() DynamicColor {
	return FromPalette(
		"primary_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		func(s DynamicScheme) float64 { return s.PrimaryPalette.KeyColor.Tone },
	)
}

func (m MaterialColor) SecondaryPaletteKeyColor() DynamicColor {
	return FromPalette(
		"secondary_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		func(s DynamicScheme) float64 { return s.SecondaryPalette.KeyColor.Tone },
	)
}

func (m MaterialColor) TertiaryPaletteKeyColor() DynamicColor {
	return FromPalette(
		"tertiary_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		func(s DynamicScheme) float64 { return s.TertiaryPalette.KeyColor.Tone },
	)
}

func (m MaterialColor) NeutralPaletteKeyColor() DynamicColor {
	return FromPalette(
		"neutral_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		func(s DynamicScheme) float64 { return s.NeutralPalette.KeyColor.Tone },
	)
}

func (m MaterialColor) NeutralVariantPaletteKeyColor() DynamicColor {
	return FromPalette(
		"neutral_variant_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.NeutralVariantPalette },
		func(s DynamicScheme) float64 { return s.NeutralVariantPalette.KeyColor.Tone },
	)
}

// The material spec

func (m MaterialColor) Background() DynamicColor {
	return DynamicColor{
		Name:         "background",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:         func(s DynamicScheme) float64 { return ternary(s.IsDark, 6.0, 98.0) },
		IsBackground: true,
	}
}

func (m MaterialColor) OnBackground() DynamicColor {
	return DynamicColor{
		Name:          "on_background",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 98.0, 10.0) },
		IsBackground:  false,
		ContrastCurve: ContrastCurve(3.0, 3.0, 4.5, 7.0),
		Background:    func(s DynamicScheme) DynamicColor { return m.Background() },
	}
}

func (M MaterialColor) Surface() DynamicColor {
	return DynamicColor{
		Name:         "surface",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:         func(s DynamicScheme) float64 { return ternary(s.IsDark, 6.0, 98.0) },
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceDim() DynamicColor {
	return DynamicColor{
		Name:    "surface_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			return ternary(s.IsDark, 6.0, ContrastCurve(87.0, 87.0, 80.0, 75.0).Get(s.ContrastLevel))
		},
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceBright() DynamicColor {
	return DynamicColor{
		Name:    "surface_bright",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			return ternary(s.IsDark, ContrastCurve(24.0, 24.0, 29.0, 34.0).Get(s.ContrastLevel), 98.0)
		},
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceContainerLowest() DynamicColor {
	return DynamicColor{
		Name:    "surface_container_lowest",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			return ternary(s.IsDark, ContrastCurve(4.0, 4.0, 2.0, 0).Get(s.ContrastLevel), 100.0)
		},
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceContainerLow() DynamicColor {
	return DynamicColor{
		Name:    "surface_container_low",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return ContrastCurve(10.0, 10.0, 11.0, 12.0).Get(s.ContrastLevel)
			} else {
				return ContrastCurve(96.0, 96.0, 96.0, 95.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceContainer() DynamicColor {
	return DynamicColor{
		Name:    "surface_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return ContrastCurve(12.0, 12.0, 16.0, 20.0).Get(s.ContrastLevel)
			} else {
				return ContrastCurve(94.0, 94.0, 92.0, 90.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceContainerHigh() DynamicColor {
	return DynamicColor{
		Name:    "surface_container_high",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return ContrastCurve(17.0, 17.0, 21.0, 25.0).Get(s.ContrastLevel)
			} else {
				return ContrastCurve(92.0, 92.0, 88.0, 85.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColor) SurfaceContainerHighest() DynamicColor {
	return DynamicColor{
		Name:    "surface_container_highest",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return ContrastCurve(22.0, 22.0, 26.0, 30.0).Get(s.ContrastLevel)
			} else {
				return ContrastCurve(90.0, 90.0, 84.0, 80.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColor) OnSurface() DynamicColor {
	return DynamicColor{
		Name:          "on_surface",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 90.0, 10.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) SurfaceVariant() DynamicColor {
	return DynamicColor{
		Name:         "surface_variant",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.NeutralVariantPalette },
		Tone:         func(s DynamicScheme) float64 { return ternary(s.IsDark, 30.0, 90.0) },
		IsBackground: true,
	}
}

func (m MaterialColor) OnSurfaceVariant() DynamicColor {
	return DynamicColor{
		Name:          "on_surface_variant",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.NeutralVariantPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 80.0, 30.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) InverseSurface() DynamicColor {
	return DynamicColor{
		Name:         "inverse_surface",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:         func(s DynamicScheme) float64 { return ternary(s.IsDark, 90.0, 20.0) },
		IsBackground: false,
	}
}

func (m MaterialColor) InverseOnSurface() DynamicColor {
	return DynamicColor{
		Name:          "inverse_on_surface",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 20.0, 95.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.InverseSurface() },
		ContrastCurve: ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) Outline() DynamicColor {
	return DynamicColor{
		Name:          "outline",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.NeutralVariantPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 60.0, 50.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.5, 3.0, 4.5, 7.0),
	}
}

func (m MaterialColor) OutlineVariant() DynamicColor {
	return DynamicColor{
		Name:          "outline_variant",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.NeutralVariantPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 30.0, 80.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
	}
}

func (m MaterialColor) Shadow() DynamicColor {
	return DynamicColor{
		Name:         "shadow",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:         func(s DynamicScheme) float64 { return 0 },
		IsBackground: false,
	}
}

func (m MaterialColor) Scrim() DynamicColor {
	return DynamicColor{
		Name:         "scrim",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone:         func(s DynamicScheme) float64 { return 0 },
		IsBackground: false,
	}
}

func (m MaterialColor) SurfaceTint() DynamicColor {
	return DynamicColor{
		Name:         "surface_tint",
		Palette:      func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone:         func(s DynamicScheme) float64 { return ternary(s.IsDark, 80.0, 40.0) },
		IsBackground: true,
	}
}

func (m MaterialColor) Primary() DynamicColor {
	return DynamicColor{
		Name:    "primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 100.0, 0.0)
			}
			return ternary(s.IsDark, 80.0, 80.0)
		},
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 7.0),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryContainer(), m.Primary(), 10, Nearer, false)
		},
	}
}

func (m MaterialColor) OnPrimary() DynamicColor {
	return DynamicColor{
		Name:    "on_primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 10.0, 90.0)
			}
			return ternary(s.IsDark, 20.0, 100.0)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.Primary() },
		ContrastCurve: ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) PrimaryContainer() DynamicColor {
	return DynamicColor{
		Name:    "primary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsFidelity(s) {
				return s.SourceColorHct.Tone
			}
			if IsMonochrome(s) {
				return ternary(s.IsDark, 85.0, 25.0)
			}
			return ternary(s.IsDark, 30.0, 90.0)
		},
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryContainer(), m.Primary(), 10, Nearer, false)
		},
	}
}

func (m MaterialColor) OnPrimaryContainer() DynamicColor {
	return DynamicColor{
		Name:    "on_primary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsFidelity(s) {
				return ForegroundTone(m.PrimaryContainer().GetTone(s), 4.5)
			}
			if IsMonochrome(s) {
				return ternary(s.IsDark, 0.0, 100.0)
			}
			return ternary(s.IsDark, 90.0, 30.0)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.PrimaryContainer() },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) InversePrimary() DynamicColor {
	return DynamicColor{
		Name:          "inverse_primary",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 40.0, 80.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.InverseSurface() },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 7.0),
	}
}

func (m MaterialColor) Secondary() DynamicColor {
	return DynamicColor{
		Name:          "secondary",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 80.0, 40.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 7.0),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryContainer(), m.Secondary(), 10, Nearer, false)
		},
	}
}

func (m MaterialColor) OnSecondary() DynamicColor {
	return DynamicColor{
		Name:    "on_secondary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 10.0, 100.0)
			}
			return ternary(s.IsDark, 20.0, 100.0)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.Secondary() },
		ContrastCurve: ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) SecondaryContainer() DynamicColor {
	return DynamicColor{
		Name:    "secondary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			initialTone := ternary(s.IsDark, 30.0, 90.0)
			if IsMonochrome(s) {
				return ternary(s.IsDark, 30.0, 85.0)
			}
			if !IsFidelity(s) {
				return initialTone
			}
			return FindDesiredChromaByTone(
				s.SecondaryPalette.Hue,
				s.SecondaryPalette.Chroma,
				initialTone,
				ternary(s.IsDark, false, true),
			)
		},
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryContainer(), m.Secondary(), 10, Nearer, false)
		},
	}
}

func (m MaterialColor) OnSecondaryContainer() DynamicColor {
	return DynamicColor{
		Name:    "on_secondary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 90.0, 10.0)
			}
			if !IsFidelity(s) {
				return ternary(s.IsDark, 90.0, 30.0)
			}
			return ForegroundTone(m.SecondaryContainer().Tone(s), 4.5)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.SecondaryContainer() },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) Tertiary() DynamicColor {
	return DynamicColor{
		Name:    "tertiary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 90.0, 25.0)
			}
			return ternary(s.IsDark, 80.0, 40.0)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) DynamicColor {
			return HighestSurface(s)
		},
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 7.0),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryContainer(), m.Tertiary(), 10.0, Nearer, false)
		},
	}
}

func (m MaterialColor) OnTertiary() DynamicColor {
	return DynamicColor{
		Name:    "on_tertiary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 10.0, 90.0)
			}
			return ternary(s.IsDark, 20.0, 100.0)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.Tertiary() },
		ContrastCurve: ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) TertiaryContainer() DynamicColor {
	return DynamicColor{
		Name:    "tertiary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 60.0, 49.0)
			}
			if !IsFidelity(s) {
				return ternary(s.IsDark, 30.0, 90.0)
			}
			proposed := s.TertiaryPalette.Tone(s.SourceColorHct.Tone).ToHct()
			return dislike.FixIfDisliked(proposed).Tone
		},
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryContainer(), m.Tertiary(), 10.0, Nearer, false)
		},
	}
}

func (m MaterialColor) OnTertiaryContainer() DynamicColor {
	return DynamicColor{
		Name:    "on_tertiary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 0.0, 100.0)
			}
			if !IsFidelity(s) {
				return ternary(s.IsDark, 90.0, 30.0)
			}
			return ForegroundTone(m.TertiaryContainer().Tone(s), 4.5)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.TertiaryContainer() },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) Error() DynamicColor {
	return DynamicColor{
		Name:          "error",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 80.0, 40.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 7.0),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.ErrorContainer(), m.Error(), 10.0, Nearer, false)
		},
	}
}

func (m MaterialColor) OnError() DynamicColor {
	return DynamicColor{
		Name:          "on_error",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 20.0, 100.0) },
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.Error() },
		ContrastCurve: ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) ErrorContainer() DynamicColor {
	return DynamicColor{
		Name:          "error_container",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(s.IsDark, 30.0, 90.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.ErrorContainer(), m.Error(), 10.0, Nearer, false)
		},
	}
}

func (m MaterialColor) OnErrorContainer() DynamicColor {
	return DynamicColor{
		Name:    "on_error_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return ternary(s.IsDark, 90.0, 10.0)
			}
			return ternary(s.IsDark, 90.0, 30.0)
		},
		IsBackground:  false,
		Background:    func(s DynamicScheme) DynamicColor { return m.ErrorContainer() },
		ContrastCurve: ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) PrimaryFixed() DynamicColor {
	return DynamicColor{
		Name:          "primary_fixed",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 40.0, 90.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryFixed(), m.PrimaryFixedDim(), 10.0, Lighter, true)
		},
	}
}

func (m MaterialColor) PrimaryFixedDim() DynamicColor {
	return DynamicColor{
		Name:          "primary_fixed_dim",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 30.0, 80.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryFixed(), m.PrimaryFixedDim(), 10.0, Lighter, true)
		},
	}
}

func (m MaterialColor) OnPrimaryFixed() DynamicColor {
	return DynamicColor{
		Name:             "on_primary_fixed",
		Palette:          func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone:             func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 100.0, 10.0) },
		IsBackground:     false,
		Background:       func(s DynamicScheme) DynamicColor { return m.PrimaryFixedDim() },
		SecondBackground: func(s DynamicScheme) DynamicColor { return m.PrimaryFixed() },
		ContrastCurve:    ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) OnPrimaryFixedVariant() DynamicColor {
	return DynamicColor{
		Name:             "on_primary_fixed_variant",
		Palette:          func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone:             func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 90.0, 30.0) },
		IsBackground:     false,
		Background:       func(s DynamicScheme) DynamicColor { return m.PrimaryFixedDim() },
		SecondBackground: func(s DynamicScheme) DynamicColor { return m.PrimaryFixed() },
		ContrastCurve:    ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) SecondaryFixed() DynamicColor {
	return DynamicColor{
		Name:          "secondary_fixed",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 80.0, 90.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryFixed(), m.SecondaryFixedDim(), 10.0, Lighter, true)
		},
	}
}

func (m MaterialColor) SecondaryFixedDim() DynamicColor {
	return DynamicColor{
		Name:          "secondary_fixed_dim",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 70.0, 80.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryFixed(), m.SecondaryFixedDim(), 10.0, Lighter, true)
		},
	}
}

func (m MaterialColor) OnSecondaryFixed() DynamicColor {
	return DynamicColor{
		Name:             "on_secondary_fixed",
		Palette:          func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone:             func(s DynamicScheme) float64 { return 10.0 },
		IsBackground:     false,
		Background:       func(s DynamicScheme) DynamicColor { return m.SecondaryFixedDim() },
		SecondBackground: func(s DynamicScheme) DynamicColor { return m.SecondaryFixed() },
		ContrastCurve:    ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) OnSecondaryFixedVariant() DynamicColor {
	return DynamicColor{
		Name:             "on_secondary_fixed_variant",
		Palette:          func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone:             func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 25.0, 30.0) },
		IsBackground:     false,
		Background:       func(s DynamicScheme) DynamicColor { return m.SecondaryFixedDim() },
		SecondBackground: func(s DynamicScheme) DynamicColor { return m.SecondaryFixed() },
		ContrastCurve:    ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}

func (m MaterialColor) TertiaryFixed() DynamicColor {
	return DynamicColor{
		Name:          "tertiary_fixed",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 40.0, 90.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryFixed(), m.TertiaryFixedDim(), 10.0, Lighter, true)
		},
	}
}

func (m MaterialColor) TertiaryFixedDim() DynamicColor {
	return DynamicColor{
		Name:          "tertiary_fixed_dim",
		Palette:       func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone:          func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 30.0, 80.0) },
		IsBackground:  true,
		Background:    func(s DynamicScheme) DynamicColor { return HighestSurface(s) },
		ContrastCurve: ContrastCurve(1.0, 1.0, 3.0, 4.5),
		ToneDeltaPair: func(s DynamicScheme) ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryFixed(), m.TertiaryFixedDim(), 10.0, Lighter, true)
		},
	}
}

func (m MaterialColor) OnTertiaryFixed() DynamicColor {
	return DynamicColor{
		Name:             "on_tertiary_fixed",
		Palette:          func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone:             func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 100.0, 10.0) },
		IsBackground:     false,
		Background:       func(s DynamicScheme) DynamicColor { return m.TertiaryFixedDim() },
		SecondBackground: func(s DynamicScheme) DynamicColor { return m.TertiaryFixed() },
		ContrastCurve:    ContrastCurve(4.5, 7.0, 11.0, 21.0),
	}
}

func (m MaterialColor) OnTertiaryFixedVariant() DynamicColor {
	return DynamicColor{
		Name:             "on_tertiary_fixed_variant",
		Palette:          func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone:             func(s DynamicScheme) float64 { return ternary(IsMonochrome(s), 90.0, 30.0) },
		IsBackground:     false,
		Background:       func(s DynamicScheme) DynamicColor { return m.TertiaryFixedDim() },
		SecondBackground: func(s DynamicScheme) DynamicColor { return m.TertiaryFixed() },
		ContrastCurve:    ContrastCurve(3.0, 4.5, 7.0, 11.0),
	}
}
