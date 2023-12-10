package dictionary

import (
	"math"
)

// Represents a word, its spelling and its pronunciation.
type Word struct {
	// the spelling of the word
	spelling string

	// pronunciations of this word
	pronunciations []*Pronunciation
	isFiller       bool
}

// The Word representing the unknown word.
func NewUnknownWord() *Word {
	w := new(Word)
	w.spelling = "<unk>"
	w.pronunciations = NewUnknownPronunciation()

	return w
}

func NewWord(spelling string, pronunciations []*Pronunciation, isFiller bool) *Word {
	w := new(Word)
	w.spelling = spelling
	w.pronunciations = pronunciations
	w.isFiller = isFiller
	return w
}

// Returns the spelling of the word.
func (w Word) GetSpelling() string {
	return w.spelling
}

// Determines if this is a filler word
func (w Word) IsFiller() bool {
	return w.isFiller
}

// Returns true if this word is an end of sentence word
func (w Word) IsSentenceEndWord() bool {
	return SENTENCE_END_SPELLING == w.spelling
}

// Returns true if this word is a start of sentence word
func (w Word) IsSentenceStartWord() bool {
	return SENTENCE_START_SPELLING == w.spelling
}

// Retrieves the pronunciations of this word
func (w Word) GetPronunciations() []*Pronunciation {
	return w.pronunciations
}

// Get the highest probability pronunciation for a word
func (w Word) GetMostLikelyPronunciation() *Pronunciation {
	bestScore := math.Inf(-1)
	var best *Pronunciation

	for _, pronunciation := range w.pronunciations {
		prob := pronunciation.GetProbability()
		if prob > bestScore {
			bestScore = prob
			best = pronunciation
		}
	}

	return best
}

func (w Word) HashCode() int {
	return w.spelling.HashCode()
}

func (w Word) Equals(obj *Word) bool {
	return w.spelling == obj.spelling
}

func (w Word) String() string {
	return w.spelling
}
