package main

import (
	"fmt"
	"slices"

	"github.com/Nadim147c/material/v2/color"
)

func gradient(start, end color.ARGB, count int) []color.ARGB {
	if count < 2 {
		return []color.ARGB{start}
	}

	from, to := start.ToOkLab(), end.ToOkLab()
	step := 1.0 / float64(count-1)

	colors := make([]color.ARGB, count)
	for i := range count {
		ratio := step * float64(i)
		l := from.L + (to.L-from.L)*ratio
		a := from.A + (to.A-from.A)*ratio
		b := from.B + (to.B-from.B)*ratio
		colors[i] = color.NewOkLab(l, a, b).ToARGB()
	}

	return colors
}

func main() {
	start := color.ARGBFromHexMust("#33ff33")
	end := color.ARGBFromHexMust("#4444ff")

	colors := gradient(start, end, 100)
	for c := range slices.Values(colors) {
		fmt.Print(c.AnsiBg(" "))
	}
	fmt.Println()
}
