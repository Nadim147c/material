package quantizer

import (
	"context"

	"github.com/Nadim147c/material/color"
)

const (
	indexBits    int64 = 5
	bitsToRemove int64 = 8 - indexBits
	histSize     int64 = 32 // 32 bins
	cubeSize     int64 = 33 // 32 bins + 1 for cumulative indexing
	totalSize    int64 = 35937
)

func index(r, g, b int64) int64 {
	return (r << (indexBits * 2)) + (r << (indexBits + 1)) + (g << indexBits) + r + g + b
}

type direction int

const (
	directionRed direction = iota
	directionGreen
	directionBlue
)

type box struct {
	r0, r1 int64
	g0, g1 int64
	b0, b1 int64
	vol    float64
}

type maximizeResult struct {
	cutLocation int64
	maximum     int64
}

type quantizerWu struct {
	weights  []int64
	momentsR []int64
	momentsG []int64
	momentsB []int64
	moments  []int64
	cubes    []box
}

// QuantizeWu is an image quantizer that divides the image's pixels into
// clusters by recursively cutting an RGB cube, based on the weight of pixels in
// each area of the cube.
//
// The algorithm was described by Xiaolin Wu in Graphic Gems II, published in
// 1991.
func QuantizeWu(input []color.ARGB, maxColor int) []color.ARGB {
	// ignore error because background context won't return any error
	qw, _ := QuantizeWuContext(context.Background(), input, maxColor)
	return qw
}

// QuantizeWuWithContext is QuantizeWu with context.Context support.
//
// Deprecated: Use QuantizeWuContext
func QuantizeWuWithContext(
	ctx context.Context,
	input []color.ARGB,
	maxColor int,
) ([]color.ARGB, error) {
	return QuantizeWuContext(ctx, input, maxColor)
}

// QuantizeWuContext is QuantizeWu with context.Context support.
func QuantizeWuContext(
	ctx context.Context,
	input []color.ARGB,
	maxColor int,
) ([]color.ARGB, error) {
	q := quantizerWu{
		weights:  make([]int64, totalSize),
		momentsR: make([]int64, totalSize),
		momentsG: make([]int64, totalSize),
		momentsB: make([]int64, totalSize),
		moments:  make([]int64, totalSize),
	}

	return q.Quantize(ctx, input, maxColor)
}

func (q *quantizerWu) Quantize(
	ctx context.Context,
	input []color.ARGB,
	maxColor int,
) ([]color.ARGB, error) {
	q.BuildHistogram(ctx, input)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	q.ComputeMoments(ctx)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	r := q.CreateBoxes(ctx, maxColor)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return q.CreateResult(r), nil
}

func (q *quantizerWu) BuildHistogram(ctx context.Context, pixels []color.ARGB) {
	for pixel, c := range QuantizeMap(pixels) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		count := int64(c)
		red := int64(pixel.Red())
		green := int64(pixel.Green())
		blue := int64(pixel.Blue())

		ri := red>>bitsToRemove + 1
		gi := green>>bitsToRemove + 1
		bi := blue>>bitsToRemove + 1

		i := index(ri, gi, bi)

		q.weights[i] += count
		q.momentsR[i] += count * red
		q.momentsG[i] += count * green
		q.momentsB[i] += count * blue
		q.moments[i] += count * (red*red + green*green + blue*blue)
	}
}

func (q *quantizerWu) CreateResult(maxColor int) []color.ARGB {
	colors := make([]color.ARGB, 0)
	for i := range maxColor {
		cube := q.cubes[i]
		weight := q.Volume(&cube, q.weights)
		if weight > 0 {
			r := uint32(q.Volume(&cube, q.momentsR) / weight)
			g := uint32(q.Volume(&cube, q.momentsG) / weight)
			b := uint32(q.Volume(&cube, q.momentsB) / weight)
			c := color.ARGB(
				(255 << 24) | ((r & 0x0FF) << 16) | ((g & 0x0FF) << 8) | (b & 0x0FF),
			)
			colors = append(colors, c)
		}
	}
	return colors
}

func (q *quantizerWu) CreateBoxes(ctx context.Context, maxColors int) int {
	q.cubes = make([]box, maxColors)
	volumeVariance := make([]int64, maxColors)

	q.cubes[0] = box{
		r0: 0, g0: 0, b0: 0,
		r1: histSize, g1: histSize, b1: histSize,
	}

	generatedColorCount := maxColors
	next := 0
	for i := 1; i < maxColors; i++ {
		select {
		case <-ctx.Done():
			return generatedColorCount
		default:
		}
		if q.Cut(&q.cubes[next], &q.cubes[i]) {
			volumeVariance[next] = 0
			if q.cubes[next].vol > 1 {
				volumeVariance[next] = q.Variance(&q.cubes[next])
			}
			volumeVariance[i] = 0
			if q.cubes[i].vol > 1 {
				volumeVariance[i] = q.Variance(&q.cubes[i])
			}
		} else {
			volumeVariance[next] = 0.0
			i--
		}

		next = 0
		temp := volumeVariance[0]
		for j := 1; j <= i; j++ {
			if volumeVariance[j] > temp {
				temp = volumeVariance[j]
				next = j
			}
		}
		if temp <= 0.0 {
			generatedColorCount = i + 1
			break
		}
	}
	return generatedColorCount
}

func (q *quantizerWu) ComputeMoments(ctx context.Context) {
	area := make([]int64, cubeSize)
	areaR := make([]int64, cubeSize)
	areaG := make([]int64, cubeSize)
	areaB := make([]int64, cubeSize)
	area2 := make([]int64, cubeSize)

	for r := int64(1); r < cubeSize; r++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for g := int64(1); g < cubeSize; g++ {
			var line, line2, lineR, lineG, lineB int64
			for b := int64(1); b < cubeSize; b++ {
				idx := index(r, g, b)

				line += q.weights[idx]
				lineR += q.momentsR[idx]
				lineG += q.momentsG[idx]
				lineB += q.momentsB[idx]
				line2 += q.moments[idx]

				area[b] += line
				areaR[b] += lineR
				areaG[b] += lineG
				areaB[b] += lineB
				area2[b] += line2

				previousIndex := index(r-1, g, b)
				q.weights[idx] = q.weights[previousIndex] + area[b]
				q.momentsR[idx] = q.momentsR[previousIndex] + areaR[b]
				q.momentsG[idx] = q.momentsG[previousIndex] + areaG[b]
				q.momentsB[idx] = q.momentsB[previousIndex] + areaB[b]
				q.moments[idx] = q.moments[previousIndex] + area2[b]
			}
		}
	}
}

func (q *quantizerWu) Cut(one *box, two *box) bool {
	wholeR := q.Volume(one, q.momentsR)
	wholeG := q.Volume(one, q.momentsG)
	wholeB := q.Volume(one, q.momentsB)
	wholeW := q.Volume(one, q.weights)

	maxRResult := q.Maximize(one, directionRed, one.r0+1, one.r1,
		wholeR, wholeG, wholeB, wholeW)
	maxGResult := q.Maximize(one, directionGreen, one.g0+1, one.g1,
		wholeR, wholeG, wholeB, wholeW)
	maxBResult := q.Maximize(one, directionBlue, one.b0+1, one.b1,
		wholeR, wholeG, wholeB, wholeW)

	var direction direction

	maxR := maxRResult.maximum
	maxG := maxGResult.maximum
	maxB := maxBResult.maximum
	if maxR >= maxG && maxR >= maxB {
		if maxRResult.cutLocation < 0 {
			return false
		}
		direction = directionRed
	} else if maxG >= maxR && maxG >= maxB {
		direction = directionGreen
	} else {
		direction = directionBlue
	}

	two.r1 = one.r1
	two.g1 = one.g1
	two.b1 = one.b1

	switch direction {
	case directionRed:
		one.r1 = maxRResult.cutLocation
		two.r0 = one.r1
		two.g0 = one.g0
		two.b0 = one.b0
	case directionGreen:
		one.g1 = maxGResult.cutLocation
		two.r0 = one.r0
		two.g0 = one.g1
		two.b0 = one.b0
	case directionBlue:
		one.b1 = maxBResult.cutLocation
		two.r0 = one.r0
		two.g0 = one.g0
		two.b0 = one.b1
	default:
		panic("unexpected direction")
	}

	one.vol = float64((one.r1 - one.r0) * (one.g1 - one.g0) * (one.b1 - one.b0))
	two.vol = float64((two.r1 - two.r0) * (two.g1 - two.g0) * (two.b1 - two.b0))
	return true
}

func (q *quantizerWu) Variance(cube *box) int64 {
	dr := q.Volume(cube, q.momentsR)
	dg := q.Volume(cube, q.momentsG)
	db := q.Volume(cube, q.momentsB)
	xx := q.moments[index(cube.r1, cube.g1, cube.b1)] -
		q.moments[index(cube.r1, cube.g1, cube.b0)] -
		q.moments[index(cube.r1, cube.g0, cube.b1)] +
		q.moments[index(cube.r1, cube.g0, cube.b0)] -
		q.moments[index(cube.r0, cube.g1, cube.b1)] +
		q.moments[index(cube.r0, cube.g1, cube.b0)] +
		q.moments[index(cube.r0, cube.g0, cube.b1)] -
		q.moments[index(cube.r0, cube.g0, cube.b0)]
	hypotenuse := dr*dr + dg*dg + db*db
	volume := q.Volume(cube, q.weights)
	return xx - hypotenuse/volume
}

func (q *quantizerWu) bottom(
	cube *box,
	direction direction,
	moment []int64,
) int64 {
	switch direction {
	case directionRed:
		return (-moment[index(cube.r0, cube.g1, cube.b1)] +
			moment[index(cube.r0, cube.g1, cube.b0)] +
			moment[index(cube.r0, cube.g0, cube.b1)] -
			moment[index(cube.r0, cube.g0, cube.b0)])
	case directionGreen:
		return (-moment[index(cube.r1, cube.g0, cube.b1)] +
			moment[index(cube.r1, cube.g0, cube.b0)] +
			moment[index(cube.r0, cube.g0, cube.b1)] -
			moment[index(cube.r0, cube.g0, cube.b0)])
	case directionBlue:
		return (-moment[index(cube.r1, cube.g1, cube.b0)] +
			moment[index(cube.r1, cube.g0, cube.b0)] +
			moment[index(cube.r0, cube.g1, cube.b0)] -
			moment[index(cube.r0, cube.g0, cube.b0)])
	default:
		panic("Raise condition")
	}
}

func (q *quantizerWu) top(
	cube *box,
	direction direction,
	position int64,
	moment []int64,
) int64 {
	switch direction {
	case directionRed:
		return (moment[index(position, cube.g1, cube.b1)] -
			moment[index(position, cube.g1, cube.b0)] -
			moment[index(position, cube.g0, cube.b1)] +
			moment[index(position, cube.g0, cube.b0)])
	case directionGreen:
		return (moment[index(cube.r1, position, cube.b1)] -
			moment[index(cube.r1, position, cube.b0)] -
			moment[index(cube.r0, position, cube.b1)] +
			moment[index(cube.r0, position, cube.b0)])
	case directionBlue:
		return (moment[index(cube.r1, cube.g1, position)] -
			moment[index(cube.r1, cube.g0, position)] -
			moment[index(cube.r0, cube.g1, position)] +
			moment[index(cube.r0, cube.g0, position)])
	default:
		panic("Raise condition")
	}
}

func (q *quantizerWu) Maximize(
	cube *box, direction direction, first int64, last int64,
	wholeR int64, wholeG int64, wholeB int64, wholeW int64,
) maximizeResult {
	bottomR := q.bottom(cube, direction, q.momentsR)
	bottomG := q.bottom(cube, direction, q.momentsG)
	bottomB := q.bottom(cube, direction, q.momentsB)
	bottomW := q.bottom(cube, direction, q.weights)

	var maxVal int64
	var cut int64 = -1

	var halfR, halfG, halfB, halfW int64
	for i := first; i < last; i++ {
		halfR = bottomR + q.top(cube, direction, i, q.momentsR)
		halfG = bottomG + q.top(cube, direction, i, q.momentsG)
		halfB = bottomB + q.top(cube, direction, i, q.momentsB)
		halfW = bottomW + q.top(cube, direction, i, q.weights)
		if halfW == 0 {
			continue
		}

		temp := (halfR*halfR + halfG*halfG + halfB*halfB) / halfW

		halfR = wholeR - halfR
		halfG = wholeG - halfG
		halfB = wholeB - halfB
		halfW = wholeW - halfW
		if halfW == 0 {
			continue
		}

		temp += (halfR*halfR + halfG*halfG + halfB*halfB) / halfW

		if temp > maxVal {
			maxVal = temp
			cut = i
		}
	}

	return maximizeResult{cut, maxVal}
}

func (q *quantizerWu) Volume(cube *box, moment []int64) int64 {
	return (moment[index(cube.r1, cube.g1, cube.b1)] -
		moment[index(cube.r1, cube.g1, cube.b0)] -
		moment[index(cube.r1, cube.g0, cube.b1)] +
		moment[index(cube.r1, cube.g0, cube.b0)] -
		moment[index(cube.r0, cube.g1, cube.b1)] +
		moment[index(cube.r0, cube.g1, cube.b0)] +
		moment[index(cube.r0, cube.g0, cube.b1)] -
		moment[index(cube.r0, cube.g0, cube.b0)])
}
