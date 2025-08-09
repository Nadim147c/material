package dynamic

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
	"github.com/Nadim147c/material/palettes"
)

type Version int

const (
	V2021 Version = 2021
	V2025 Version = 2025
)

type DynamicScheme struct {
	SourceColorHct color.Hct
	Variant        Variant
	IsDark         bool
	Platform       Platform
	Version        Version
	ContrastLevel  float64

	PrimaryPalette        palettes.TonalPalette
	SecondaryPalette      palettes.TonalPalette
	TertiaryPalette       palettes.TonalPalette
	NeutralPalette        palettes.TonalPalette
	NeutralVariantPalette palettes.TonalPalette
	ErrorPalette          palettes.TonalPalette
	MaterialColor         MaterialColorSpec
}

func NewDynamicScheme(
	sourceColorHct color.Hct,
	variant Variant,
	contrastLevel float64,
	isDark bool,
	platform Platform,
	version Version,
	primaryPalette *palettes.TonalPalette,
	secondaryPalette *palettes.TonalPalette,
	tertiaryPalette *palettes.TonalPalette,
	neutralPalette *palettes.TonalPalette,
	neutralVariantPalette *palettes.TonalPalette,
	errorPalette *palettes.TonalPalette,
) DynamicScheme {
	var palettesDelegate DynamicSchemePalettesDelegate = &DynamicSchemePalettesDelegateImpl2021{}
	var colorSpec MaterialColorSpec = &MaterialSpec2021{}
	if version == V2025 {
		palettesDelegate = &DynamicSchemePalettesDelegateImpl2025{}
		colorSpec = &MaterialSpec2025{}
	}
	if primaryPalette == nil {
		primaryPalette = palettesDelegate.GetPrimaryPalette(variant, sourceColorHct, isDark, Phone, contrastLevel)
	}
	if secondaryPalette == nil {
		secondaryPalette = palettesDelegate.GetSecondaryPalette(variant, sourceColorHct, isDark, Phone, contrastLevel)
	}
	if tertiaryPalette == nil {
		tertiaryPalette = palettesDelegate.GetTertiaryPalette(variant, sourceColorHct, isDark, Phone, contrastLevel)
	}
	if neutralPalette == nil {
		neutralPalette = palettesDelegate.GetNeutralPalette(variant, sourceColorHct, isDark, Phone, contrastLevel)
	}
	if neutralVariantPalette == nil {
		neutralVariantPalette = palettesDelegate.GetNeutralVariantPalette(variant, sourceColorHct, isDark, Phone, contrastLevel)
	}
	if errorPalette == nil {
		errorPalette = palettes.FromHueAndChroma(25.0, 84.0)
	}

	return DynamicScheme{
		SourceColorHct:        sourceColorHct,
		Variant:               variant,
		IsDark:                isDark,
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
func GetPiecewiseHue(sourceColorHct color.Hct, hueBreakpoints []float64, hues []float64) float64 {
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
func GetRotatedHue(sourceColorHct color.Hct, hueBreakpoints []float64, rotations []float64) float64 {
	rotation := GetPiecewiseHue(sourceColorHct, hueBreakpoints, rotations)
	if min(len(hueBreakpoints)-1, len(rotations)) <= 0 {
		// No valid range; apply no rotation.
		rotation = 0
	}

	return num.NormalizeDegree(sourceColorHct.Hue + rotation)
}

func (d DynamicScheme) SourceColorArgb() color.ARGB {
	return d.SourceColorHct.ToARGB()
}

func (d DynamicScheme) ToColorMap() map[string]*DynamicColor {
	return map[string]*DynamicColor{
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
