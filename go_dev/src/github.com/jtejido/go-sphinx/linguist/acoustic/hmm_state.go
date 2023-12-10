package acoustic

import (
	"github.com/jtejido/go-sphinx/frontend"
	"github.com/jtejido/go-sphinx/linguist/acoustic/tiedstate/model"
)

/** Represents a single state in an HMM */
type HMMState interface {

	/**
	 * Gets the HMM associated with this state
	 *
	 * @return the HMM
	 */
	HMM() HMM

	/**
	 * Returns the mixture components associated with this Gaussian
	 *
	 * @return the array of mixture components
	 */
	MixtureComponents() model.MixtureComponent

	/**
	 * Gets the id of the mixture
	 *
	 * @return the id
	 */
	MixtureId() int64

	/**
	 *
	 * @return the mixture weights vector
	 */
	LogMixtureWeights() []float32

	/**
	 * Gets the state
	 *
	 * @return the state
	 */
	State() int

	/**
	 * Gets the score for this HMM state
	 *
	 * @param data the data to be scored
	 * @return the acoustic score for this state.
	 */
	Score(data frontend.Data) float32

	CalculateComponentScore(data frontend.Data) []float32

	/**
	 * Determines if this HMMState is an emitting state
	 *
	 * @return true if the state is an emitting state
	 */
	IsEmitting() bool

	/**
	 * Retrieves the state of successor states for this state
	 *
	 * @return the set of successor state arcs
	 */
	Successors() *HMMStateArc

	/**
	 * Determines if this state is an exit state of the HMM
	 *
	 * @return true if the state is an exit state
	 */
	IsExitState() bool
}
