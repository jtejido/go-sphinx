package tiedstate

import "github.com/jtejido/go-sphinx/frontend"

/** Represents a set of acoustic data that can be scored against a feature */
type Senone interface {

	/**
	 * Calculates the score for this senone based upon the given feature.
	 *
	 * @param feature the feature vector to score this senone against
	 * @return the score for this senone in LogMath log base
	 */
	Score(feature frontend.Data) float32

	/**
	 * Calculates the component scores for the mixture components in this senone based upon the given feature.
	 *
	 * @param feature the feature vector to score this senone against
	 * @return the scores for this senone in LogMath log base
	 */
	CalculateComponentScore(feature frontend.Data) []float32

	/**
	 * Gets the ID for this senone
	 *
	 * @return the senone id
	 */
	ID() int64

	/**
	 * Dumps a senone
	 *
	 * @param msg an annotation for the dump
	 */
	Dump(msg string)

	/**
	 * Returns the mixture components associated with this Gaussian
	 *
	 * @return the array of mixture components
	 */
	MixtureComponents() []*MixtureComponent

	/**
	 *
	 * @return the mixture weights vector
	 */
	LogMixtureWeights() []float32
}
