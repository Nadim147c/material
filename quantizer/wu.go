package quantizer

import (
	"math"

	"github.com/Nadim147c/goyou/color"
)

const (
	indexBits int = 5
	histSize  int = 32  // 32 bins
	cubeSize  int = 33  // 32 bins + 1 for cumulative indexing
	maxColors int = 256 // or any desired palette size
)

func index(r, g, b int) int {
	i := (r * cubeSize * cubeSize) + (g * cubeSize) + b
	return i
}

type direction int

const (
	directionRed direction = iota
	directionGreen
	directionBlue
)

type box struct {
	R0, R1 int
	G0, G1 int
	B0, B1 int
	vol    uint32
}

type colorCube struct {
	Boxes [maxColors]box
	Next  int // how many boxes created so far
}

type maximizeResult struct {
	cutLocation int
	maximum     int
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
	size := (cubeSize * cubeSize * cubeSize) + (cubeSize * cubeSize) + cubeSize + histSize
	q := quantizerWu{
		weights:  make([]float64, size),
		momentsR: make([]float64, size),
		momentsG: make([]float64, size),
		momentsB: make([]float64, size),
		moments:  make([]float64, size),
	}

	return q.Quantize(input, maxColor)
}

func (q *quantizerWu) Quantize(input pixels, maxColor int) pixels {
	q.buildHistogram(input)
	q.computeMoments()
	r := q.createBoxes(maxColor)
	return q.createResult(r)
}

func (q *quantizerWu) buildHistogram(pixels []color.Color) {
	for pixel, count := range QuantizeMap(pixels) {
		_, r, g, b := pixel.Values() // ignore alpha

		ri := int(r) * histSize / 256
		gi := int(g) * histSize / 256
		bi := int(b) * histSize / 256

		i := index(ri, gi, bi)

		q.weights[i]++
		q.momentsR[i] += float64(count + int(r))
		q.momentsG[i] += float64(count + int(g))
		q.momentsB[i] += float64(count + int(b))
		q.moments[i] += float64(count) + float64(r*r+g*g+b*b)
	}
}

func (q *quantizerWu) createResult(maxColor int) pixels {
	colors := make(pixels, 0)
	for i := range maxColor {
		cube := q.cubes[i]
		weight := q.volume(&cube, q.weights)
		if weight > 0 {
			r := uint32(math.Round(q.volume(&cube, q.momentsR) / weight))
			g := uint32(math.Round(q.volume(&cube, q.momentsG) / weight))
			b := uint32(math.Round(q.volume(&cube, q.momentsB) / weight))
			color := color.Color((255 << 24) | ((r & 0x0FF) << 16) | ((g & 0x0FF) << 8) | (b & 0x0FF))
			colors = append(colors, color)
		}
	}
	return colors
}

func (q *quantizerWu) createBoxes(maxColors int) int {
	q.cubes = make([]box, maxColors)
	volumeVariance := make([]float64, maxColors)

	q.cubes[0] = box{
		R0: 0, G0: 0, B0: 0,
		R1: cubeSize, G1: cubeSize, B1: cubeSize,
	}

	generatedColorCount := maxColors
	next := 0
	for i := 1; i < maxColors; i++ {
		if q.cut(&q.cubes[next], &q.cubes[i]) {
			q.cubes[next].vol = 0
			if volumeVariance[next] > 1 {
				q.cubes[next].vol = q.variance(&q.cubes[next])
			}
			q.cubes[i].vol = 0
			if volumeVariance[i] > 1 {
				q.cubes[i].vol = q.variance(&q.cubes[i])
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

func (q *quantizerWu) computeMoments() {
	area := make([]float64, cubeSize)
	areaR := make([]float64, cubeSize)
	areaG := make([]float64, cubeSize)
	areaB := make([]float64, cubeSize)
	area2 := make([]float64, cubeSize)

	for r := 1; r < cubeSize; r++ {
		// Reset area accumulators
		for i := range cubeSize {
			area[i] = 0
			areaR[i] = 0
			areaG[i] = 0
			areaB[i] = 0
			area2[i] = 0
		}

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

func (q *quantizerWu) cut(one *box, two *box) bool {
	wholeR := q.volume(one, q.momentsR)
	wholeG := q.volume(one, q.momentsG)
	wholeB := q.volume(one, q.momentsB)
	wholeW := q.volume(one, q.weights)

	maxRResult := q.maximize(one, directionRed, one.R0+1, one.R1,
		wholeR, wholeG, wholeB, wholeW)
	maxGResult := q.maximize(one, directionGreen, one.G0+1, one.G1,
		wholeR, wholeG, wholeB, wholeW)
	maxBResult := q.maximize(one, directionBlue, one.B0+1, one.B1,
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

	two.R1 = one.R1
	two.G1 = one.G1
	two.B1 = one.B1

	switch direction {
	case directionRed:
		one.R1 = maxRResult.cutLocation
		two.R0 = one.R1
		two.G0 = one.G0
		two.B0 = one.B0
		break
	case directionGreen:
		one.G1 = maxGResult.cutLocation
		two.R0 = one.R0
		two.G0 = one.G1
		two.B0 = one.B0
		break
	case directionBlue:
		one.B1 = maxBResult.cutLocation
		two.R0 = one.R0
		two.G0 = one.G0
		two.B0 = one.B1
		break
	default:
		panic("unexpected direction")
	}

	one.vol = uint32((one.R1 - one.R0) * (one.G1 - one.G0) * (one.B1 - one.B0))
	two.vol = uint32((two.R1 - two.R0) * (two.G1 - two.G0) * (two.B1 - two.B0))
	return true
}

func (q *quantizerWu) variance(cube *box) uint32 {
	dr := q.volume(cube, q.momentsR)
	dg := q.volume(cube, q.momentsG)
	db := q.volume(cube, q.momentsB)
	xx := q.moments[index(cube.R1, cube.G1, cube.B1)] -
		q.moments[index(cube.R1, cube.G1, cube.B0)] -
		q.moments[index(cube.R1, cube.G0, cube.B1)] +
		q.moments[index(cube.R1, cube.G0, cube.B0)] -
		q.moments[index(cube.R0, cube.G1, cube.B1)] +
		q.moments[index(cube.R0, cube.G1, cube.B0)] +
		q.moments[index(cube.R0, cube.G0, cube.B1)] -
		q.moments[index(cube.R0, cube.G0, cube.B0)]
	hypotenuse := dr*dr + dg*dg + db*db
	volume := q.volume(cube, q.weights)
	return uint32(xx - float64(hypotenuse)/float64(volume))
}

func (q *quantizerWu) bottom(cube *box, direction direction, moment []float64) float64 {
	switch direction {
	case directionRed:
		return (-moment[index(cube.R0, cube.G1, cube.B1)] +
			moment[index(cube.R0, cube.G1, cube.B0)] +
			moment[index(cube.R0, cube.G0, cube.B1)] -
			moment[index(cube.R0, cube.G0, cube.B0)])
	case directionGreen:
		return (-moment[index(cube.R1, cube.G0, cube.B1)] +
			moment[index(cube.R1, cube.G0, cube.B0)] +
			moment[index(cube.R0, cube.G0, cube.B1)] -
			moment[index(cube.R0, cube.G0, cube.B0)])
	case directionBlue:
		return (-moment[index(cube.R1, cube.G1, cube.B0)] +
			moment[index(cube.R1, cube.G0, cube.B0)] +
			moment[index(cube.R0, cube.G1, cube.B0)] -
			moment[index(cube.R0, cube.G0, cube.B0)])
	default:
		panic("Raise condition")
	}
}

func (q *quantizerWu) top(cube *box, direction direction, position int, moment []float64) float64 {
	switch direction {
	case directionRed:
		return (moment[index(position, cube.G1, cube.B1)] -
			moment[index(position, cube.G1, cube.B0)] -
			moment[index(position, cube.G0, cube.B1)] +
			moment[index(position, cube.G0, cube.B0)])
	case directionGreen:
		return (moment[index(cube.R1, position, cube.B1)] -
			moment[index(cube.R1, position, cube.B0)] -
			moment[index(cube.R0, position, cube.B1)] +
			moment[index(cube.R0, position, cube.B0)])
	case directionBlue:
		return (moment[index(cube.R1, cube.G1, position)] -
			moment[index(cube.R1, cube.G0, position)] -
			moment[index(cube.R0, cube.G1, position)] +
			moment[index(cube.R0, cube.G0, position)])
	default:
		panic("Raise condition")
	}
}

func (q *quantizerWu) maximize(
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

	m := int(math.Round(maxVal))

	return maximizeResult{cut, m}
}

func (q *quantizerWu) volume(cube *box, moment []float64) float64 {
	return moment[index(cube.R1, cube.G1, cube.B1)] -
		moment[index(cube.R1, cube.G1, cube.B0)] -
		moment[index(cube.R1, cube.G0, cube.B1)] +
		moment[index(cube.R1, cube.G0, cube.B0)] -
		moment[index(cube.R0, cube.G1, cube.B1)] +
		moment[index(cube.R0, cube.G1, cube.B0)] +
		moment[index(cube.R0, cube.G0, cube.B1)] -
		moment[index(cube.R0, cube.G0, cube.B0)]
}
