package pruner

import (
	"github.com/jtejido/go-sphinx/decoder/search"
)

// Provides a mechanism for pruning a set of StateTokens
type Pruner interface {
	// Starts the pruner
	StartRecognition()

	// prunes the given set of states
	Prune(search.ActiveList) search.ActiveList

	// Performs post-recognition cleanup.
	StopRecognition()

	// Allocates resources necessary for this pruner
	Allocate()

	// Deallocates resources necessary for this pruner
	Deallocate()
}
