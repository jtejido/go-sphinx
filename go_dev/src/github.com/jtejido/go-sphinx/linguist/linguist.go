package linguist

const (
	// Word insertion probability property
	DEFAULT_WORD_INSERTION_PROBABILITY = 1.0

	// Unit insertion probability property
	DEFAULT_UNIT_INSERTION_PROBABILITY = 1.0

	// Silence insertion probability property
	DEFAULT_SILENCE_INSERTION_PROBABILITY = 1.0

	// Filler insertion probability property
	DEFAULT_FILLER_INSERTION_PROBABILITY = 1.0

	// The property that defines the language weight for the search
	DEFAULT_LANGUAGE_WEIGHT = 1.0
)

type Linguist interface {
	// Retrieves search graph.  The search graph represents the search space to be used to guide the search.
	GetSearchGraph() SearchGraph

	// Called before a recognition. This method gives a linguist the opportunity to prepare itself before a recognition
	// begins.
	StartRecognition()

	// Called after a recognition. This method gives a linguist the opportunity to clean up after a recognition has been
	// completed.
	StopRecognition()

	// Allocates the linguist. Resources allocated by the linguist are allocated here. This method may take many seconds
	// to complete depending upon the linguist.
	Allocate()

	// Deallocates the linguist. Any resources allocated by this linguist are released.
	Deallocate()
}
