package quantizer

import (
	"math"
	"slices"

	"github.com/Nadim147c/goyou/color"
)

const (
	indexBits int = 5
	histSize  int = 32  // 32 bins
	cubeSize  int = 33  // 32 bins + 1 for cumulative indexing
	maxColors int = 256 // or any desired palette size
	totalSize int = 35937
)

func index(r, g, b int) int {
	return (r << (indexBits * 2)) + (r << (indexBits + 1)) + (g << indexBits) + r + g + b
}

type direction int

const (
	directionRed direction = iota
	directionGreen
	directionBlue
)

type box struct {
	r0, r1 int
	g0, g1 int
	b0, b1 int
	vol    float64
}

type colorCube struct {
	Boxes [maxColors]box
	Next  int // how many boxes created so far
}

type maximizeResult struct {
	cutLocation int
	maximum     float64
}

type createBoxesResult struct {
	MaxColors           int
	GeneratedColorCount int
}

type quantizerWu struct {
	weights  []float64
	momentsR []float64
	momentsG []float64
	momentsB []float64
	moments  []float64
	cubes    []box
}

func QuantizeWu(input pixels, maxColor int) pixels {
	q := quantizerWu{
		weights:  make([]float64, totalSize),
		momentsR: make([]float64, totalSize),
		momentsG: make([]float64, totalSize),
		momentsB: make([]float64, totalSize),
		moments:  make([]float64, totalSize),
	}

	return q.Quantize(input, maxColor)
}

func (q *quantizerWu) Quantize(input pixels, maxColor int) pixels {
	q.BuildHistogram(input)
	q.ComputeMoments()
	r := q.CreateBoxes(maxColor)
	return q.CreateResult(r)
}

func (q *quantizerWu) BuildHistogram(pixels []color.ARGB) {
	for pixel := range slices.Values(pixels) {
		_, r, g, b := pixel.Values() // ignore alpha

		ri := int(r) * histSize / 256
		gi := int(g) * histSize / 256
		bi := int(b) * histSize / 256

		i := index(ri, gi, bi)

		q.weights[i]++
		q.momentsR[i] += float64((r))
		q.momentsG[i] += float64((g))
		q.momentsB[i] += float64((b))
		q.moments[i] += float64(r*r + g*g + b*b)
	}
}

func (q *quantizerWu) CreateResult(maxColor int) pixels {
	colors := make(pixels, 0)
	for i := range maxColor {
		cube := q.cubes[i]
		weight := q.Volume(&cube, q.weights)
		if weight > 0 {
			r := uint32(math.Round(q.Volume(&cube, q.momentsR) / weight))
			g := uint32(math.Round(q.Volume(&cube, q.momentsG) / weight))
			b := uint32(math.Round(q.Volume(&cube, q.momentsB) / weight))
			color := color.ARGB((255 << 24) | ((r & 0x0FF) << 16) | ((g & 0x0FF) << 8) | (b & 0x0FF))
			colors = append(colors, color)
		}
	}
	return colors
}

func (q *quantizerWu) CreateBoxes(maxColors int) int {
	q.cubes = make([]box, maxColors)
	volumeVariance := make([]float64, maxColors)

	q.cubes[0] = box{
		r0: 0, g0: 0, b0: 0,
		r1: histSize, g1: histSize, b1: histSize,
	}

	generatedColorCount := maxColors
	next := 0
	for i := 1; i < maxColors; i++ {
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

func (q *quantizerWu) ComputeMoments() {
	area := make([]float64, cubeSize)
	areaR := make([]float64, cubeSize)
	areaG := make([]float64, cubeSize)
	areaB := make([]float64, cubeSize)
	area2 := make([]float64, cubeSize)

	for r := 1; r < cubeSize; r++ {
		for g := 1; g < cubeSize; g++ {
			var line, lineR, lineG, lineB float64
			line2 := float64(0)

			for b := 1; b < cubeSize; b++ {
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
		break
	case directionGreen:
		one.g1 = maxGResult.cutLocation
		two.r0 = one.r0
		two.g0 = one.g1
		two.b0 = one.b0
		break
	case directionBlue:
		one.b1 = maxBResult.cutLocation
		two.r0 = one.r0
		two.g0 = one.g0
		two.b0 = one.b1
		break
	default:
		panic("unexpected direction")
	}

	one.vol = float64((one.r1 - one.r0) * (one.g1 - one.g0) * (one.b1 - one.b0))
	two.vol = float64((two.r1 - two.r0) * (two.g1 - two.g0) * (two.b1 - two.b0))
	return true
}

func (q *quantizerWu) Variance(cube *box) float64 {
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

func (q *quantizerWu) bottom(cube *box, direction direction, moment []float64) float64 {
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

func (q *quantizerWu) top(cube *box, direction direction, position int, moment []float64) float64 {
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
	cube *box, direction direction, first int, last int,
	wholeR float64, wholeG float64, wholeB float64, wholeW float64,
) maximizeResult {
	bottomR := q.bottom(cube, direction, q.momentsR)
	bottomG := q.bottom(cube, direction, q.momentsG)
	bottomB := q.bottom(cube, direction, q.momentsB)
	bottomW := q.bottom(cube, direction, q.weights)

	maxVal := 0.0
	cut := -1

	var halfR, halfG, halfB, halfW float64
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

func (q *quantizerWu) Volume(cube *box, moment []float64) float64 {
	return (moment[index(cube.r1, cube.g1, cube.b1)] -
		moment[index(cube.r1, cube.g1, cube.b0)] -
		moment[index(cube.r1, cube.g0, cube.b1)] +
		moment[index(cube.r1, cube.g0, cube.b0)] -
		moment[index(cube.r0, cube.g1, cube.b1)] +
		moment[index(cube.r0, cube.g1, cube.b0)] +
		moment[index(cube.r0, cube.g0, cube.b1)] -
		moment[index(cube.r0, cube.g0, cube.b0)])
}
