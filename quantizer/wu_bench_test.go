package quantizer

import (
	"math/rand"
	"testing"

	"github.com/Nadim147c/material/v2/color"
)

func generateTestImage(width, height int) []color.ARGB {
	pixels := make([]color.ARGB, width*height)
	for i := range pixels {
		// Generate random colors
		pixels[i] = color.ARGB(uint32(rand.Intn(0x1000000)) | 0xFF000000)
	}
	return pixels
}

func BenchmarkQuantizeWu_SmallImage(b *testing.B) {
	pixels := generateTestImage(100, 100)

	b.ReportAllocs()

	for b.Loop() {
		_ = QuantizeWu(pixels, 8)
	}
}

func BenchmarkQuantizeWu_MediumImage(b *testing.B) {
	pixels := generateTestImage(512, 512)

	b.ReportAllocs()

	for b.Loop() {
		_ = QuantizeWu(pixels, 16)
	}
}

func BenchmarkQuantizeWu_LargeImage(b *testing.B) {
	pixels := generateTestImage(1024, 768)

	b.ReportAllocs()

	for b.Loop() {
		_ = QuantizeWu(pixels, 32)
	}
}

func BenchmarkQuantizeWu_VeryLargeImage(b *testing.B) {
	pixels := generateTestImage(4000, 4000)

	b.ReportAllocs()

	for b.Loop() {
		_ = QuantizeWu(pixels, 32)
	}
}

func BenchmarkQuantizeWu_VaryingColors(b *testing.B) {
	pixels := generateTestImage(512, 512)

	tests := []struct {
		name   string
		colors int
	}{
		{"4 colors", 4},
		{"8 colors", 8},
		{"16 colors", 16},
		{"32 colors", 32},
		{"64 colors", 64},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_ = QuantizeWu(pixels, tc.colors)
			}
		})
	}
}
