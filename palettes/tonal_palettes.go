package palettes

import "github.com/Nadim147c/goyou/color"

type TonalPalette struct {
	cache    map[float64]color.Color
	Hue      float64
	Chroma   float64
	KeyColor *color.Hct
}

func NewFromColor(color color.Color) *TonalPalette {
	hct := color.ToHct()
	return NewFromHct(hct)
}

func NewFromHct(hct *color.Hct) *TonalPalette {
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

func (tp *TonalPalette) Tone(tone float64) color.Color {
	if tp.cache == nil {
		tp.cache = make(map[float64]color.Color)
	}

	argb, ok := tp.cache[tone]
	if !ok {
		argb = color.NewHct(tp.Hue, tp.Chroma, tone).ToColor()
		tp.cache[tone] = argb
	}
	return argb
}

func (tp *TonalPalette) GetHct(tone float64) *color.Hct {
	return tp.Tone(tone).ToHct()
}
