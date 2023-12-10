package search

type BaseActiveListFactory struct {
	absoluteBeamWidth    int
	logRelativeBeamWidth float64
}

type ActiveListFactory interface {
	// Creates a new empty version of an active list with the same general properties.
	NewInstance() ActiveList
}
