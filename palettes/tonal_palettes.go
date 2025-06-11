package palettes

import "github.com/Nadim147c/material/color"

type TonalPalette struct {
	cache    map[float64]color.ARGB
	Hue      float64
	Chroma   float64
	KeyColor color.Hct
}

func NewFromARGB(color color.ARGB) *TonalPalette {
	hct := color.ToHct()
	return NewFromHct(hct)
}

func NewFromHct(hct color.Hct) *TonalPalette {
	return &TonalPalette{
		Hue:      hct.Hue,
		Chroma:   hct.Chroma,
		KeyColor: hct,
	}
}

func FromHueAndChroma(hue, chroma float64) *TonalPalette {
	keyColor := NewKeyColor(hue, chroma).Create()
	return NewFromHct(keyColor)
}

func (tp *TonalPalette) Tone(tone float64) color.ARGB {
	if tp.cache == nil {
		tp.cache = make(map[float64]color.ARGB)
	}

	argb, ok := tp.cache[tone]
	if !ok {
		argb = color.NewHct(tp.Hue, tp.Chroma, tone).ToARGB()
		tp.cache[tone] = argb
	}
	return argb
}

func (tp *TonalPalette) Get(tone float64) color.ARGB {
	return tp.Tone(tone)
}

func (tp *TonalPalette) GetHct(tone float64) color.Hct {
	return tp.Tone(tone).ToHct()
}

// IsBlue determines if a hue is in the blue range.
func (tp *TonalPalette) IsBlue() bool {
	return tp.Hue >= 250 && tp.Hue < 270
}

// IsYellow determines if a hue is in the yellow range.
func (tp *TonalPalette) IsYellow() bool {
	return tp.Hue >= 105 && tp.Hue < 125
}

// IsCyan determines if a hue is in the cyan range.
func (tp *TonalPalette) IsCyan() bool {
	return tp.Hue >= 170 && tp.Hue < 207
}
