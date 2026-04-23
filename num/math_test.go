package num

import (
	"fmt"
	"testing"
)

func TestClamp(t *testing.T) {
	cases := []struct {
		low, high, value, expected float64
	}{
		{0, 1, 1.5, 1},
		{0, 100, 1.5, 1.5},
		{0, 0.5, 1.5, 0.5},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("clamp %.3f between [%.3f, %.3f]", tt.value, tt.low, tt.high)
		t.Run(name, func(t *testing.T) {
			got := Clamp(tt.low, tt.high, tt.value)
			if tt.expected != got {
				t.Errorf("expected %.3f but got %.3f", tt.expected, got)
			}
		})
	}
}

func TestNormalizeDegree(t *testing.T) {
	cases := []struct {
		value, expected float64
	}{
		{-5, 355},
		{180, 180},
		{100, 100},
		{720, 0},
		{520, 160},
		{-520, 200},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("normalize %.3f degree", tt.value)
		t.Run(name, func(t *testing.T) {
			got := NormalizeDegree(tt.value)
			if tt.expected != got {
				t.Errorf("expected %.3f but got %.3f", tt.expected, got)
			}
		})
	}
}

func TestRotationDirection(t *testing.T) {
	cases := []struct {
		from, to, expected float64
	}{
		{5, 10, Clockwise},
		{50, 10, CounterClockwise},
		{10, 10, NoRotation},
		{-50, 10, Clockwise},
		{350, -30, CounterClockwise},
	}

	names := map[float64]string{
		Clockwise:        "clockwise",
		CounterClockwise: "counter-clockwise",
		NoRotation:       "no-rotation",
	}

	for _, tt := range cases {
		name := fmt.Sprintf("get roration from %.3f to %.3f", tt.from, tt.to)
		t.Run(name, func(t *testing.T) {
			got := RotationDirection(tt.from, tt.to)
			if tt.expected != got {
				t.Errorf("expected %s but got %s", names[tt.expected], names[got])
			}
		})
	}
}

func TestDifferenceDegrees(t *testing.T) {
	cases := []struct {
		from, to, expected float64
	}{
		{5, 10, 5},
		{50, 10, 40},
		{10, 10, 0},
		{-50, 10, 60},
		{350, -30, 20},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("get roration from %.3f to %.3f", tt.from, tt.to)
		t.Run(name, func(t *testing.T) {
			got := DifferenceDegrees(tt.from, tt.to)
			if tt.expected != got {
				t.Errorf("expected %.4f but got %.4f", tt.expected, got)
			}
		})
	}
}
