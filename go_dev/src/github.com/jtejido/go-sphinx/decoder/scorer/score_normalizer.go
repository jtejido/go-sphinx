package scorer

// Describes all API-elements that are necessary  to normalize token-scores after these have been computed by an
// AcousticScorer.
type ScoreNormalizer interface {

	// Normalizes the scores of a set of Tokens.
	Normalize([]Scoreable, Scoreable)
}
