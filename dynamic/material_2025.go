package dynamic

import (
	"github.com/Nadim147c/goyou/num"
	"github.com/Nadim147c/goyou/palettes"
)

// tMaxC
// Paramters:
//
//	lowerBound: 0
//	upperBound: 100
//	chromaMultiplier: 1
func tMaxC(palette palettes.TonalPalette, lowerBound float64, upperBound float64, chromaMultiplier float64) float64 {
	answer := FindDesiredChromaByTone(palette.Hue, palette.Chroma*chromaMultiplier, 100, true)
	return num.Clamp(lowerBound, upperBound, answer)
}

// tMaxC
// Paramters:
//
//	lowerBound: 0
//	upperBound: 100
func tMinC(palette palettes.TonalPalette, lowerBound float64, upperBound float64) float64 {
	answer := FindDesiredChromaByTone(palette.Hue, palette.Chroma, 0, false)
	return num.Clamp(lowerBound, upperBound, answer)
}

// GetCurve returns the contrast curve for a given default contrast.
func GetCurve(defaultContrast float64) *ContrastCurve {
	switch defaultContrast {
	case 1.5:
		return NewContrastCurve(1.5, 1.5, 3, 4.5)
	case 3:
		return NewContrastCurve(3, 3, 4.5, 7)
	case 4.5:
		return NewContrastCurve(4.5, 4.5, 7, 11)
	case 6:
		return NewContrastCurve(6, 6, 7, 11)
	case 7:
		return NewContrastCurve(7, 7, 11, 21)
	case 9:
		return NewContrastCurve(9, 9, 11, 21)
	case 11:
		return NewContrastCurve(11, 11, 21, 21)
	case 21:
		return NewContrastCurve(21, 21, 21, 21)
	default:
		return NewContrastCurve(defaultContrast, defaultContrast, 7, 21)
	}
}

type MaterialColorSpec2025 struct {
	MaterialColorSpec2021
}

var _ MaterialColorSpec = (*MaterialColorSpec2025)(nil)

func (m MaterialColorSpec2025) Surface() *DynamicColor {
	return &DynamicColor{
		Name:    "surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				if s.IsDark {
					return 4
				}

				if s.NeutralPalette.IsBlue() {
					return 99
				} else if s.Variant == Vibrant {
					return 97
				}

				return 98
			}
			return 0
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2025) SurfaceDim() *DynamicColor {
	return &DynamicColor{
		Name:    "surface_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 4
			}
			if s.NeutralPalette.IsYellow() {
				return 90
			} else if s.Variant == Vibrant {
				return 85
			}
			return 87
		},
		IsBackground: true,

		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.IsDark {
				switch s.Variant {
				case Neutral:
					return 2.5
				case TonalSpot:
					return 1.7
				case Expressive:
					return ternary(s.NeutralPalette.IsBlue(), 2.7, 1.75)
				case Vibrant:
					return 1.36
				}
			}
			return 1
		},
	}
}

func (m MaterialColorSpec2025) SurfaceBright() *DynamicColor {
	return &DynamicColor{
		Name:    "surface_bright",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 18
			}
			if s.NeutralPalette.IsBlue() {
				return 99
			} else if s.Variant == Vibrant {
				return 97
			}
			return 98
		},
		IsBackground: true,
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if !s.IsDark {
				return 1
			}
			switch s.Variant {
			case Neutral:
				return 2.5
			case TonalSpot:
				return 1.7
			case Vibrant:
				return 1.36
			default:
				return 1
			}
		},
	}
}

func (m MaterialColorSpec2025) SurfaceContainerLowest() *DynamicColor {
	return &DynamicColor{
		Name: "surface_container_lowest",
		// palette: (s) => s.neutralPalette,
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			return ternary(s.IsDark, 0.0, 100.0)
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2025) SurfaceContainerLow() *DynamicColor {
	return &DynamicColor{
		Name:    "surface_container_low",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform != Phone {
				return 15
			}
			if s.IsDark {
				return 6
			} else if s.NeutralPalette.IsYellow() {
				return 98
			} else if s.Variant == Vibrant {
				return 95
			}
			return 96
		},
		IsBackground: true,
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform != Phone {
				return 1
			}
			switch s.Variant {
			case Neutral:
				return 1.3
			case TonalSpot:
				return 1.25
			case Vibrant:
				return 1.08
			default:
				return 1
			}
		},
	}
}

func (m MaterialColorSpec2025) SurfaceContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "surface_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform != Phone {
				return 20
			}
			if s.IsDark {
				return 9
			}
			if s.NeutralPalette.IsYellow() {
				return 96
			} else if s.Variant == Vibrant {
				return 92
			}
			return 94
		},
		IsBackground: true,
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform != Phone {
				return 1
			}
			switch s.Variant {
			case Neutral:
				return 1.6
			case TonalSpot:
				return 1.4
			case Expressive:
				return ternary(s.NeutralPalette.IsYellow(), 1.6, 1.3)
			case Vibrant:
				return 1.15
			default:
				return 1
			}
		},
	}
}

func (m MaterialColorSpec2025) SurfaceContainerHigh() *DynamicColor {
	return &DynamicColor{
		Name:    "surface_container_high",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				if s.IsDark {
					return 12
				} else {
					if s.NeutralPalette.IsYellow() {
						return 94
					} else if s.Variant == Vibrant {
						return 90
					} else {
						return 92
					}
				}
			} else {
				return 25
			}
		},
		IsBackground: true,
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 1.9
				case TonalSpot:
					return 1.5
				case Expressive:
					return ternary(s.NeutralPalette.IsYellow(), 1.95, 1.45)
				case Vibrant:
					return 1.22
				}
			}
			return 1
		},
	}
}

func (m MaterialColorSpec2025) SurfaceContainerHighest() *DynamicColor {
	return &DynamicColor{
		Name:    "surface_container_highest",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 15
			} else {
				if s.NeutralPalette.IsYellow() {
					return 92
				} else if s.Variant == Vibrant {
					return 88
				} else {
					return 90
				}
			}
		},
		IsBackground: true,
		ChromaMultiplier: func(s DynamicScheme) float64 {
			switch s.Variant {
			case Neutral:
				return 2.2
			case TonalSpot:
				return 1.7
			case Expressive:
				return ternary(s.NeutralPalette.IsYellow(), 2.3, 1.6)
			case Vibrant:
				return 1.29
			default:
				return 1
			}
		},
	}
}

func (m MaterialColorSpec2025) OnSurface() *DynamicColor {
	return &DynamicColor{
		Name:    "on_surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Variant == Vibrant {
				return tMaxC(s.NeutralPalette, 0, 100, 1.1)
			} else {
				// For all other variants, the initial tone should be the default
				// tone, which is the same as the background color.
				return GetInitialToneFromBackground(func(s DynamicScheme) *DynamicColor {
					if s.Platform == Phone {
						return m.HighestSurface(s)
					} else {
						return m.SurfaceContainerHigh()
					}
				})(s)
			}
		},
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						return ternary(s.IsDark, 3.0, 2.3)
					}
					return 1.6
				}
			}
			return 1
		},
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainerHigh()
			}
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.IsDark {
				return GetCurve(11)
			}
			return GetCurve(9)
		},
	}
}

func (m MaterialColorSpec2025) OnSurfaceVariant() *DynamicColor {
	return &DynamicColor{
		Name:    "on_surface_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						return ternary(s.IsDark, 3.0, 2.3)
					}
					return 1.6
				}
			}
			return 1
		},
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainerHigh()
			}
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) Outline() *DynamicColor {
	return &DynamicColor{
		Name: "outline",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						return ternary(s.IsDark, 3.0, 2.3)
					} else {
						return 1.6
					}
				}
			}
			return 1
		},
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainer()
			}
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return ternary(s.Platform == Phone, GetCurve(3), GetCurve(4.5))
		},
	}
}

func (m MaterialColorSpec2025) OutlineVariant() *DynamicColor {
	return &DynamicColor{
		Name: "outline_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.NeutralPalette
		},
		ChromaMultiplier: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				switch s.Variant {
				case Neutral:
					return 2.2
				case TonalSpot:
					return 1.7
				case Expressive:
					if s.NeutralPalette.IsYellow() {
						return ternary(s.IsDark, 3.0, 2.3)
					}
					return 1.6
				}
			}
			return 1
		},
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return ternary(s.Platform == Phone, GetCurve(1.5), GetCurve(3))
		},
	}
}

func (m MaterialColorSpec2025) InverseSurface() *DynamicColor {
	return &DynamicColor{
		Name:    "inverse_surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.IsDark {
				return 98
			}
			return 4
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2025) InverseOnSurface() *DynamicColor {
	return &DynamicColor{
		Name:    "inverse_on_surface",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.NeutralPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.InverseSurface()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) Primary() *DynamicColor {
	return &DynamicColor{
		Name:    "primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			switch s.Variant {
			case Neutral:
				if s.Platform == Phone {
					return ternary(s.IsDark, 80.0, 40.0)
				} else {
					return 90
				}
			case TonalSpot:
				if s.Platform == Phone {
					if s.IsDark {
						return 80
					} else {
						return tMaxC(s.PrimaryPalette, 0, 10, 1)
					}
				} else {
					return tMaxC(s.PrimaryPalette, 0, 90, 1)
				}
			case Expressive:
				if s.PrimaryPalette.IsYellow() {
					return tMaxC(s.PrimaryPalette, 0, 25, 1)
				} else if s.PrimaryPalette.IsCyan() {
					return tMaxC(s.PrimaryPalette, 0, 88, 1)
				} else {
					return tMaxC(s.PrimaryPalette, 0, 98, 1)
				}
			default: // VIBRANT
				if s.PrimaryPalette.IsCyan() {
					return tMaxC(s.PrimaryPalette, 0, 88, 1)
				} else {
					return tMaxC(s.PrimaryPalette, 0, 98, 1)
				}
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			} else {
				return m.SurfaceContainerHigh()
			}
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.PrimaryContainer(),
					m.Primary(),
					5,
					ToneRelativeDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) PrimaryDim() *DynamicColor {
	return &DynamicColor{
		Name: "primary_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.PrimaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			switch s.Variant {
			case Neutral:
				return 85
			case TonalSpot:
				return tMaxC(s.PrimaryPalette, 0, 90, 1)
			default:
				return tMaxC(s.PrimaryPalette, 0, 100, 1)
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(m.PrimaryDim(), m.Primary(), 5, ToneDarker, true)
		},
	}
}

func (m MaterialColorSpec2025) OnPrimary() *DynamicColor {
	return &DynamicColor{
		Name:    "on_primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Primary()
			} else {
				return m.PrimaryDim()
			}
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) PrimaryContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "primary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Watch {
				return 30
			} else if s.Variant == Neutral {
				return ternary(s.IsDark, 30.0, 90.0)
			} else if s.Variant == TonalSpot {
				if s.IsDark {
					return tMinC(s.PrimaryPalette, 35, 93)
				} else {
					return tMaxC(s.PrimaryPalette, 0, 90, 1)
				}
			} else if s.Variant == Expressive {
				if s.IsDark {
					return tMaxC(s.PrimaryPalette, 30, 93, 1)
				} else {
					if s.PrimaryPalette.IsCyan() {
						return tMaxC(s.PrimaryPalette, 78, 88, 1)
					} else {
						return tMaxC(s.PrimaryPalette, 78, 90, 1)
					}
				}
			} else { // VIBRANT
				if s.IsDark {
					return tMinC(s.PrimaryPalette, 66, 93)
				} else {
					if s.PrimaryPalette.IsCyan() {
						return tMaxC(s.PrimaryPalette, 66, 88, 1)
					} else {
						return tMaxC(s.PrimaryPalette, 66, 93, 1)
					}
				}
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{}
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return nil
			}
			return NewToneDeltaPair(
				m.PrimaryContainer(),
				m.PrimaryDim(),
				10,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) OnPrimaryContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "on_primary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) PrimaryFixed() *DynamicColor {
	return &DynamicColor{
		Name:    "primary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			s.IsDark = false
			return m.PrimaryContainer().GetTone(s)
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2025) PrimaryFixedDim() *DynamicColor {
	return &DynamicColor{
		Name:    "primary_fixed_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			return m.PrimaryFixed().GetTone(s)
		},
		IsBackground: true,
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.PrimaryFixedDim(),
				m.PrimaryFixed(),
				5,
				ToneDarker,
				true,
				ConstraintExact,
			)
		},
	}
}

func (m MaterialColorSpec2025) OnPrimaryFixed() *DynamicColor {
	return &DynamicColor{
		Name:    "on_primary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryFixedDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) OnPrimaryFixedVariant() *DynamicColor {
	return &DynamicColor{
		Name:    "on_primary_fixed_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.PrimaryFixedDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
	}
}

func (m MaterialColorSpec2025) InversePrimary() *DynamicColor {
	return &DynamicColor{
		Name:    "inverse_primary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.PrimaryPalette },
		Tone: func(s DynamicScheme) float64 {
			return tMaxC(s.PrimaryPalette, 0, 100, 1)
		},
		Background: func(s DynamicScheme) *DynamicColor {
			return m.InverseSurface()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) Secondary() *DynamicColor {
	return &DynamicColor{
		Name:    "secondary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Watch {
				if s.Variant == Neutral {
					return 90
				}
				return tMaxC(s.SecondaryPalette, 0, 90, 1)
			} else if s.Variant == Neutral {
				if s.IsDark {
					return tMinC(s.SecondaryPalette, 0, 98)
				}
				return tMaxC(s.SecondaryPalette, 0, 100, 1)
			} else if s.Variant == Vibrant {
				if s.IsDark {
					return tMaxC(s.SecondaryPalette, 0, 90, 1)
				}
				return tMaxC(s.SecondaryPalette, 0, 98, 1)
			} else { // EXPRESSIVE and TONAL_SPOT
				if s.IsDark {
					return 80
				}
				return tMaxC(s.SecondaryPalette, 0, 100, 1)
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.SecondaryContainer(),
					m.Secondary(),
					5,
					ToneRelativeLighter,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) SecondaryDim() *DynamicColor {
	return &DynamicColor{
		Name:    "secondary_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Variant == Neutral {
				return 85
			}
			return tMaxC(s.SecondaryPalette, 0, 90, 1)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.SecondaryDim(),
				m.Secondary(),
				5,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
	}
}

func (m MaterialColorSpec2025) OnSecondary() *DynamicColor {
	return &DynamicColor{
		Name:    "on_secondary",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Secondary()
			}
			return m.SecondaryDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) SecondaryContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "secondary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Watch {
				return 30
			} else if s.Variant == Vibrant {
				if s.IsDark {
					return tMinC(s.SecondaryPalette, 30, 40)
				}
				return tMaxC(s.SecondaryPalette, 84, 90, 1)
			} else if s.Variant == Expressive {
				if s.IsDark {
					return 15
				}
				return tMaxC(s.SecondaryPalette, 90, 95, 1)
			} else {
				if s.IsDark {
					return 25
				}
				return 90
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{}
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Watch {
				return NewToneDeltaPair(
					m.SecondaryContainer(),
					m.SecondaryDim(),
					10,
					ToneDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) OnSecondaryContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "on_secondary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) SecondaryFixed() *DynamicColor {
	return &DynamicColor{
		Name:    "secondary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			s.IsDark = false
			return m.SecondaryContainer().GetTone(s)
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2025) SecondaryFixedDim() *DynamicColor {
	return &DynamicColor{
		Name:    "secondary_fixed_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.SecondaryPalette },
		Tone: func(s DynamicScheme) float64 {
			return m.SecondaryFixed().GetTone(s)
		},
		IsBackground: true,
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.SecondaryFixedDim(),
				m.SecondaryFixed(),
				5,
				ToneDarker,
				true,
				ConstraintExact,
			)
		},
	}
}

func (m MaterialColorSpec2025) OnSecondaryFixed() *DynamicColor {
	return &DynamicColor{
		Name: "on_secondary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryFixedDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve { return GetCurve(7) },
	}
}

func (m MaterialColorSpec2025) OnSecondaryFixedVariant() *DynamicColor {
	return &DynamicColor{
		Name: "on_secondary_fixed_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.SecondaryPalette
		},
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SecondaryFixedDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
	}
}

func (m MaterialColorSpec2025) Tertiary() *DynamicColor {
	return &DynamicColor{
		Name: "tertiary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Watch {
				if s.Variant == TonalSpot {
					return tMaxC(s.TertiaryPalette, 0, 90, 1)
				}
				return tMaxC(s.TertiaryPalette, 0, 100, 1)
			} else if s.Variant == Expressive || s.Variant == Vibrant {
				limit := 100.0
				if s.TertiaryPalette.IsYellow() {
					limit = 88
				} else if s.IsDark {
					limit = 98
				}
				return tMaxC(s.TertiaryPalette, 0, limit, 1)
			} else {
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 98, 1)
				}
				return tMaxC(s.TertiaryPalette, 0, 100, 1)
			}
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.TertiaryContainer(),
					m.Tertiary(),
					5,
					ToneRelativeLighter,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) TertiaryDim() *DynamicColor {
	return &DynamicColor{
		Name: "tertiary_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Tone: func(s DynamicScheme) float64 {
			if s.Variant == TonalSpot {
				return tMaxC(s.TertiaryPalette, 0, 90, 1)
			}
			return tMaxC(s.TertiaryPalette, 0, 100, 1)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.TertiaryDim(),
				m.Tertiary(),
				5,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
	}
}

func (m MaterialColorSpec2025) OnTertiary() *DynamicColor {
	return &DynamicColor{
		Name: "on_tertiary",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Background: func(s DynamicScheme) *DynamicColor {
			return ternary(s.Platform == Phone, m.Tertiary(), m.TertiaryDim())
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return ternary(s.Platform == Phone, GetCurve(6), GetCurve(7))
		},
	}
}

func (m MaterialColorSpec2025) TertiaryContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "tertiary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Watch {
				if s.Variant == TonalSpot {
					return tMaxC(s.TertiaryPalette, 0, 90, 1)
				}
				return tMaxC(s.TertiaryPalette, 0, 100, 1)
			}

			switch s.Variant {
			case Neutral:
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 93, 1)
				}
				return tMaxC(s.TertiaryPalette, 0, 96, 1)

			case TonalSpot:
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 93, 1)
				}
				return tMaxC(s.TertiaryPalette, 0, 100, 1)

			case Expressive:
				upper := 100.0
				if s.TertiaryPalette.IsCyan() {
					upper = 88
				} else if s.IsDark {
					upper = 93
				}
				return tMaxC(s.TertiaryPalette, 75, upper, 1)

			case Vibrant:
				if s.IsDark {
					return tMaxC(s.TertiaryPalette, 0, 93, 1)
				}
				return tMaxC(s.TertiaryPalette, 72, 100, 1)
			}

			return tMaxC(s.TertiaryPalette, 0, 100, 1) // fallback
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{} // undefined
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Watch {
				return NewToneDeltaPair(
					m.TertiaryContainer(),
					m.TertiaryDim(),
					10,
					ToneDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) OnTertiaryContainer() *DynamicColor {
	return &DynamicColor{
		Name: "on_tertiary_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette {
			return s.TertiaryPalette
		},
		Background: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return ternary(s.Platform == Phone, GetCurve(6), GetCurve(7))
		},
	}
}

func (m MaterialColorSpec2025) TertiaryFixed() *DynamicColor {
	return &DynamicColor{
		Name:    "tertiary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Tone: func(s DynamicScheme) float64 {
			temp := s
			temp.IsDark = false
			return m.TertiaryContainer().Tone(temp)
		},
		IsBackground: true,
	}
}

func (m MaterialColorSpec2025) OnTertiaryFixed() *DynamicColor {
	return &DynamicColor{
		Name:    "on_tertiary_fixed",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryFixedDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) OnTertiaryFixedVariant() *DynamicColor {
	return &DynamicColor{
		Name:    "on_tertiary_fixed_variant",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.TertiaryPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.TertiaryFixedDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
	}
}

func (m MaterialColorSpec2025) Error() *DynamicColor {
	return &DynamicColor{
		Name:    "error",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Phone {
				if s.IsDark {
					return tMinC(s.ErrorPalette, 0, 98)
				}
				return tMaxC(s.ErrorPalette, 0, 100, 1)
			}
			return tMinC(s.ErrorPalette, 0, 100)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Phone {
				return NewToneDeltaPair(
					m.ErrorContainer(),
					m.Error(),
					5,
					ToneRelativeLighter,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) ErrorDim() *DynamicColor {
	return &DynamicColor{
		Name:    "error_dim",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s DynamicScheme) float64 {
			return tMinC(s.ErrorPalette, 0, 100)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			return m.SurfaceContainerHigh()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			return GetCurve(4.5)
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			return NewToneDeltaPair(
				m.ErrorDim(),
				m.Error(),
				5,
				ToneDarker,
				true,
				ConstraintFarther,
			)
		},
	}
}

func (m MaterialColorSpec2025) OnError() *DynamicColor {
	return &DynamicColor{
		Name:    "on_error",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.Error()
			}
			return m.ErrorDim()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(6)
			}
			return GetCurve(7)
		},
	}
}

func (m MaterialColorSpec2025) ErrorContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "error_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Tone: func(s DynamicScheme) float64 {
			if s.Platform == Watch {
				return 30
			}
			if s.IsDark {
				return tMinC(s.ErrorPalette, 30, 93)
			}
			return tMaxC(s.ErrorPalette, 0, 90, 1)
		},
		IsBackground: true,
		Background: func(s DynamicScheme) *DynamicColor {
			if s.Platform == Phone {
				return m.HighestSurface(s)
			}
			return &DynamicColor{} // or `return nil` if *DynamicColor is a pointer type
		},
		ToneDeltaPair: func(s DynamicScheme) *ToneDeltaPair {
			if s.Platform == Watch {
				return NewToneDeltaPair(
					m.ErrorContainer(),
					m.ErrorDim(),
					10,
					ToneDarker,
					true,
					ConstraintFarther,
				)
			}
			return nil
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone && s.ContrastLevel > 0 {
				return GetCurve(1.5)
			}
			return nil
		},
	}
}

func (m MaterialColorSpec2025) OnErrorContainer() *DynamicColor {
	return &DynamicColor{
		Name:    "on_error_container",
		Palette: func(s DynamicScheme) palettes.TonalPalette { return s.ErrorPalette },
		Background: func(s DynamicScheme) *DynamicColor {
			return m.ErrorContainer()
		},
		ContrastCurve: func(s DynamicScheme) *ContrastCurve {
			if s.Platform == Phone {
				return GetCurve(4.5)
			}
			return GetCurve(7)
		},
	}
}
