package api

import (
	"github.com/jtejido/go-sphinx/result"
	"io"
)

// High-level wrapper for Result instance.
type SpeechResult struct {
	result  *result.Result
	lattice *result.Lattice
}

// Constructs recognition result based on Result object.
// Accepts recognition result returned by Recognizer.
func NewSpeechResult(result *result.Result) *SpeechResult {
	sr := new(SpeechResult)
	sr.result = result
	if sr.result.ToCreateLattice() {
		sr.lattice = NewLattice(result)
		result.LatticeOptimizer{lattice}.Optimize()
		sr.lattice.ComputeNodePosteriors(1.0)
	} else {
		sr.lattice = nil
	}

	return sr
}

// Returns slice of words of the recognition result.
// Within the list words are ordered by time frame.
func (sr *SpeechResult) GetWords() []*result.WordResult {
	if sr.lattice != nil {
		return sr.lattice.GetWordResultPath()
	}

	return sr.result.GetTimedBestResult(false)
}

// Returns string representation of the result.
func (sr *SpeechResult) GetHypothesis() io.Reader {
	return sr.result.GetBestResultNoFiller()
}

// Returns lattice for the recognition result.
func (sr *SpeechResult) GetLattice() *result.Lattice {
	return sr.lattice
}

// Return Result object of current SpeechResult
func (sr *SpeechResult) GetResult() *result.Result {
	return sr.result
}
