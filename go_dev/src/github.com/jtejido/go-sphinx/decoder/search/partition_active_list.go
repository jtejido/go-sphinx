package search

import (
	"math"
)

type PartitionActiveList struct {
	size, absoluteBeamWidth int
	logRelativeBeamWidth    float64
	bestToken               *Token
	tokenList               []*Token
	partitioner             Partitioner
}

func NewPartitionActiveList(absoluteBeamWidth int, logRelativeBeamWidth float64) *PartitionActiveList {
	pal := new(PartitionActiveList)
	pal.absoluteBeamWidth = absoluteBeamWidth
	pal.logRelativeBeamWidth = logRelativeBeamWidth
	listSize := 2000
	if pal.absoluteBeamWidth > 0 {
		listSize = pal.absoluteBeamWidth / 3
	}
	pal.tokenList = make([]*Token, listSize)
	pal.partitioner = DefaultPartitioner{}
	return pal
}

func (pal *PartitionActiveList) Add(token *Token) {
	if pal.size < len(pal.tokenList) {
		pal.tokenList[pal.size] = token
		pal.size++
	} else {
		// token array too small, double the capacity
		pal.doubleCapacity()
		pal.Add(token)
	}
	if pal.bestToken == nil || token.GetScore() > bestToken.GetScore() {
		pal.bestToken = token
	}
}

func (pal *PartitionActiveList) doubleCapacity() {
	pal.tokenList = make([]*Token, len(pal.tokenList)*2)
}

// Purges excess members. Remove all nodes that fall below the relativeBeamWidth
func (pal *PartitionActiveList) Purge() *PartitionActiveList {
	if pal.absoluteBeamWidth > 0 {
		// if we have an absolute beam, then we will
		// need to sort the tokens to apply the beam
		if pal.size > pal.absoluteBeamWidth {
			pal.size = pal.partitioner.Partition(pal.tokenList, pal.size, pal.absoluteBeamWidth) + 1
		}
	}

	return pal

}

// Gets the set of all tokens
func (pal PartitionActiveList) GetTokens() []*Token {
	return pal.tokenList
}

// Returns the number of tokens on this active list
func (pal PartitionActiveList) Size() int {
	return pal.size
}

// Gets the beam threshold best upon the best scoring token
func (pal PartitionActiveList) GetBeamThreshold() float64 {
	return pal.GetBestScore() + pal.logRelativeBeamWidth
}

// Gets the best score in the list
func (pal PartitionActiveList) GetBestScore() float64 {
	bestScore := -math.MaxFloat64
	if pal.bestToken != nil {
		bestScore = pal.bestToken.GetScore()
	}
	// A sanity check
	// for (Token t : this) {
	//    if (t.getScore() > bestScore) {
	//         System.out.println("GBS: found better score "
	//             + t + " vs. " + bestScore);
	//    }
	// }
	return bestScore
}

// Sets the best scoring token for this active list
func (pal *PartitionActiveList) SetBestToken(token *Token) {
	pal.bestToken = token
}

// Gets the best scoring token for this active list
func (pal *PartitionActiveList) GetBestToken() *Token {
	return pal.bestToken
}

// Creates new instance with same properties
func (pal *PartitionActiveList) NewInstance() *PartitionActiveList {
	return NewPartitionActiveList(pal.absoluteBeamWidth, pal.logRelativeBeamWidth)
}
