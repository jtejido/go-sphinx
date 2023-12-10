package lextree

import (
	"fmt"
	"github.com/jtejido/go-sphinx/linguist"
)

// The LexTreeLinguist returns language states to the search manager. This class forms the base implementation for
// all language states returned. This LexTreeState keeps track of the probability of entering this state (a
// language+insertion probability) as well as the unit history. The unit history consists of the LexTree nodes that
// correspond to the left, center and right contexts.
type LexTreeState struct {
	node                               Node
	wordSequence                       *linguist.WordSequence
	currentSmearTerm, currentSmearProb float64
}

func NewLexTreeState(node Node, wordSequence *linguist.WordSequence, smearTerm, smearProb float64) *LexTreeState {
	lts := new(LexTreeState)
	lts.node = node
	lts.wordSequence = wordSequence
	lts.currentSmearProb = smearProb
	lts.currentSmearTerm = smearTerm

	return lts
}

func (lts *LexTreeState) GetSignature() string {
	return fmt.Sprintf("lts-%d-ws-%d", lts.node.HashCode(), wordSequence)
}

func (lts *LexTreeState) GetSmearTerm() float64 {
	return lts.currentSmearTerm
}

func (lts *LexTreeState) GetSmearProb() float64 {
	return lts.currentSmearProb
}

func (lts *LexTreeState) HashCode() int {
	hashCode := lts.wordSequence.HashCode() * 37
	hashCode += lts.node.HashCode()
	return hashCode
}

func (lts *LexTreeState) GetState() linguist.SearchState {
	return lts
}

func (lts *LexTreeState) GetProbability() float64 {
	return lts.getLanguageProbability() + lts.getInsertionProbability()
}

func (lts *LexTreeState) GetLanguageProbability() float64 {
	return lts.logOne
}

func (lts *LexTreeState) GetInsertionProbability() float64 {
	return lts.logOne
}
