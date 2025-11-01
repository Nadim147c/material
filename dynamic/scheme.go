package dynamic

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/num"
	"github.com/Nadim147c/material/v2/palettes"
)

// Scheme represents a dynamic color scheme constructed from a set of values
// representing the current UI state (dark/light theme, theme style, etc.),
// and provides a set of TonalPalettes that create colors fitting the theme
// style. Used by Color to resolve into actual colors.
type Scheme struct {
	// SourceColorHct is the source color of the theme as an HCT color
	SourceColorHct color.Hct
	// Variant is the style variant of the theme (e.g., monochrome, tonal spot,
	// etc.)
	Variant Variant
	// Dark indicates whether the scheme is in dark mode (true) or light mode
	// (false)
	Dark bool
	// Platform specifies the platform on which this scheme is intended to be
	// used
	Platform Platform
	// Version specifies the version of the Material Design spec (2021 or 2025)
	Version Version
	// ContrastLevel represents the contrast level from -1 to 1, where:
	// -1 = minimum contrast, 0 = standard contrast, 1 = maximum contrast
	ContrastLevel float64
	// PrimaryPalette produces colors for primary UI elements. Usually colorful.
	PrimaryPalette palettes.TonalPalette
	// SecondaryPalette produces colors for secondary UI elements. Usually less
	// colorful.
	SecondaryPalette palettes.TonalPalette
	// TertiaryPalette produces colors for tertiary UI elements. Usually a
	// different hue from primary and colorful.
	TertiaryPalette palettes.TonalPalette
	// NeutralPalette produces neutral colors for backgrounds and surfaces.
	// Usually not colorful at all.
	NeutralPalette palettes.TonalPalette
	// NeutralVariantPalette produces neutral variant colors for backgrounds and
	// surfaces. Usually not colorful, but slightly more colorful than Neutral
	// palette.
	NeutralVariantPalette palettes.TonalPalette
	// ErrorPalette produces colors for error states. Usually reddish and
	// colorful.
	ErrorPalette palettes.TonalPalette
	// MaterialColor provides the material color specification implementation
	// for the given version (2021 or 2025)
	MaterialColor MaterialColorSpec
}

// NewDynamicScheme creates a new dynamic color scheme.
//
// Parameters:
//   - sourceColorHct: The source color of the theme as an HCT color
//   - variant: The variant, or style, of the theme
//   - contrastLevel: Value from -1 to 1. -1 represents minimum contrast, 0
//     represents standard (the design as specified), and 1 represents maximum
//     contrast
//   - dark: Whether the scheme is in dark mode or light mode
//   - platform: The platform on which this scheme is intended to be used
//   - version: The version of the design spec that this scheme is based on
//
// - optPalettes (optional): Up to six *palettes.TonalPalette arguments in this
// order:
//  1. primaryPalette
//  2. secondaryPalette
//  3. tertiaryPalette
//  4. neutralPalette
//  5. neutralVariantPalette
//  6. errorPalette
//
// Missing palettes are generated automatically. If errorPalette is not
// provided,
// a default reddish palette (hue 25.0, chroma 84.0) is used.
//
// The function automatically selects the appropriate palette delegate and
// material
// color specification based on the version (2021 or 2025).
func NewDynamicScheme(
	sourceColorHct color.Hct,
	variant Variant,
	contrastLevel float64,
	dark bool,
	platform Platform,
	version Version,
	optPalettes ...*palettes.TonalPalette,
) *Scheme {
	selectPalette := func(i int) *palettes.TonalPalette {
		if i < len(optPalettes) {
			return optPalettes[i]
		}
		return nil
	}

	primaryPalette := selectPalette(0)
	secondaryPalette := selectPalette(1)
	tertiaryPalette := selectPalette(2)
	neutralPalette := selectPalette(3)
	neutralVariantPalette := selectPalette(4)
	errorPalette := selectPalette(5)

	var palettesDelegate SchemePalettesDelegate = &schemePalettesDelegateImpl2021{}
	var colorSpec MaterialColorSpec = &MaterialSpec2021{}
	if version == Version2025 {
		palettesDelegate = &schemePalettesDelegateImpl2025{}
		colorSpec = &MaterialSpec2025{}
	}

	if primaryPalette == nil {
		primaryPalette = palettesDelegate.GetPrimaryPalette(
			variant,
			sourceColorHct,
			dark,
			platform,
			contrastLevel,
		)
	}
	if secondaryPalette == nil {
		secondaryPalette = palettesDelegate.GetSecondaryPalette(
			variant,
			sourceColorHct,
			dark,
			platform,
			contrastLevel,
		)
	}
	if tertiaryPalette == nil {
		tertiaryPalette = palettesDelegate.GetTertiaryPalette(
			variant,
			sourceColorHct,
			dark,
			platform,
			contrastLevel,
		)
	}
	if neutralPalette == nil {
		neutralPalette = palettesDelegate.GetNeutralPalette(
			variant,
			sourceColorHct,
			dark,
			platform,
			contrastLevel,
		)
	}
	if neutralVariantPalette == nil {
		neutralVariantPalette = palettesDelegate.GetNeutralVariantPalette(
			variant,
			sourceColorHct,
			dark,
			platform,
			contrastLevel,
		)
	}
	if errorPalette == nil {
		errorPalette = palettes.FromHueAndChroma(25.0, 84.0)
	}

	return &Scheme{
		SourceColorHct:        sourceColorHct,
		Variant:               variant,
		Dark:                  dark,
		Platform:              platform,
		Version:               version,
		ContrastLevel:         contrastLevel,
		PrimaryPalette:        *primaryPalette,
		SecondaryPalette:      *secondaryPalette,
		TertiaryPalette:       *tertiaryPalette,
		NeutralPalette:        *neutralPalette,
		NeutralVariantPalette: *neutralVariantPalette,
		ErrorPalette:          *errorPalette,
		MaterialColor:         colorSpec,
	}
}

// GetPiecewiseHue returns a new hue based on a piece wise function and the
// input color's hue.
func GetPiecewiseHue(
	sourceColorHct color.Hct,
	hueBreakpoints []float64,
	hues []float64,
) float64 {
	size := min(len(hues), len(hueBreakpoints)-1)
	sourceHue := sourceColorHct.Hue
	for i := range size {
		if sourceHue >= hueBreakpoints[i] && sourceHue < hueBreakpoints[i+1] {
			return num.NormalizeDegree(hues[i])
		}
	}
	// No match found, return the source hue.
	return sourceHue
}

// GetRotatedHue returns a shifted hue based on a piece wise function and the
// input hue.
func GetRotatedHue(
	sourceColorHct color.Hct,
	hueBreakpoints []float64,
	rotations []float64,
) float64 {
	rotation := GetPiecewiseHue(sourceColorHct, hueBreakpoints, rotations)
	if min(len(hueBreakpoints)-1, len(rotations)) <= 0 {
		// No valid range; apply no rotation.
		rotation = 0
	}

	return num.NormalizeDegree(sourceColorHct.Hue + rotation)
}

// SourceColorARGB returns ARGB version of source color.
func (d Scheme) SourceColorARGB() color.ARGB {
	return d.SourceColorHct.ToARGB()
}

// ToColorMap creates a map of color name as key and *Color as value.
func (d Scheme) ToColorMap() map[string]*Color {
	return map[string]*Color{
		"primary_palette_key_color":         d.MaterialColor.PrimaryPaletteKeyColor(),
		"secondary_palette_key_color":       d.MaterialColor.SecondaryPaletteKeyColor(),
		"tertiary_palette_key_color":        d.MaterialColor.TertiaryPaletteKeyColor(),
		"neutral_palette_key_color":         d.MaterialColor.NeutralPaletteKeyColor(),
		"neutral_variant_palette_key_color": d.MaterialColor.NeutralVariantPaletteKeyColor(),
		"background":                        d.MaterialColor.Background(),
		"on_background":                     d.MaterialColor.OnBackground(),
		"surface":                           d.MaterialColor.Surface(),
		"surface_dim":                       d.MaterialColor.SurfaceDim(),
		"surface_bright":                    d.MaterialColor.SurfaceBright(),
		"surface_container_lowest":          d.MaterialColor.SurfaceContainerLowest(),
		"surface_container_low":             d.MaterialColor.SurfaceContainerLow(),
		"surface_container":                 d.MaterialColor.SurfaceContainer(),
		"surface_container_high":            d.MaterialColor.SurfaceContainerHigh(),
		"surface_container_highest":         d.MaterialColor.SurfaceContainerHighest(),
		"on_surface":                        d.MaterialColor.OnSurface(),
		"surface_variant":                   d.MaterialColor.SurfaceVariant(),
		"on_surface_variant":                d.MaterialColor.OnSurfaceVariant(),
		"inverse_surface":                   d.MaterialColor.InverseSurface(),
		"inverse_on_surface":                d.MaterialColor.InverseOnSurface(),
		"outline":                           d.MaterialColor.Outline(),
		"outline_variant":                   d.MaterialColor.OutlineVariant(),
		"shadow":                            d.MaterialColor.Shadow(),
		"scrim":                             d.MaterialColor.Scrim(),
		"surface_tint":                      d.MaterialColor.SurfaceTint(),
		"primary":                           d.MaterialColor.Primary(),
		"on_primary":                        d.MaterialColor.OnPrimary(),
		"primary_container":                 d.MaterialColor.PrimaryContainer(),
		"primary_dim":                       d.MaterialColor.PrimaryDim(),
		"on_primary_container":              d.MaterialColor.OnPrimaryContainer(),
		"inverse_primary":                   d.MaterialColor.InversePrimary(),
		"secondary":                         d.MaterialColor.Secondary(),
		"on_secondary":                      d.MaterialColor.OnSecondary(),
		"secondary_container":               d.MaterialColor.SecondaryContainer(),
		"secondary_dim":                     d.MaterialColor.SecondaryDim(),
		"on_secondary_container":            d.MaterialColor.OnSecondaryContainer(),
		"tertiary":                          d.MaterialColor.Tertiary(),
		"on_tertiary":                       d.MaterialColor.OnTertiary(),
		"tertiary_container":                d.MaterialColor.TertiaryContainer(),
		"tertiary_dim":                      d.MaterialColor.TertiaryDim(),
		"on_tertiary_container":             d.MaterialColor.OnTertiaryContainer(),
		"error":                             d.MaterialColor.Error(),
		"on_error":                          d.MaterialColor.OnError(),
		"error_container":                   d.MaterialColor.ErrorContainer(),
		"error_dim":                         d.MaterialColor.ErrorDim(),
		"on_error_container":                d.MaterialColor.OnErrorContainer(),
		"primary_fixed":                     d.MaterialColor.PrimaryFixed(),
		"primary_fixed_dim":                 d.MaterialColor.PrimaryFixedDim(),
		"on_primary_fixed":                  d.MaterialColor.OnPrimaryFixed(),
		"on_primary_fixed_variant":          d.MaterialColor.OnPrimaryFixedVariant(),
		"secondary_fixed":                   d.MaterialColor.SecondaryFixed(),
		"secondary_fixed_dim":               d.MaterialColor.SecondaryFixedDim(),
		"on_secondary_fixed":                d.MaterialColor.OnSecondaryFixed(),
		"on_secondary_fixed_variant":        d.MaterialColor.OnSecondaryFixedVariant(),
		"tertiary_fixed":                    d.MaterialColor.TertiaryFixed(),
		"tertiary_fixed_dim":                d.MaterialColor.TertiaryFixedDim(),
		"on_tertiary_fixed":                 d.MaterialColor.OnTertiaryFixed(),
		"on_tertiary_fixed_variant":         d.MaterialColor.OnTertiaryFixedVariant(),
	}
}
