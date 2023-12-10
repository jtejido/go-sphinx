package acoustic

import "fmt"

/**
 * Represents a transition to single state in an HMM
 * <p>
 * All probabilities are maintained in linear base
 */
type HMMStateArc struct {
	hmmState    HMMState
	probability float32
}

/**
 * Constructs an HMMStateArc
 *
 * @param hmmState    destination state for this arc
 * @param probability the probability for this transition
 */
func NewHMMStateArc(hmmState HMMState, probability float32) *HMMStateArc {
	return &HMMStateArc{
		hmmState:    hmmState,
		probability: probability,
	}
}

/**
 * Gets the HMM associated with this state
 *
 * @return the HMM
 */
func (a *HMMStateArc) HMMState() HMMState {
	return a.hmmState
}

/**
 * Gets log transition probability
 *
 * @return the probability in the LogMath log domain
 */
func (a *HMMStateArc) LogProbability() float32 {
	return a.probability
}

/** returns a string representation of the arc */
func (a *HMMStateArc) String() string {
	return fmt.Sprintf("HSA %v prob %f", a.hmmState, a.probability)
}
