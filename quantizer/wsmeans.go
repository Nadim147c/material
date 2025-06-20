package quantizer

import (
	"math"
	"math/rand"
	"slices"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
)

const (
	MaxIterations       int     = 10
	MinMovementDistance float64 = 3.0
)

type distanceAndIndex struct {
	distance float64
	index    int
}

type (
	pixels    = []color.ARGB
	pixelsLab = []color.Lab

	QuantizedMap = map[color.ARGB]int
)

func QuantizeWsMeans(input pixels, startingClusters []color.Lab, maxColors int) QuantizedMap {
	// Get color frequncies
	freq := make(map[color.ARGB]int)
	for c := range slices.Values(input) {
		freq[c]++
	}

	// Number of unique color in the image/pixels array
	pointCount := len(freq)
	points := make(pixelsLab, pointCount)
	counts := make([]int, pointCount)
	i := 0
	for k, v := range freq {
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

	distanceToIndexMatrix := make([][]distanceAndIndex, clusterCount)
	for i := range distanceToIndexMatrix {
		distanceToIndexMatrix[i] = make([]distanceAndIndex, len(clusters))
	}

	for iteration := range MaxIterations {
		// Step 1: Compute cluster-to-cluster distances
		for i := range clusterCount {
			for j := range clusterCount {
				distance := clusters[i].DistanceSquared(clusters[j])
				distanceToIndexMatrix[j][i].distance = distance
				distanceToIndexMatrix[j][i].index = i
				distanceToIndexMatrix[i][j].distance = distance
				distanceToIndexMatrix[i][j].index = j
			}

			slices.SortFunc(distanceToIndexMatrix[i], func(a, b distanceAndIndex) int {
				return num.SignCmp(a.distance, b.distance)
			})

			for j := range clusterCount {
				indexMatrix[i][j] = distanceToIndexMatrix[i][j].index
			}
		}

		pointsMoved := 0
		for i, point := range points {
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
				distanceChange := math.Abs((math.Sqrt(minimumDistance) - math.Sqrt(previousDistance)))
				if distanceChange > MinMovementDistance {
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

		argbToPopulation := make(QuantizedMap)

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

		return argbToPopulation

	}

	result := make(QuantizedMap)
	for lab := range slices.Values(clusters) {
		result[lab.ToARGB()]++
	}
	return result
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
