package color

import (
	"math"
	"testing"
)

func TestMakeViewingConditions_Default(t *testing.T) {
	adaptingLuminance := (200.0 / math.Pi) * YFromLstar(50.0) / 100.0
	backgroundLstar := 50.0
	surround := 2.0
	discountingIlluminant := false

	v := MakeViewingConditions(adaptingLuminance, backgroundLstar, surround, discountingIlluminant)

	t.Logf("ViewingConditions:")
	t.Logf("  N:      %.6f", v.N)
	t.Logf("  Aw:     %.6f", v.Aw)
	t.Logf("  Nbb:    %.6f", v.Nbb)
	t.Logf("  Ncb:    %.6f", v.Ncb)
	t.Logf("  C:      %.6f", v.C)
	t.Logf("  Nc:     %.6f", v.Nc)
	t.Logf("  RgbD:   [%.6f, %.6f, %.6f]", v.RgbD[0], v.RgbD[1], v.RgbD[2])
	t.Logf("  Fl:     %.6f", v.Fl)
	t.Logf("  FlRoot: %.6f", v.FlRoot)
	t.Logf("  Z:      %.6f", v.Z)
}
