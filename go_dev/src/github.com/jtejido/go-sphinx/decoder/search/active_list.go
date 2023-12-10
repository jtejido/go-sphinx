package search

// An active list is maintained as a sorted list
// Note that all scores are represented in LogMath logbase
type ActiveList interface {

	// Adds the given token to the list, keeping track of the lowest scoring token
	Add(*Token)

	// Purges the active list of excess members returning a (potentially new) active list
	Purge() ActiveList

	// Returns the size of this list
	Size() int

	// Gets the list of all tokens
	GetTokens() []*Token

	// Gets the beam threshold best upon the best scoring token
	GetBeamThreshold() float64

	// Gets the best score in the list
	GetBestScore() float64

	// Sets the best scoring token for this active list
	SetBestToken(*Token)

	// Gets the best scoring token for this active list
	GetBestToken() *Token

	// Creates a new empty version of this active list with the same general properties.
	ActiveListFactory
}
