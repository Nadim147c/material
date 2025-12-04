package material

import (
	"context"
	"errors"
	"image"
	gocolor "image/color"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/Nadim147c/material/v2/blend"
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/material/v2/quantizer"
	"github.com/Nadim147c/material/v2/score"
)

// Colors is generated material you colors
type Colors struct {
	m      map[string]color.ARGB
	Scheme *dynamic.Scheme `json:"scheme,omitzero"`

	CustomColors map[string]CustomColor `json:"custom"`

	Background              color.ARGB `json:"background"`
	Error                   color.ARGB `json:"error"`
	ErrorContainer          color.ARGB `json:"error_container"`
	ErrorDim                color.ARGB `json:"error_dim"`
	InverseOnSurface        color.ARGB `json:"inverse_on_surface"`
	InversePrimary          color.ARGB `json:"inverse_primary"`
	InverseSurface          color.ARGB `json:"inverse_surface"`
	OnBackground            color.ARGB `json:"on_background"`
	OnError                 color.ARGB `json:"on_error"`
	OnErrorContainer        color.ARGB `json:"on_error_container"`
	OnPrimary               color.ARGB `json:"on_primary"`
	OnPrimaryContainer      color.ARGB `json:"on_primary_container"`
	OnPrimaryFixed          color.ARGB `json:"on_primary_fixed"`
	OnPrimaryFixedVariant   color.ARGB `json:"on_primary_fixed_variant"`
	OnSecondary             color.ARGB `json:"on_secondary"`
	OnSecondaryContainer    color.ARGB `json:"on_secondary_container"`
	OnSecondaryFixed        color.ARGB `json:"on_secondary_fixed"`
	OnSecondaryFixedVariant color.ARGB `json:"on_secondary_fixed_variant"`
	OnSurface               color.ARGB `json:"on_surface"`
	OnSurfaceVariant        color.ARGB `json:"on_surface_variant"`
	OnTertiary              color.ARGB `json:"on_tertiary"`
	OnTertiaryContainer     color.ARGB `json:"on_tertiary_container"`
	OnTertiaryFixed         color.ARGB `json:"on_tertiary_fixed"`
	OnTertiaryFixedVariant  color.ARGB `json:"on_tertiary_fixed_variant"`
	Outline                 color.ARGB `json:"outline"`
	OutlineVariant          color.ARGB `json:"outline_variant"`
	Primary                 color.ARGB `json:"primary"`
	PrimaryContainer        color.ARGB `json:"primary_container"`
	PrimaryDim              color.ARGB `json:"primary_dim"`
	PrimaryFixed            color.ARGB `json:"primary_fixed"`
	PrimaryFixedDim         color.ARGB `json:"primary_fixed_dim"`
	Scrim                   color.ARGB `json:"scrim"`
	Secondary               color.ARGB `json:"secondary"`
	SecondaryContainer      color.ARGB `json:"secondary_container"`
	SecondaryDim            color.ARGB `json:"secondary_dim"`
	SecondaryFixed          color.ARGB `json:"secondary_fixed"`
	SecondaryFixedDim       color.ARGB `json:"secondary_fixed_dim"`
	Shadow                  color.ARGB `json:"shadow"`
	Surface                 color.ARGB `json:"surface"`
	SurfaceBright           color.ARGB `json:"surface_bright"`
	SurfaceContainer        color.ARGB `json:"surface_container"`
	SurfaceContainerHigh    color.ARGB `json:"surface_container_high"`
	SurfaceContainerHighest color.ARGB `json:"surface_container_highest"`
	SurfaceContainerLow     color.ARGB `json:"surface_container_low"`
	SurfaceContainerLowest  color.ARGB `json:"surface_container_lowest"`
	SurfaceDim              color.ARGB `json:"surface_dim"`
	SurfaceTint             color.ARGB `json:"surface_tint"`
	SurfaceVariant          color.ARGB `json:"surface_variant"`
	Tertiary                color.ARGB `json:"tertiary"`
	TertiaryContainer       color.ARGB `json:"tertiary_container"`
	TertiaryDim             color.ARGB `json:"tertiary_dim"`
	TertiaryFixed           color.ARGB `json:"tertiary_fixed"`
	TertiaryFixedDim        color.ARGB `json:"tertiary_fixed_dim"`
}

// Map returns map with color name in snake case as name and color.ARGB as value
func (c *Colors) Map() map[string]color.ARGB {
	m := maps.Clone(c.m)
	for k, v := range c.CustomColors {
		key := strings.ToLower(k)
		m[key] = v.Color
		m["on_"+key] = v.Color
	}
	return m
}

// createColors converts a map of color names to Color pointers into a Colors
// struct
func createColors(
	scheme *dynamic.Scheme,
	custom map[string]CustomColorOption,
) *Colors {
	m := scheme.ToColorMap()
	primary := calc(scheme, m["primary"])

	customColors := make(map[string]CustomColor, len(custom))
	for name, opt := range custom {
		customColors[name] = createCustomColor(opt, scheme.Dark, primary)
	}

	return &Colors{
		Scheme:                  scheme,
		Primary:                 primary,
		CustomColors:            customColors,
		Background:              calc(scheme, m["background"]),
		Error:                   calc(scheme, m["error"]),
		ErrorContainer:          calc(scheme, m["error_container"]),
		ErrorDim:                calc(scheme, m["error_dim"]),
		InverseOnSurface:        calc(scheme, m["inverse_on_surface"]),
		InversePrimary:          calc(scheme, m["inverse_primary"]),
		InverseSurface:          calc(scheme, m["inverse_surface"]),
		OnBackground:            calc(scheme, m["on_background"]),
		OnError:                 calc(scheme, m["on_error"]),
		OnErrorContainer:        calc(scheme, m["on_error_container"]),
		OnPrimary:               calc(scheme, m["on_primary"]),
		OnPrimaryContainer:      calc(scheme, m["on_primary_container"]),
		OnPrimaryFixed:          calc(scheme, m["on_primary_fixed"]),
		OnPrimaryFixedVariant:   calc(scheme, m["on_primary_fixed_variant"]),
		OnSecondary:             calc(scheme, m["on_secondary"]),
		OnSecondaryContainer:    calc(scheme, m["on_secondary_container"]),
		OnSecondaryFixed:        calc(scheme, m["on_secondary_fixed"]),
		OnSecondaryFixedVariant: calc(scheme, m["on_secondary_fixed_variant"]),
		OnSurface:               calc(scheme, m["on_surface"]),
		OnSurfaceVariant:        calc(scheme, m["on_surface_variant"]),
		OnTertiary:              calc(scheme, m["on_tertiary"]),
		OnTertiaryContainer:     calc(scheme, m["on_tertiary_container"]),
		OnTertiaryFixed:         calc(scheme, m["on_tertiary_fixed"]),
		OnTertiaryFixedVariant:  calc(scheme, m["on_tertiary_fixed_variant"]),
		Outline:                 calc(scheme, m["outline"]),
		OutlineVariant:          calc(scheme, m["outline_variant"]),
		PrimaryContainer:        calc(scheme, m["primary_container"]),
		PrimaryDim:              calc(scheme, m["primary_dim"]),
		PrimaryFixed:            calc(scheme, m["primary_fixed"]),
		PrimaryFixedDim:         calc(scheme, m["primary_fixed_dim"]),
		Scrim:                   calc(scheme, m["scrim"]),
		Secondary:               calc(scheme, m["secondary"]),
		SecondaryContainer:      calc(scheme, m["secondary_container"]),
		SecondaryDim:            calc(scheme, m["secondary_dim"]),
		SecondaryFixed:          calc(scheme, m["secondary_fixed"]),
		SecondaryFixedDim:       calc(scheme, m["secondary_fixed_dim"]),
		Shadow:                  calc(scheme, m["shadow"]),
		Surface:                 calc(scheme, m["surface"]),
		SurfaceBright:           calc(scheme, m["surface_bright"]),
		SurfaceContainer:        calc(scheme, m["surface_container"]),
		SurfaceContainerHigh:    calc(scheme, m["surface_container_high"]),
		SurfaceContainerHighest: calc(scheme, m["surface_container_highest"]),
		SurfaceContainerLow:     calc(scheme, m["surface_container_low"]),
		SurfaceContainerLowest:  calc(scheme, m["surface_container_lowest"]),
		SurfaceDim:              calc(scheme, m["surface_dim"]),
		SurfaceTint:             calc(scheme, m["surface_tint"]),
		SurfaceVariant:          calc(scheme, m["surface_variant"]),
		Tertiary:                calc(scheme, m["tertiary"]),
		TertiaryContainer:       calc(scheme, m["tertiary_container"]),
		TertiaryDim:             calc(scheme, m["tertiary_dim"]),
		TertiaryFixed:           calc(scheme, m["tertiary_fixed"]),
		TertiaryFixedDim:        calc(scheme, m["tertiary_fixed_dim"]),
	}
}

// CustomColor is the custom colors generated from user defined colors
type CustomColor struct {
	Color            color.ARGB `json:"color"`
	OnColor          color.ARGB `json:"on_color"`
	ColorContainer   color.ARGB `json:"color_container"`
	OnColorContainer color.ARGB `json:"on_color_container"`
}

func tonedHctToARGB(hct color.Hct, tone float64) color.ARGB {
	hct.Tone = tone
	return hct.ToARGB()
}

func createCustomColor(
	option CustomColorOption,
	dark bool,
	to color.ARGB,
) CustomColor {
	var hct color.Hct

	if option.Blend {
		hct = blend.HctHueDirect(option.Color, to, option.Ratio)
	} else {
		hct = option.Color.ToHct()
	}

	if dark {
		return CustomColor{
			Color:            tonedHctToARGB(hct, 40),
			OnColor:          tonedHctToARGB(hct, 100),
			ColorContainer:   tonedHctToARGB(hct, 90),
			OnColorContainer: tonedHctToARGB(hct, 10),
		}
	}
	return CustomColor{
		Color:            tonedHctToARGB(hct, 80),
		OnColor:          tonedHctToARGB(hct, 20),
		ColorContainer:   tonedHctToARGB(hct, 30),
		OnColorContainer: tonedHctToARGB(hct, 90),
	}
}

// calc converts a *Color pointer to color.ARGB Returns 0 if the pointer is nil
func calc(s *dynamic.Scheme, c *dynamic.Color) color.ARGB {
	if c == nil {
		return 0
	}
	return c.GetArgb(s)
}

// Source is a function that returns source colors for material you
type Source func() ([]color.ARGB, error)

// FromImage returns Source colors from image.Image interface
func FromImage(img image.Image) Source {
	return func() ([]color.ARGB, error) {
		bounds := img.Bounds()
		pixels := make([]color.ARGB, 0, bounds.Dx()*bounds.Dy())
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := img.At(x, y)
				argb := color.ARGBFromInterface(c)
				pixels = append(pixels, argb)
			}
		}
		return pixels, nil
	}
}

// FromColor returns a Source from a single color.Color interface
func FromColor(c gocolor.Color) Source {
	argb := color.ARGBFromInterface(c)
	return func() ([]color.ARGB, error) {
		return []color.ARGB{argb}, nil
	}
}

// FromColors returns a Source from a slice of color.Color interfaces
func FromColors(s []gocolor.Color) Source {
	colors := make([]color.ARGB, len(s))
	for i, ci := range s {
		colors[i] = color.ARGBFromInterface(ci)
	}
	return func() ([]color.ARGB, error) {
		return colors, nil
	}
}

// FromBytes returns Source colors from a byte slice.
//
// WARNING: DO NOT pass image encoded file buffer. byte slice should be a
// sequence of r, g, b bytes (3 bytes per color).
func FromBytes(b []byte) Source {
	size := len(b)
	colors := make([]color.ARGB, 0, size/3)
	for i := 0; i < size-2; i += 3 {
		colors = append(colors, color.ARGBFromRGB(b[i], b[i+1], b[i+2]))
	}
	return func() ([]color.ARGB, error) {
		return colors, nil
	}
}

// FromARGB returns Source colors from a slice of ARGB values
func FromARGB(argbs []color.ARGB) Source {
	return func() ([]color.ARGB, error) {
		return argbs, nil
	}
}

// FromReader returns Source colors from an io.Reader containing RGB bytes.
// Reads until EOF and extracts RGB triplets from the stream.
//
// WARNING: Do NOT pass image encoded file readers (e.g., PNG, JPEG). Use
// FromImage with image.Decode for encoded image files.
func FromReader(r io.Reader) Source {
	return func() ([]color.ARGB, error) {
		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		size := len(b)
		colors := make([]color.ARGB, 0, size/3)
		for i := 0; i < size-2; i += 3 {
			colors = append(colors, color.ARGBFromRGB(b[i], b[i+1], b[i+2]))
		}
		return colors, nil
	}
}

// FromHex returns a Source from a hex color string (e.g., "#FF5733" or
// "FF5733")
func FromHex(hex string) Source {
	argb, err := color.ARGBFromHex(hex)
	return func() ([]color.ARGB, error) {
		return []color.ARGB{argb}, err
	}
}

// FromHexes returns a Source from multiple hex color strings
func FromHexes(hexes []string) Source {
	return func() ([]color.ARGB, error) {
		colors := make([]color.ARGB, 0, len(hexes))
		for _, hex := range hexes {
			argb, err := color.ARGBFromHex(hex)
			if err != nil {
				return colors, err
			}
			colors = append(colors, argb)
		}
		return colors, nil
	}
}

// Combine returns a Source that merges multiple Sources into one
func Combine(sources ...Source) Source {
	return func() ([]color.ARGB, error) {
		var result []color.ARGB
		for _, source := range sources {
			colors, err := source()
			if err != nil {
				return result, err
			}
			result = append(result, colors...)
		}
		return result, nil
	}
}

// Filter returns a Source that filters colors using a predicate function
func Filter(source Source, predicate func(color.ARGB) bool) Source {
	return func() ([]color.ARGB, error) {
		colors, err := source()
		if err != nil {
			return colors, err
		}
		filtered := make([]color.ARGB, 0, len(colors))
		for _, c := range colors {
			if predicate(c) {
				filtered = append(filtered, c)
			}
		}
		return filtered, nil
	}
}

// CustomColorOption is used define custom color
type CustomColorOption struct {
	Blend bool
	Ratio float64
	Color color.ARGB
}

// Settings is the dynamic schema configuration
type Settings struct {
	Context  context.Context  `json:"-"` // context shouldn't be encoded
	Contrast float64          `json:"contrast"`
	Dark     bool             `json:"dark"`
	Platform dynamic.Platform `json:"platform"`
	Variant  dynamic.Variant  `json:"variant"`
	Version  dynamic.Version  `json:"version"`

	Custom map[string]CustomColorOption `json:"-"`
}

// Option is a func modifes the dynamic scheme settings
type Option func(s *Settings)

// WithContext returns an Option that set the context
func WithContext(ctx context.Context) Option {
	return func(s *Settings) { s.Context = ctx }
}

// WithContrast returns an Option that sets the contrast level
func WithContrast(c float64) Option {
	return func(s *Settings) { s.Contrast = c }
}

// WithDark returns an Option that sets the dark mode flag
func WithDark(d bool) Option {
	return func(s *Settings) { s.Dark = d }
}

// WithPlatform returns an Option that sets the platform
func WithPlatform(p dynamic.Platform) Option {
	return func(s *Settings) { s.Platform = p }
}

// WithVariant returns an Option that sets the variant
func WithVariant(v dynamic.Variant) Option {
	return func(s *Settings) { s.Variant = v }
}

// WithVersion returns an Option that sets the version
func WithVersion(v dynamic.Version) Option {
	return func(s *Settings) { s.Version = v }
}

// WithCustomColor returns an Option that adds a custom color.
func WithCustomColor(name string, c gocolor.Color) Option {
	return func(o *Settings) {
		if o.Custom == nil {
			o.Custom = map[string]CustomColorOption{}
		}
		o.Custom[name] = CustomColorOption{
			Color: color.ARGBFromInterface(c),
		}
	}
}

// WithCustomColorBlend returns an Option that adds a custom color which will be
// blended with primary color by geven ratio. Ratio range [0, 1].
func WithCustomColorBlend(name string, c gocolor.Color, ratio float64) Option {
	return func(o *Settings) {
		if o.Custom == nil {
			o.Custom = map[string]CustomColorOption{}
		}
		o.Custom[name] = CustomColorOption{
			Blend: true,
			Ratio: ratio,
			Color: color.ARGBFromInterface(c),
		}
	}
}

// WithSettings settings all values of settings
func WithSettings(s Settings) Option {
	return func(o *Settings) { *o = s }
}

var errNoColorFound = errors.New("no source colors")

// Generate generates material you colors
func Generate(src Source, options ...Option) (*Colors, error) {
	colors, err := src()
	if err != nil {
		return nil, err
	}

	if len(colors) == 0 {
		return nil, errors.New("no source colors")
	}

	cfg := &Settings{
		Contrast: 0,
		Dark:     false,
		Variant:  VariantExpressive,
		Version:  Version2025,
		Platform: PlatformPhone,
	}
	for opt := range slices.Values(options) {
		opt(cfg)
	}

	if cfg.Context == nil {
		cfg.Context = context.Background()
	}

	source := colors[0]

	if len(colors) != 1 {
		quantized, err := quantizer.QuantizeCelebiContext(
			cfg.Context, colors, 5,
		)
		if err != nil {
			return nil, err
		}
		if len(quantized) == 0 {
			return nil, errNoColorFound
		}

		scored := score.Score(quantized)
		if len(scored) == 0 {
			return nil, errNoColorFound
		}
		source = scored[0]
	}

	scheme := dynamic.NewDynamicScheme(
		source.ToHct(),
		cfg.Variant,
		cfg.Contrast,
		cfg.Dark,
		cfg.Platform,
		cfg.Version,
	)

	return createColors(scheme, cfg.Custom), nil
}
