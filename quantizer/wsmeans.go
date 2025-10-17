package quantizer

import (
	"context"
	"math"
	"math/rand"
	"slices"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
)

const (
	maxIterations       int     = 10
	minMovementDistance float64 = 3.0
)

type distanceIndex struct {
	distance float64
	index    int
}

type (
	pixels    = []color.ARGB
	pixelsLab = []color.Lab
)

// QuantizedMap is a map where ARGB is key and their frequencies as int
type QuantizedMap = map[color.ARGB]int

// QuantizeWsMeans is an image quantizer that improves on the speed of a
// standard K-Means algorithm by implementing several optimizations, including
// deduping identical pixels and a triangle inequality rule that reduces the
// number of comparisons needed to identify which cluster a point should be
// moved to.
//
// Wsmeans stands for Weighted Square Means.
//
// This algorithm was designed by M. Emre Celebi, and was found in their 2011
// paper, Improving the Performance of K-Means for Color Quantization.
// https://arxiv.org/abs/1101.0395
func QuantizeWsMeans(
	input pixels,
	startingClusters []color.Lab,
	maxColors int,
) QuantizedMap {
	// ignore error because background context won't return any error
	qm, _ := QuantizeWsMeansContext(
		context.Background(),
		input,
		startingClusters,
		maxColors,
	)
	return qm
}

// QuantizeWsMeansWithContext is QuantizeWsMeans with context.Context support.
//
// Deprecated: Use QuantizeWsWithContext
func QuantizeWsMeansWithContext(
	ctx context.Context,
	input pixels,
	startingClusters []color.Lab,
	maxColors int,
) (QuantizedMap, error) {
	return QuantizeWsMeansContext(ctx, input, startingClusters, maxColors)
}

// QuantizeWsMeansContext is QuantizeWsMeans with context.Context support.
func QuantizeWsMeansContext(
	ctx context.Context,
	input pixels,
	startingClusters []color.Lab,
	maxColors int,
) (QuantizedMap, error) {
	// Check context at the start
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Get color frequencies
	freq := QuantizedMap{}
	for c := range slices.Values(input) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		freq[c]++
	}

	// Number of unique color in the image/pixels array
	pointCount := len(freq)
	points := make(pixelsLab, pointCount)
	counts := make([]int, pointCount)
	i := 0
	for k, v := range freq {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		points[i] = k.ToLab()
		counts[i] = v
		i++
	}

	clusterCount := min(maxColors, pointCount)
	if len(startingClusters) != 0 {
		clusterCount = min(clusterCount, len(startingClusters))
	}

	clusters := make([]color.Lab, 0)

	for cluster := range slices.Values(startingClusters) {
		clusters = append(clusters, cluster)
	}

	clustersNeeded := clusterCount - len(clusters)
	if len(startingClusters) == 0 && clustersNeeded > 0 {
		clusters = append(clusters, randomLabClusters(clustersNeeded)...)
	}

	clusterIndices := make([]int, pointCount)
	for i := range clusterIndices {
		clusterIndices[i] = rand.Intn(clusterCount)
	}

	indexMatrix := make([][]int, clusterCount)
	for i := range indexMatrix {
		indexMatrix[i] = make([]int, clusterCount)
	}

	distanceToIndexMatrix := make([][]distanceIndex, clusterCount)
	for i := range distanceToIndexMatrix {
		distanceToIndexMatrix[i] = make([]distanceIndex, len(clusters))
	}

	for iteration := range maxIterations {
		// Check context at start of each iteration
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Step 1: Compute cluster-to-cluster distances
		for i := range clusterCount {
			// Check context in outer loop
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			for j := range clusterCount {
				distance := clusters[i].DistanceSquared(clusters[j])
				distanceToIndexMatrix[j][i].distance = distance
				distanceToIndexMatrix[j][i].index = i
				distanceToIndexMatrix[i][j].distance = distance
				distanceToIndexMatrix[i][j].index = j
			}

			slices.SortFunc(
				distanceToIndexMatrix[i],
				func(a, b distanceIndex) int {
					return num.SignCmp(a.distance, b.distance)
				},
			)

			for j := range clusterCount {
				indexMatrix[i][j] = distanceToIndexMatrix[i][j].index
			}
		}

		pointsMoved := 0
		for i, point := range points {
			// Check context periodically during point processing
			if i%100 == 0 { // Check every 100 points to avoid excessive overhead
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}
			}

			previousClusterIndex := clusterIndices[i]
			previousCluster := clusters[previousClusterIndex]
			previousDistance := point.DistanceSquared(previousCluster)
			minimumDistance := previousDistance
			newClusterIndex := -1

			for j := range clusterCount {
				if distanceToIndexMatrix[previousClusterIndex][j].distance >= 4*previousDistance {
					continue
				}
				distance := point.DistanceSquared(clusters[j])
				if distance < minimumDistance {
					minimumDistance = distance
					newClusterIndex = j
				}
			}

			if newClusterIndex != -1 {
				distanceChange := math.Abs(
					(math.Sqrt(minimumDistance) - math.Sqrt(previousDistance)),
				)
				if distanceChange > minMovementDistance {
					pointsMoved++
					clusterIndices[i] = newClusterIndex
				}
			}
		}

		if pointsMoved == 0 && iteration != 0 {
			break
		}

		component0Sums := make([]float64, clusterCount) // L
		component1Sums := make([]float64, clusterCount) // a
		component2Sums := make([]float64, clusterCount) // b
		pixelCountSums := make([]float64, clusterCount)

		// Accumulate weighted components
		for i := range pointCount {
			// Check context periodically during accumulation
			if i%1000 == 0 { // Check every 1000 points
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}
			}

			clusterIndex := clusterIndices[i]
			point := points[i]
			count := float64(counts[i])

			pixelCountSums[clusterIndex] += count
			component0Sums[clusterIndex] += point.L * count
			component1Sums[clusterIndex] += point.A * count
			component2Sums[clusterIndex] += point.B * count
		}

		// Compute new cluster centers
		for i := range clusterCount {
			count := pixelCountSums[i]
			if count == 0 {
				clusters[i] = color.NewLab(0.0, 0.0, 0.0)
				continue
			}
			l := component0Sums[i] / count
			a := component1Sums[i] / count
			b := component2Sums[i] / count
			clusters[i] = color.NewLab(l, a, b)
		}

		// Check if we should return results early (this seems to be the
		// original logic)
		argbToPopulation := QuantizedMap{}

		for i := range clusterCount {
			count := int(pixelCountSums[i])
			if count == 0 {
				continue
			}

			colorInt := clusters[i].ToARGB()
			if _, exists := argbToPopulation[colorInt]; exists {
				continue
			}

			argbToPopulation[colorInt] = count
		}

		return argbToPopulation, nil
	}

	// Final result if we exit the loop normally
	result := QuantizedMap{}
	for lab := range slices.Values(clusters) {
		result[lab.ToARGB()]++
	}
	return result, nil
}

func randomLabClusters(n int) []color.Lab {
	clusters := make([]color.Lab, n)
	for i := range n {
		l := rand.Float64() * 100.0
		a := rand.Float64()*200.0 - 100.0
		b := rand.Float64()*200.0 - 100.0
		clusters[i] = color.NewLab(l, a, b)
	}
	return clusters
}
