package dynamic

import (
	"math"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dislike"
	"github.com/Nadim147c/material/palettes"
)

// IsFidelity returns whether the scheme is a fidelity scheme
func IsFidelity(scheme DynamicScheme) bool {
	return scheme.Variant == Fidelity || scheme.Variant == Content
}

// IsMonochrome returns whether the scheme is monochrome
func IsMonochrome(scheme DynamicScheme) bool {
	return scheme.Variant == Monochrome
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

type MaterialColorSpec interface {
	Background() *DynamicColor
	Error() *DynamicColor
	ErrorContainer() *DynamicColor
	ErrorDim() *DynamicColor
	HighestSurface(s DynamicScheme) *DynamicColor
	InverseOnSurface() *DynamicColor
	InversePrimary() *DynamicColor
	InverseSurface() *DynamicColor
	NeutralPaletteKeyColor() *DynamicColor
	NeutralVariantPaletteKeyColor() *DynamicColor
	OnBackground() *DynamicColor
	OnError() *DynamicColor
	OnErrorContainer() *DynamicColor
	OnPrimary() *DynamicColor
	OnPrimaryContainer() *DynamicColor
	OnPrimaryFixed() *DynamicColor
	OnPrimaryFixedVariant() *DynamicColor
	OnSecondary() *DynamicColor
	OnSecondaryContainer() *DynamicColor
	OnSecondaryFixed() *DynamicColor
	OnSecondaryFixedVariant() *DynamicColor
	OnSurface() *DynamicColor
	OnSurfaceVariant() *DynamicColor
	OnTertiary() *DynamicColor
	OnTertiaryContainer() *DynamicColor
	OnTertiaryFixed() *DynamicColor
	OnTertiaryFixedVariant() *DynamicColor
	Outline() *DynamicColor
	OutlineVariant() *DynamicColor
	Primary() *DynamicColor
	PrimaryContainer() *DynamicColor
	PrimaryDim() *DynamicColor
	PrimaryFixed() *DynamicColor
	PrimaryFixedDim() *DynamicColor
	PrimaryPaletteKeyColor() *DynamicColor
	Scrim() *DynamicColor
	Secondary() *DynamicColor
	SecondaryContainer() *DynamicColor
	SecondaryDim() *DynamicColor
	SecondaryFixed() *DynamicColor
	SecondaryFixedDim() *DynamicColor
	SecondaryPaletteKeyColor() *DynamicColor
	Shadow() *DynamicColor
	Surface() *DynamicColor
	SurfaceBright() *DynamicColor
	SurfaceContainer() *DynamicColor
	SurfaceContainerHigh() *DynamicColor
	SurfaceContainerHighest() *DynamicColor
	SurfaceContainerLow() *DynamicColor
	SurfaceContainerLowest() *DynamicColor
	SurfaceDim() *DynamicColor
	SurfaceTint() *DynamicColor
	SurfaceVariant() *DynamicColor
	Tertiary() *DynamicColor
	TertiaryContainer() *DynamicColor
	TertiaryDim() *DynamicColor
	TertiaryFixed() *DynamicColor
	TertiaryFixedDim() *DynamicColor
	TertiaryPaletteKeyColor() *DynamicColor
}

type MaterialColorSpec2021 struct{}

var _ MaterialColorSpec = (*MaterialColorSpec2021)(nil)

// HighestSurface returns the highest surface color based on dark mode
func (m MaterialColorSpec2021) HighestSurface(s DynamicScheme) *DynamicColor {
	if s.IsDark {
		return m.SurfaceBright()
	}
	return m.SurfaceDim()
}

func (m MaterialColorSpec2021) PrimaryPaletteKeyColor() *DynamicColor {
	return FromPalette(
		"primary_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		func(s DynamicScheme) float64 { return s.PrimaryPalette.KeyColor.Tone },
	)
}

func (m MaterialColorSpec2021) SecondaryPaletteKeyColor() *DynamicColor {
	return FromPalette(
		"secondary_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		func(s DynamicScheme) float64 { return s.SecondaryPalette.KeyColor.Tone },
	)
}

func (m MaterialColorSpec2021) TertiaryPaletteKeyColor() *DynamicColor {
	return FromPalette(
		"tertiary_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		func(s DynamicScheme) float64 { return s.TertiaryPalette.KeyColor.Tone },
	)
}

func (m MaterialColorSpec2021) NeutralPaletteKeyColor() *DynamicColor {
	return FromPalette(
		"neutral_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		func(s DynamicScheme) float64 { return s.NeutralPalette.KeyColor.Tone },
	)
}

func (m MaterialColorSpec2021) NeutralVariantPaletteKeyColor() *DynamicColor {
	return FromPalette(
		"neutral_variant_palette_key_color",
		func(s DynamicScheme) palettes.TonalPalette { return s.NeutralVariantPalette },
		func(s DynamicScheme) float64 { return s.NeutralVariantPalette.KeyColor.Tone },
	)
}

// The material spec

func (m MaterialColorSpec2021) Background() *DynamicColor {
	return &DynamicColor{
		Name: "background",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 6.0
			}
			return 98.0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) OnBackground() *DynamicColor {
	return &DynamicColor{
		Name: "on_background",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 90.0
			}
			return 10.0
		},
		IsBackground: false,
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 3.0, 4.5, 7.0)
		},
		Background: func(s DynamicScheme) *DynamicColor { return m.Background() },
	}
}

func (m MaterialColorSpec2021) Surface() *DynamicColor {
	return &DynamicColor{
		Name: "surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 6.0
			}
			return 98.0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceDim() *DynamicColor {
	return &DynamicColor{
		Name: "surface_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			cc := NewContrastCurve(87.0, 87.0, 80.0, 75.0).Get(s.ContrastLevel)
			if s.IsDark {
				return 6.0
			}
			return cc
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceBright() *DynamicColor {
	return &DynamicColor{
		Name: "surface_bright",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			cc := NewContrastCurve(24.0, 24.0, 29.0, 34.0).Get(s.ContrastLevel)
			if s.IsDark {
				return cc
			}
			return 98.0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceContainerLowest() *DynamicColor {
	return &DynamicColor{
		Name: "surface_container_lowest",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			cc := NewContrastCurve(4.0, 4.0, 2.0, 0).Get(s.ContrastLevel)
			if s.IsDark {
				return cc
			}
			return 100.0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceContainerLow() *DynamicColor {
	return &DynamicColor{
		Name: "surface_container_low",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return NewContrastCurve(10.0, 10.0, 11.0, 12.0).Get(s.ContrastLevel)
			} else {
				return NewContrastCurve(96.0, 96.0, 96.0, 95.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceContainer() *DynamicColor {
	return &DynamicColor{
		Name: "surface_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return NewContrastCurve(12.0, 12.0, 16.0, 20.0).Get(s.ContrastLevel)
			} else {
				return NewContrastCurve(94.0, 94.0, 92.0, 90.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceContainerHigh() *DynamicColor {
	return &DynamicColor{
		Name: "surface_container_high",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return NewContrastCurve(17.0, 17.0, 21.0, 25.0).Get(s.ContrastLevel)
			} else {
				return NewContrastCurve(92.0, 92.0, 88.0, 85.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) SurfaceContainerHighest() *DynamicColor {
	return &DynamicColor{
		Name: "surface_container_highest",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return NewContrastCurve(22.0, 22.0, 26.0, 30.0).Get(s.ContrastLevel)
			} else {
				return NewContrastCurve(90.0, 90.0, 84.0, 80.0).Get(s.ContrastLevel)
			}
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) OnSurface() *DynamicColor {
	return &DynamicColor{
		Name: "on_surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 90.0
			}
			return 10.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) SurfaceVariant() *DynamicColor {
	return &DynamicColor{
		Name: "surface_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralVariantPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 30.0
			}
			return 90.0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) OnSurfaceVariant() *DynamicColor {
	return &DynamicColor{
		Name: "on_surface_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralVariantPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 80.0
			}
			return 30.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) InverseSurface() *DynamicColor {
	return &DynamicColor{
		Name: "inverse_surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 90.0
			}
			return 20.0
		},
		IsBackground: false,
	}
}

func (m MaterialColorSpec2021) InverseOnSurface() *DynamicColor {
	return &DynamicColor{
		Name: "inverse_on_surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 20.0
			}
			return 95.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.InverseSurface()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) Outline() *DynamicColor {
	return &DynamicColor{
		Name: "outline",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralVariantPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 60.0
			}
			return 50.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.5, 3.0, 4.5, 7.0)
		},
	}
}

func (m MaterialColorSpec2021) OutlineVariant() *DynamicColor {
	return &DynamicColor{
		Name: "outline_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralVariantPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 30.0
			}
			return 80.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
	}
}

func (m MaterialColorSpec2021) Shadow() *DynamicColor {
	return &DynamicColor{
		Name: "shadow",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone:         func(s DynamicScheme) float64 { return 0 },
		IsBackground: false,
	}
}

func (m MaterialColorSpec2021) Scrim() *DynamicColor {
	return &DynamicColor{
		Name: "scrim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		Tone:         func(s DynamicScheme) float64 { return 0 },
		IsBackground: false,
	}
}

func (m MaterialColorSpec2021) SurfaceTint() *DynamicColor {
	return &DynamicColor{
		Name: "surface_tint",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 80.0
			}
			return 40.0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2021) Primary() *DynamicColor {
	return &DynamicColor{
		Name: "primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 100.0
				}
				return 0.0
			}
			if s.IsDark {
				return 80.0
			}
			return 80.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 7.0)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryContainer(), m.Primary(), 10, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) OnPrimary() *DynamicColor {
	return &DynamicColor{
		Name: "on_primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 10.0
				}
				return 90.0
			}
			if s.IsDark {
				return 20.0
			}
			return 100.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.Primary()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) PrimaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "primary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsFidelity(s) {
				return s.SourceColorHct.Tone
			}
			if IsMonochrome(s) {
				if s.IsDark {
					return 85.0
				}
				return 25.0
			}
			if s.IsDark {
				return 30.0
			}
			return 90.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryContainer(), m.Primary(), 10, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) PrimaryDim() *DynamicColor {
	return nil
}

func (m MaterialColorSpec2021) OnPrimaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "on_primary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsFidelity(s) {
				return ForegroundTone(m.PrimaryContainer().GetTone(s), 4.5)
			}
			if IsMonochrome(s) {
				if s.IsDark {
					return 0.0
				}
				return 100.0
			}
			if s.IsDark {
				return 90.0
			}
			return 30.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) InversePrimary() *DynamicColor {
	return &DynamicColor{
		Name: "inverse_primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 40.0
			}
			return 80.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.InverseSurface()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 7.0)
		},
	}
}

func (m MaterialColorSpec2021) Secondary() *DynamicColor {
	return &DynamicColor{
		Name: "secondary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 80.0
			}
			return 40.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 7.0)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryContainer(), m.Secondary(), 10, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) OnSecondary() *DynamicColor {
	return &DynamicColor{
		Name: "on_secondary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 10.0
				}
				return 100.0
			}
			if s.IsDark {
				return 20.0
			}
			return 100.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.Secondary()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) SecondaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "secondary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			initialTone := 90.0
			if s.IsDark {
				initialTone = 30.0
			}
			if IsMonochrome(s) {
				if s.IsDark {
					return 30.0
				}
				return 85.0
			}
			if !IsFidelity(s) {
				return initialTone
			}
			return FindDesiredChromaByTone(
				s.SecondaryPalette.Hue,
				s.SecondaryPalette.Chroma,
				initialTone,
				!s.IsDark,
			)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryContainer(), m.Secondary(), 10, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) SecondaryDim() *DynamicColor {
	return nil
}

func (m MaterialColorSpec2021) OnSecondaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "on_secondary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 90.0
				}
				return 10.0
			}
			if !IsFidelity(s) {
				if s.IsDark {
					return 90.0
				}
				return 30.0
			}
			return ForegroundTone(m.SecondaryContainer().Tone(s), 4.5)
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) Tertiary() *DynamicColor {
	return &DynamicColor{
		Name: "tertiary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 90.0
				}
				return 25.0
			}
			if s.IsDark {
				return 80.0
			}
			return 40.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 7.0)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryContainer(), m.Tertiary(), 10.0, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) OnTertiary() *DynamicColor {
	return &DynamicColor{
		Name: "on_tertiary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 10.0
				}
				return 90.0
			}
			if s.IsDark {
				return 20.0
			}
			return 100.0
		},
		IsBackground: false,
		Background:   func(s DynamicScheme) *DynamicColor { return m.Tertiary() },
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) TertiaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "tertiary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 60.0
				}
				return 49.0
			}
			if !IsFidelity(s) {
				if s.IsDark {
					return 30.0
				}
				return 90.0
			}
			proposed := s.TertiaryPalette.Tone(s.SourceColorHct.Tone).ToHct()
			return dislike.FixIfDisliked(proposed).Tone
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryContainer(), m.Tertiary(), 10.0, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) TertiaryDim() *DynamicColor {
	return nil
}

func (m MaterialColorSpec2021) OnTertiaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "on_tertiary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 0.0
				}
				return 100.0
			}
			if !IsFidelity(s) {
				if s.IsDark {
					return 90.0
				}
				return 30.0
			}
			return ForegroundTone(m.TertiaryContainer().Tone(s), 4.5)
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) Error() *DynamicColor {
	return &DynamicColor{
		Name: "error",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.ErrorPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 80.0
			}
			return 40.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 7.0)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.ErrorContainer(), m.Error(), 10.0, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) OnError() *DynamicColor {
	return &DynamicColor{
		Name: "on_error",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.ErrorPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 20.0
			}
			return 100.0
		},
		IsBackground: false,
		Background:   func(s DynamicScheme) *DynamicColor { return m.Error() },
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) ErrorContainer() *DynamicColor {
	return &DynamicColor{
		Name: "error_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.ErrorPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 30.0
			}
			return 90.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.ErrorContainer(), m.Error(), 10.0, ToneNearer, false)
		},
	}
}

func (m MaterialColorSpec2021) ErrorDim() *DynamicColor {
	return nil
}

func (m MaterialColorSpec2021) OnErrorContainer() *DynamicColor {
	return &DynamicColor{
		Name: "on_error_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.ErrorPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				if s.IsDark {
					return 90.0
				}
				return 10.0
			}
			if s.IsDark {
				return 90.0
			}
			return 30.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.ErrorContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) PrimaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "primary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 40.0
			} else {
				return 90.0
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryFixed(), m.PrimaryFixedDim(), 10.0, ToneLighter, true)
		},
	}
}

func (m MaterialColorSpec2021) PrimaryFixedDim() *DynamicColor {
	return &DynamicColor{
		Name: "primary_fixed_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 30.0
			}
			return 80.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryFixed(), m.PrimaryFixedDim(), 10.0, ToneLighter, true)
		},
	}
}

func (m MaterialColorSpec2021) OnPrimaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "on_primary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 100.0
			}
			return 10.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryFixedDim()
		},
		SecondBackground: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryFixed()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) OnPrimaryFixedVariant() *DynamicColor {
	return &DynamicColor{
		Name: "on_primary_fixed_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 90.0
			}
			return 30.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryFixedDim()
		},
		SecondBackground: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryFixed()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) SecondaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "secondary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 80.0
			}
			return 90.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryFixed(), m.SecondaryFixedDim(), 10.0, ToneLighter, true)
		},
	}
}

func (m MaterialColorSpec2021) SecondaryFixedDim() *DynamicColor {
	return &DynamicColor{
		Name: "secondary_fixed_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 70.0
			}
			return 80.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.SecondaryFixed(), m.SecondaryFixedDim(), 10.0, ToneLighter, true)
		},
	}
}

func (m MaterialColorSpec2021) OnSecondaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "on_secondary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone:         func(s DynamicScheme) float64 { return 10.0 },
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryFixedDim()
		},
		SecondBackground: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryFixed()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) OnSecondaryFixedVariant() *DynamicColor {
	return &DynamicColor{
		Name: "on_secondary_fixed_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 25.0
			}
			return 30.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryFixedDim()
		},
		SecondBackground: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryFixed()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}

func (m MaterialColorSpec2021) TertiaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "tertiary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 40.0
			}
			return 90.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryFixed(), m.TertiaryFixedDim(), 10.0, ToneLighter, true)
		},
	}
}

func (m MaterialColorSpec2021) TertiaryFixedDim() *DynamicColor {
	return &DynamicColor{
		Name: "tertiary_fixed_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 30.0
			}
			return 80.0
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.HighestSurface(s)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(1.0, 1.0, 3.0, 4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.TertiaryFixed(), m.TertiaryFixedDim(), 10.0, ToneLighter, true)
		},
	}
}

func (m MaterialColorSpec2021) OnTertiaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "on_tertiary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 100.0
			}
			return 10.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryFixedDim()
		},
		SecondBackground: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryFixed()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(4.5, 7.0, 11.0, 21.0)
		},
	}
}

func (m MaterialColorSpec2021) OnTertiaryFixedVariant() *DynamicColor {
	return &DynamicColor{
		Name: "on_tertiary_fixed_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if IsMonochrome(s) {
				return 90.0
			}
			return 30.0
		},
		IsBackground: false,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryFixedDim()
		},
		SecondBackground: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryFixed()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return NewContrastCurve(3.0, 4.5, 7.0, 11.0)
		},
	}
}
