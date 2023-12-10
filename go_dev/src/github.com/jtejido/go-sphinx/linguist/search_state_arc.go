package linguist

// Represents a single state in a language search space
type SearchStateArc interface {
	// Gets a successor to this search state
	GetState() SearchState

	// Gets the composite probability of entering this state
	GetProbability() float64

	// Gets the language probability of entering this state
	GetLanguageProbability() float64

	// Gets the insertion probability of entering this state
	GetInsertionProbability() float64
}
