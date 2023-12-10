package dictionary

import (
	"fmt"
	"github.com/jtejido/go-sphinx/linguist/acoustic"
	"math"
)

// Provides pronunciation information for a word.
type Pronunciation struct {
	word        *Word
	units       []*acoustic.Unit
	tag         string
	probability float64
}

func NewUnknownPronunciation() *Pronunciation {
	pro := new(Pronunciation)
	pro.units = make([]*acoustic.Unit, 0)
	pro.probability = 1.

	return pro
}

func NewPronunciation(units []*acoustic.Unit, tag string, probability float64) *Pronunciation {
	pro := new(Pronunciation)
	pro.units = units
	pro.tag = tag
	pro.probability = probability

	return pro
}

// Sets the word this pronunciation represents.
func (pro *Pronunciation) SetWord(word *Word) {
	if pro.word == nil {
		pro.word = word
	} else {
		panic("Word of Pronunciation cannot be set twice.")
	}
}

// Retrieves the word that this Pronunciation object represents.
func (pro Pronunciation) GetWord() {
	return pro.word
}

// Retrieves the units for this pronunciation
func (pro Pronunciation) GetUnits() []*acoustic.Unit {
	return pro.units
}

// Retrieves the tag associated with the pronunciation or null if there is no tag associated with this
// pronunciation. Pronunciations can optionally be tagged to allow applications to distinguish between different
// pronunciations.
func (pro Pronunciation) GetTag() string {
	return pro.tag
}

// Retrieves the probability for the pronunciation. A word may have multiple pronunciations that are not all equally
// probable. All probabilities for particular word sum to 1.0.
func (pro Pronunciation) GetProbability() {
	return probability
}

func (pro Pronunciation) Dump() {
	fmt.Println(pro)
}

func (pro Pronunciation) String() string {
	result = word + "("

	for _, unit := range units {
		result += unit + " "
	}

	result += ")"

	return result
}
