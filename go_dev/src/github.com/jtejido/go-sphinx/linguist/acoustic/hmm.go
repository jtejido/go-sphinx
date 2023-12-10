package acoustic

/**
 * Represents a hidden-markov-model. An HMM consists of a unit (context dependent or independent), a transition matrix
 * from state to state, and a sequence of senones associated with each state. This representation of an HMM is a
 * specialized left-to-right markov model. No backward transitions are allowed.
 */

type HMM interface {

	/**
	 * Gets the  unit associated with this HMM
	 *
	 * @return the unit associated with this HMM
	 */
	Unit() *Unit

	/**
	 * Gets the  base unit associated with this HMM
	 *
	 * @return the unit associated with this HMM
	 */
	BaseUnit() *Unit

	/**
	 * @param which the state of interest
	 * @return hmm state
	 */
	State(which int) HMMState

	/**
	 * Returns the order of the HMM
	 *
	 * @return the order of the HMM
	 */
	Order() int

	/**
	 * Retrieves the position of this HMM.
	 *
	 * @return the position for this HMM
	 */
	Position() HMMPosition

	/**
	 * Gets the initial states (with probabilities) for this HMM
	 *
	 * @return the set of arcs that transition to the initial states for this HMM
	 */
	InitialState() HMMState
}
