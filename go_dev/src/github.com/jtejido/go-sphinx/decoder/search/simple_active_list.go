package search

import (
	"github.com/jtejido/go-sphinx/decoder/scorer"
	"math"
	"sort"
)

const (
	DEFAULT_CHECK_PRIOR_LISTS_EMPTY = false
)

// An active list that tries to be simple and correct. This type of active list will be slow, but should exhibit
// correct behavior. Faster versions of the ActiveList exist (HeapActiveList, TreeActiveList).
// This class is not thread safe and should only be used by a single thread.
//
// Note that all scores are maintained in the LogMath log domain.
type SimpleActiveList struct {
	absoluteBeamWidth    int
	logRelativeBeamWidth float64
	bestToken            *Token
	tokenList            []*Token
}

func NewDefaultSimpleActiveList() *SimpleActiveList {
	sal := new(SimpleActiveList)
	sal.absoluteBeamWidth = 20000
	sal.logRelativeBeamWidth = 1e-60
	sal.tokenList = make([]*Token, 0)
	return sal
}

func NewSimpleActiveList(absoluteBeamWidth int, logRelativeBeamWidth float64) *SimpleActiveList {
	sal := new(SimpleActiveList)
	sal.absoluteBeamWidth = absoluteBeamWidth
	sal.logRelativeBeamWidth = logRelativeBeamWidth
	sal.tokenList = make([]*Token, 0)
	return sal
}

// Adds the given token to the list
func (sal *SimpleActiveList) Add(token *Token) {
	sal.tokenList = append(sal.tokenList, token)
	if sal.bestToken == nil || token.GetScore() > bestToken.GetScore() {
		sal.bestToken = token
	}
}

// Purges excess members. Remove all nodes that fall below the relativeBeamWidth
func (sal *SimpleActiveList) Purge() *SimpleActiveList {
	if sal.absoluteBeamWidth > 0 && len(sal.tokenList) > sal.absoluteBeamWidth {
		sort.Sort(scorer.ByScore(tokenList))
		copy(sal.tokenList[0:], sal.tokenList[sal.absoluteBeamWidth:])
		for k, n := len(sal.tokenList)-sal.absoluteBeamWidth, len(sal.tokenList); k < n; k++ {
			sal.tokenList[k] = nil
		}
		sal.tokenList = sal.tokenList[:len(sal.tokenList)-sal.absoluteBeamWidth]
	}

	return sal
}

// Gets the set of all tokens
func (sal SimpleActiveList) GetTokens() []*Token {
	return sal.tokenList
}

// Returns the number of tokens on this active list
func (sal SimpleActiveList) Size() int {
	return len(sal.tokenList)
}

// Gets the beam threshold best upon the best scoring token
func (sal SimpleActiveList) GetBeamThreshold() float64 {
	return sal.GetBestScore() + sal.logRelativeBeamWidth
}

// Gets the best score in the list
func (sal SimpleActiveList) GetBestScore() float64 {
	bestScore = -math.MaxFloat64
	if sal.bestToken != nil {
		bestScore = sal.bestToken.GetScore()
	}

	return bestScore
}

// Sets the best scoring token for this active list
func (sal *SimpleActiveList) SetBestToken(token *Token) {
	sal.bestToken = token
}

// Gets the best scoring token for this active list
func (sal *SimpleActiveList) GetBestToken() *Token {
	return sal.bestToken
}

// Creates new instance with same properties
func (sal *SimpleActiveList) NewInstance() *SimpleActiveList {
	return NewSimpleActiveList(sal.absoluteBeamWidth, sal.logRelativeBeamWidth)
}
