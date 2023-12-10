package linguist

import (
	"github.com/jtejido/go-sphinx/linguist/dictionary"
	"math"
)

// This class can be used to keep track of a word sequence.
type WordSequence struct {
	words    []*dictionary.Word // making this private ensures immutability
	hashCode int64
}

func NewEmptyWordSequence() *WordSequence {
	ws := new(WordSequence)
	ws.words = make([]*dictionary.Word, 0)
	ws.hashCode = -1

	return ws
}

func NewWordSequenceBySize(size int) *WordSequence {
	ws := new(WordSequence)
	ws.words = make([]*dictionary.Word, size)
	ws.hashCode = -1

	return ws
}

func NewWordSequenceByWordSlice(words []*dictionary.Word) *WordSequence {
	ws := new(WordSequence)
	ws.words = words
	ws.hashCode = -1
	ws.check()
	return ws
}

func (ws WordSequence) check() {
	for _, word := range ws.words {
		if word == nil {
			panic("WordSequence should not have null Words.")
		}
	}
}

// Returns a new word sequence with the given word added to the sequence
func (ws *WordSequence) AddWord(word *dictionary.Word, maxSize int) *WordSequence {

	var nextSize int

	if maxSize <= 0 {
		return NewEmptyWordSequence()
	}

	if (ws.Size() + 1) > maxSize {

		nextSize = maxSize
	} else {
		nextSize = ws.Size() + 1
	}

	next := NewWordSequenceBySize(nextSize)
	nextIndex := nextSize - 1
	thisIndex := ws.Size() - 1
	next.words[nextIndex] = word
	nextIndex--

	for nextIndex >= 0 && thisIndex >= 0 {
		next.words[nextIndex] = this.words[thisIndex]
		nextIndex--
		thisIndex--
	}

	next.check()

	return next
}

// Returns the oldest words in the sequence (the newest word is omitted)
func (ws WordSequence) GetOldest() *WordSequence {
	next := NewEmptyWordSequence()

	if ws.Size() >= 1 {
		next = NewWordSequenceBySize(len(ws.words) - 1)
		copy(next.words, ws.words[0:])
	}

	return next
}

// Returns the newest words in the sequence (the old word is omitted)
func (ws WordSequence) GetNewest() *WordSequence {
	next := NewEmptyWordSequence()

	if ws.Size() >= 1 {
		next = NewWordSequenceBySize(len(ws.words) - 1)
		copy(next.words, ws.words[1:])
	}

	return next
}

// Returns a word sequence that is no longer than the given size, that is
// filled in with the newest words from this sequence
func (ws *WordSequence) Trim(maxSize int) *WordSequence {
	if maxSize <= 0 || ws.Size() == 0 {
		return NewEmptyWordSequence()
	} else if maxSize == ws.Size() {
		return ws
	}

	if maxSize > ws.Size() {
		maxSize = ws.Size()
	}

	next := NewWordSequenceBySize(maxSize)
	thisIndex := len(ws.words) - 1
	nextIndex := len(next.words) - 1

	for i := 0; i < maxSize; i++ {
		next.words[nextIndex] = this.words[thisIndex]
		nextIndex--
		thisIndex--
	}

	return next
}

// Returns the n-th word in this sequence
func (ws WordSequence) GetWord(n int) *dictionary.Word {
	if n > len(ws.words) {
		panic("n greater than number of words")
	}

	return ws.words[n]
}

// Returns the number of words in this sequence
func (ws WordSequence) Size() int {
	return len(ws.words)
}

// Calculates the hashcode for this object
func (ws *WordSequence) HashCode() int {
	if ws.hashCode == -1 {
		code := 123
		for i := 0; i < len(ws.words); i++ {
			code += ws.words[i].HashCode() * (2*i + 1)
		}
		ws.hashCode = code
	}
	return ws.hashCode
}

//  Returns a subsequence with both startIndex and stopIndex exclusive.
func (ws WordSequence) GetSubSequence(startIndex, stopIndex int) *WordSequence {

	subseqWords := make([]*dictionary.Word, 0)

	for i := startIndex; i < stopIndex; i++ {
		subseqWords = append(subseqWords, ws.GetWord(i))
	}

	return NewWordSequenceByWordSlice(subseqWords)
}

// Returns the words of the WordSequence
func (ws WordSequence) GetWords() []*dictionary.Word {
	return ws.words
}

func (ws WordSequence) Equals(other *WordSequence) bool {
	if len(ws.words) != len(other.words) {
		return false
	}

	for i := 0; i < len(ws.words); i++ {
		if !words[i].Equals(other.words[i]) {
			return false
		}
	}

	return true
}
