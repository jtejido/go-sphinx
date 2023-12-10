package linguist

// Represents a search graph
type SearchGraph interface {
	// Retrieves initial search state
	GetInitialState() SearchState

	// Returns the number of different state types maintained in the search graph
	GetNumStateOrder() int

	// Order of words and data tokens
	GetWordTokenFirst() bool
}
