package linguist

// Represents a single state in a language search space
type SearchState interface {
	// Gets a successor to this search state
	GetSuccessors() []SearchStateArc

	// Determines if this is an emitting state
	IsEmitting() bool

	// Determines if this is a final state
	IsFinal() bool

	// Returns a pretty version of the string representation for this object
	ToPrettyString() string

	// Returns a unique signature for this state
	GetSignature() string

	// Gets the word history for this state
	GetWordHistory() *WordSequence

	// Returns the lex tree state
	GetLexState() interface{}

	// Returns the order of this particular state
	GetOrder() int
}
