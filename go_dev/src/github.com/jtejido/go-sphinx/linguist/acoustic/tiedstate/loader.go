package tiedstate

import (
	"fmt"

	"github.com/jtejido/go-sphinx/decoder/adaptation/model"
	"github.com/jtejido/go-sphinx/linguist/acoustic"
	"github.com/jtejido/go-sphinx/util"
	"github.com/jtejido/go-sphinx/util/props"
)

/** Generic interface for a loader of acoustic models */
type Loader interface {
	props.Configurable
	/**
	 * Loads the acoustic model.
	 *
	 * @throws IOException if an error occurs while loading the model
	 */
	Load() error

	/**
	 * Gets the pool of means for this loader.
	 *
	 * @return the pool
	 */
	MeansPool() *Pool[[]float32]

	/**
	 * Gets the means transformation matrix pool.
	 *
	 * @return the pool
	 */
	MeansTransformationMatrixPool() *Pool[[][]float32]

	/**
	 * Gets the means transformation vectors pool.
	 *
	 * @return the pool
	 */
	MeansTransformationVectorPool() *Pool[[]float32]

	/**
	 * Gets the variance pool.
	 *
	 * @return the pool
	 */
	VariancePool() *Pool[[]float32]

	/**
	 * Gets the variance transformation matrix pool.
	 *
	 * @return the pool
	 */
	VarianceTransformationMatrixPool() *Pool[[][]float32]

	/**
	 * Gets the variance transformation vectors pool.
	 *
	 * @return the pool
	 */
	VarianceTransformationVectorPool() *Pool[[]float32]

	/**
	 * Gets the mixture weight pool.
	 *
	 * @return the pool
	 */
	MixtureWeights() *GaussianWeights

	/**
	 * Gets the transition matrix pool.
	 *
	 * @return the pool
	 */
	TransitionMatrixPool() *Pool[[][]float32]

	/**
	 * Gets the transformation matrix.
	 *
	 * @return the matrix
	 */
	TransformMatrix() [][]float32

	/**
	 * Gets the senone pool for this loader.
	 *
	 * @return the pool
	 */
	SenonePool() *Pool[Senone]

	/**
	 * Returns the HMM Manager for this loader.
	 *
	 * @return the HMM Manager
	 */
	HMMManager() *HMMManager

	/**
	 * Returns the map of context indepent units. The map can be accessed by unit name.
	 *
	 * @return the map of context independent units
	 */
	ContextIndependentUnits() map[string]*acoustic.Unit

	/** logs information about this loader */
	LogInfo()

	/**
	 * Returns the size of the left context for context dependent units.
	 *
	 * @return the left context size
	 */
	LeftContextSize() int

	/**
	 * Returns the size of the right context for context dependent units.
	 *
	 * @return the left context size
	 */
	RightContextSize() int

	/**
	 * @return the model properties
	 */
	Properties() *util.Properties

	/**
	 * Apply the transform
	 * @param transform transform to apply to the model
	 * @param clusters transform clusters
	 */
	Update(transform model.Transform, clusters model.ClusteredDensityFileData)
}

func assert(ok bool) {
	if !ok {
		panic("assert fail")
	}
}

func assert2(ok bool, msg string, args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(msg, args...))
	}
}
