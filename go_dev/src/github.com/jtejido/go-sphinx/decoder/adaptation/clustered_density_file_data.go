package adaptation

import (
	"math"
	"math/rand"
	"time"

	"github.com/jtejido/go-sphinx/linguist/acoustic/tiedstate"
)

/**
 * Used for clustering gaussians. The clustering is performed by Euclidean
 * distance criterion. The "k-means" clustering algorithm is used for clustering
 * the gaussians.
 *
 * @author Bogdan Petcu
 */
type ClusteredDensityFileData struct {
	numberOfClusters  int
	corespondingClass []int
}

func NewClusteredDensityFileData(loader tiedstate.Loader, numberOfClusters int) *ClusteredDensityFileData {
	res := new(ClusteredDensityFileData)
	res.numberOfClusters = numberOfClusters
	res.kMeansClustering(loader, 30)
	return res
}

func (c *ClusteredDensityFileData) NumberOfClusters() int {
	return c.numberOfClusters
}

/**
 * Used for accessing the index that is specific to a gaussian.
 *
 * @param gaussian
 *            provided in a i * numStates + gaussianIndex form.
 * @return class index
 */
func (c *ClusteredDensityFileData) ClassIndex(gaussian int) int {
	return c.corespondingClass[gaussian]
}

/**
 * Computes euclidean distance between 2 n-dimensional points.
 *
 * @param a
 *            - n-dimensional "a" point
 * @param b
 *            - n-dimensional "b" point
 * @return the euclidean distance between a and b.
 */
func (c *ClusteredDensityFileData) euclidianDistance(a, b []float32) float32 {
	var s float64
	var d float32

	for i := 0; i < len(a); i++ {
		d = a[i] - b[i]
		s += float64(d) * float64(d)
	}

	return float32(math.Sqrt(s))
}

/**
 * Checks if the two float array have the same components
 *
 * @param a
 *            - float array a
 * @param b
 *            - float array b
 * @return true if values from a are equal to the ones in b, else false.
 */
func (c *ClusteredDensityFileData) isEqual(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

/**
 * Performs k-means-clustering algorithm for clustering gaussians.
 * Clustering is done using euclidean distance criterium.
 *
 * @param maxIterations
 */
func (c *ClusteredDensityFileData) kMeansClustering(loader tiedstate.Loader, maxIterations int) {
	initialData := loader.MeansPool()
	oldCentroids := make([][]float32, c.numberOfClusters)

	centroids := make([][]float32, c.numberOfClusters)
	numberOfElements := initialData.Size()
	nrOfIterations := maxIterations
	var index int
	count := make([]int, c.numberOfClusters)
	var distance, min float64
	var currentValue, centroid []float32
	array := make([][][]float32, c.numberOfClusters)
	for i := 0; i < c.numberOfClusters; i++ {
		array[i] = make([][]float32, numberOfElements)
	}

	var converged bool
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for i := 0; i < c.numberOfClusters; i++ {
		index = rng.Intn(numberOfElements)
		centroids = append(centroids, initialData.Get(index))
		oldCentroids = append(oldCentroids, initialData.Get(index))
		count[i] = 0
	}

	index = 0

	for !converged && nrOfIterations > 0 {
		c.corespondingClass = make([]int, initialData.Size())
		array = make([][][]float32, c.numberOfClusters)
		for i := 0; i < c.numberOfClusters; i++ {
			array[i] = make([][]float32, numberOfElements)
		}
		for i := 0; i < c.numberOfClusters; i++ {
			oldCentroids[i] = centroids[i]
			count[i] = 0
		}

		for i := 0; i < initialData.Size(); i++ {
			currentValue = initialData.Get(i)
			min = float64(c.euclidianDistance(oldCentroids[0], currentValue))
			index = 0

			for k := 1; k < c.numberOfClusters; k++ {
				distance = float64(c.euclidianDistance(oldCentroids[k], currentValue))

				if distance < min {
					min = distance
					index = k
				}
			}

			array[index][count[index]] = currentValue
			c.corespondingClass[i] = index
			count[index]++

		}

		for i := 0; i < c.numberOfClusters; i++ {
			centroid = make([]float32, len(initialData.Get(0)))

			if count[i] > 0 {

				for j := 0; j < count[i]; j++ {
					for k := 0; k < len(initialData.Get(0)); k++ {
						centroid[k] += array[i][j][k]
					}
				}

				for k := 0; k < len(initialData.Get(0)); k++ {
					centroid[k] /= float32(count[i])
				}

				centroids[i] = centroid
			}
		}

		converged = true

		for i := 0; i < c.numberOfClusters; i++ {
			converged = converged && (c.isEqual(centroids[i], oldCentroids[i]))
		}

		nrOfIterations--
	}
}
