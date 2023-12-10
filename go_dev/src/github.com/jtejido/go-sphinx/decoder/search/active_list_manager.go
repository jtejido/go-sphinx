package search

const (
	// The property that specifies the absolute word beam width
	DEFAULT_ABSOLUTE_WORD_BEAM_WIDTH = 2000

	// The property that specifies the relative word beam width
	DEFAULT_RELATIVE_WORD_BEAM_WIDTH = 0.0
)

// defaults for ActiveListFactory
const (
	// property that sets the desired (or target) size for this active list.  This is sometimes referred to as the beam
	// size
	DEFAULT_ABSOLUTE_BEAM_WIDTH = -1
	// Property that sets the minimum score relative to the maximum score in the list for pruning.  Tokens with a score
	// less than relativeBeamWidth * maximumScore will be pruned from the list
	DEFAULT_RELATIVE_BEAM_WIDTH = 1e-80
	// Property that indicates whether or not the active list will implement 'strict pruning'.  When strict pruning is
	// enabled, the active list will not remove tokens from the active list until they have been completely scored.  If
	// strict pruning is not enabled, tokens can be removed from the active list based upon their entry scores. The
	DEFAULT_STRING_PRUNING = true
)

type ActiveListManager interface {

	// Adds the given token to the list
	Add(*Token)

	// Returns an Iterator of all the non-emitting ActiveLists. The iteration order is the same as the search state
	// order.
	GetNonEmittingListIterator() ActiveListIterator

	// Returns the emitting ActiveList from the manager
	GetEmittingList() ActiveList

	// Clears emitting list in manager
	ClearEmittingList()

	// Dumps out debug info for the active list manager
	Dump()

	// Sets the total number of state types to be managed
	SetNumStateOrder(int)
}

// internal use, need public?
type ActiveListIterator interface {
	HasNext() bool
	Next() ActiveList
	Remove()
}
