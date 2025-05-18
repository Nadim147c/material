package dynamic

import (
	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/goyou/num"
	"github.com/Nadim147c/goyou/palettes"
)

type DynamicScheme struct {
	SourceColorHct color.Hct
	Variant        Variant
	IsDark         bool
	ContrastLevel  float64

	PrimaryPalette        palettes.TonalPalette
	SecondaryPalette      palettes.TonalPalette
	TertiaryPalette       palettes.TonalPalette
	NeutralPalette        palettes.TonalPalette
	NeutralVariantPalette palettes.TonalPalette
	ErrorPalette          palettes.TonalPalette
}

func NewDynamicScheme(
	sourceColorHct color.Hct,
	variant Variant,
	contrastLevel float64,
	isDark bool,
	primaryPalette palettes.TonalPalette,
	secondaryPalette palettes.TonalPalette,
	tertiaryPalette palettes.TonalPalette,
	neutralPalette palettes.TonalPalette,
	neutralVariantPalette palettes.TonalPalette,
	errorPalette *palettes.TonalPalette,
) DynamicScheme {
	finalErrorPalette := palettes.FromHueAndChroma(25.0, 84.0)
	if errorPalette != nil {
		finalErrorPalette = errorPalette
	}
	return DynamicScheme{
		SourceColorHct:        sourceColorHct,
		Variant:               variant,
		IsDark:                isDark,
		ContrastLevel:         contrastLevel,
		PrimaryPalette:        primaryPalette,
		SecondaryPalette:      secondaryPalette,
		TertiaryPalette:       tertiaryPalette,
		NeutralPalette:        neutralPalette,
		NeutralVariantPalette: neutralVariantPalette,
		ErrorPalette:          *finalErrorPalette,
	}
}

func (d DynamicScheme) GetRotatedHue(sourceColor color.Hct, hues, rotations []float64) float64 {
	sourceHue := sourceColor.Hue
	if len(rotations) == 1 {
		return num.NormalizeDegree(sourceHue + rotations[0])
	}
	for i := 0; i <= len(hues)-2; i++ {
		if hues[i] < sourceHue && sourceHue < hues[i+1] {
			return num.NormalizeDegree(sourceHue + rotations[i])
		}
	}
	return sourceHue
}

func (d DynamicScheme) SourceColorArgb() color.ARGB {
	return d.SourceColorHct.ToARGB()
}

// Below are the color accessor methods.
func (d DynamicScheme) PrimaryPaletteKeyColor() color.ARGB {
	return DynamicSchemeProvider.PrimaryPaletteKeyColor().GetArgb(d)
}

func (d DynamicScheme) SecondaryPaletteKeyColor() color.ARGB {
	return DynamicSchemeProvider.SecondaryPaletteKeyColor().GetArgb(d)
}

func (d DynamicScheme) TertiaryPaletteKeyColor() color.ARGB {
	return DynamicSchemeProvider.TertiaryPaletteKeyColor().GetArgb(d)
}

func (d DynamicScheme) NeutralPaletteKeyColor() color.ARGB {
	return DynamicSchemeProvider.NeutralPaletteKeyColor().GetArgb(d)
}

func (d DynamicScheme) NeutralVariantPaletteKeyColor() color.ARGB {
	return DynamicSchemeProvider.NeutralVariantPaletteKeyColor().GetArgb(d)
}

func (d DynamicScheme) Background() color.ARGB {
	return DynamicSchemeProvider.Background().GetArgb(d)
}

func (d DynamicScheme) OnBackground() color.ARGB {
	return DynamicSchemeProvider.OnBackground().GetArgb(d)
}

func (d DynamicScheme) Surface() color.ARGB {
	return DynamicSchemeProvider.Surface().GetArgb(d)
}

func (d DynamicScheme) SurfaceDim() color.ARGB {
	return DynamicSchemeProvider.SurfaceDim().GetArgb(d)
}

func (d DynamicScheme) SurfaceBright() color.ARGB {
	return DynamicSchemeProvider.SurfaceBright().GetArgb(d)
}

func (d DynamicScheme) SurfaceContainerLowest() color.ARGB {
	return DynamicSchemeProvider.SurfaceContainerLowest().GetArgb(d)
}

func (d DynamicScheme) SurfaceContainerLow() color.ARGB {
	return DynamicSchemeProvider.SurfaceContainerLow().GetArgb(d)
}

func (d DynamicScheme) SurfaceContainer() color.ARGB {
	return DynamicSchemeProvider.SurfaceContainer().GetArgb(d)
}

func (d DynamicScheme) SurfaceContainerHigh() color.ARGB {
	return DynamicSchemeProvider.SurfaceContainerHigh().GetArgb(d)
}

func (d DynamicScheme) SurfaceContainerHighest() color.ARGB {
	return DynamicSchemeProvider.SurfaceContainerHighest().GetArgb(d)
}

func (d DynamicScheme) OnSurface() color.ARGB {
	return DynamicSchemeProvider.OnSurface().GetArgb(d)
}

func (d DynamicScheme) SurfaceVariant() color.ARGB {
	return DynamicSchemeProvider.SurfaceVariant().GetArgb(d)
}

func (d DynamicScheme) OnSurfaceVariant() color.ARGB {
	return DynamicSchemeProvider.OnSurfaceVariant().GetArgb(d)
}

func (d DynamicScheme) InverseSurface() color.ARGB {
	return DynamicSchemeProvider.InverseSurface().GetArgb(d)
}

func (d DynamicScheme) InverseOnSurface() color.ARGB {
	return DynamicSchemeProvider.InverseOnSurface().GetArgb(d)
}

func (d DynamicScheme) Outline() color.ARGB {
	return DynamicSchemeProvider.Outline().GetArgb(d)
}

func (d DynamicScheme) OutlineVariant() color.ARGB {
	return DynamicSchemeProvider.OutlineVariant().GetArgb(d)
}

func (d DynamicScheme) Shadow() color.ARGB {
	return DynamicSchemeProvider.Shadow().GetArgb(d)
}

func (d DynamicScheme) Scrim() color.ARGB {
	return DynamicSchemeProvider.Scrim().GetArgb(d)
}

func (d DynamicScheme) SurfaceTint() color.ARGB {
	return DynamicSchemeProvider.SurfaceTint().GetArgb(d)
}

func (d DynamicScheme) Primary() color.ARGB {
	return DynamicSchemeProvider.Primary().GetArgb(d)
}

func (d DynamicScheme) OnPrimary() color.ARGB {
	return DynamicSchemeProvider.OnPrimary().GetArgb(d)
}

func (d DynamicScheme) PrimaryContainer() color.ARGB {
	return DynamicSchemeProvider.PrimaryContainer().GetArgb(d)
}

func (d DynamicScheme) OnPrimaryContainer() color.ARGB {
	return DynamicSchemeProvider.OnPrimaryContainer().GetArgb(d)
}

func (d DynamicScheme) InversePrimary() color.ARGB {
	return DynamicSchemeProvider.InversePrimary().GetArgb(d)
}

func (d DynamicScheme) Secondary() color.ARGB {
	return DynamicSchemeProvider.Secondary().GetArgb(d)
}

func (d DynamicScheme) OnSecondary() color.ARGB {
	return DynamicSchemeProvider.OnSecondary().GetArgb(d)
}

func (d DynamicScheme) SecondaryContainer() color.ARGB {
	return DynamicSchemeProvider.SecondaryContainer().GetArgb(d)
}

func (d DynamicScheme) OnSecondaryContainer() color.ARGB {
	return DynamicSchemeProvider.OnSecondaryContainer().GetArgb(d)
}

func (d DynamicScheme) Tertiary() color.ARGB {
	return DynamicSchemeProvider.Tertiary().GetArgb(d)
}

func (d DynamicScheme) OnTertiary() color.ARGB {
	return DynamicSchemeProvider.OnTertiary().GetArgb(d)
}

func (d DynamicScheme) TertiaryContainer() color.ARGB {
	return DynamicSchemeProvider.TertiaryContainer().GetArgb(d)
}

func (d DynamicScheme) OnTertiaryContainer() color.ARGB {
	return DynamicSchemeProvider.OnTertiaryContainer().GetArgb(d)
}

func (d DynamicScheme) Error() color.ARGB {
	return DynamicSchemeProvider.Error().GetArgb(d)
}

func (d DynamicScheme) OnError() color.ARGB {
	return DynamicSchemeProvider.OnError().GetArgb(d)
}

func (d DynamicScheme) ErrorContainer() color.ARGB {
	return DynamicSchemeProvider.ErrorContainer().GetArgb(d)
}

func (d DynamicScheme) OnErrorContainer() color.ARGB {
	return DynamicSchemeProvider.OnErrorContainer().GetArgb(d)
}

func (d DynamicScheme) PrimaryFixed() color.ARGB {
	return DynamicSchemeProvider.PrimaryFixed().GetArgb(d)
}

func (d DynamicScheme) PrimaryFixedDim() color.ARGB {
	return DynamicSchemeProvider.PrimaryFixedDim().GetArgb(d)
}

func (d DynamicScheme) OnPrimaryFixed() color.ARGB {
	return DynamicSchemeProvider.OnPrimaryFixed().GetArgb(d)
}

func (d DynamicScheme) OnPrimaryFixedVariant() color.ARGB {
	return DynamicSchemeProvider.OnPrimaryFixedVariant().GetArgb(d)
}

func (d DynamicScheme) SecondaryFixed() color.ARGB {
	return DynamicSchemeProvider.SecondaryFixed().GetArgb(d)
}

func (d DynamicScheme) SecondaryFixedDim() color.ARGB {
	return DynamicSchemeProvider.SecondaryFixedDim().GetArgb(d)
}

func (d DynamicScheme) OnSecondaryFixed() color.ARGB {
	return DynamicSchemeProvider.OnSecondaryFixed().GetArgb(d)
}

func (d DynamicScheme) OnSecondaryFixedVariant() color.ARGB {
	return DynamicSchemeProvider.OnSecondaryFixedVariant().GetArgb(d)
}

func (d DynamicScheme) TertiaryFixed() color.ARGB {
	return DynamicSchemeProvider.TertiaryFixed().GetArgb(d)
}

func (d DynamicScheme) TertiaryFixedDim() color.ARGB {
	return DynamicSchemeProvider.TertiaryFixedDim().GetArgb(d)
}

func (d DynamicScheme) OnTertiaryFixed() color.ARGB {
	return DynamicSchemeProvider.OnTertiaryFixed().GetArgb(d)
}

func (d DynamicScheme) OnTertiaryFixedVariant() color.ARGB {
	return DynamicSchemeProvider.OnTertiaryFixedVariant().GetArgb(d)
}
