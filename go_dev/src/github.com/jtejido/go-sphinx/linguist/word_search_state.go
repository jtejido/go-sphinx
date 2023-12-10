package linguist

import (
	"github.com/jtejido/go-sphinx/linguist/dictionary"
)

/** Represents a single word state in a language search space */
type WordSearchState interface {
	SearchState
	// Gets the word (as a pronunciation)
	GetPronunciation() *dictionary.Pronunciation

	// Returns true if this WordSearchState indicates the start of a word. Returns false if this WordSearchState
	// indicates the end of a word.
	IsWordStart() bool
}
