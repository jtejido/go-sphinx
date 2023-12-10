package search

import (
	"math"
)

const (
	DEFAULT_MAX_DEPTH = 50
)

//  Partitions a list of tokens according to the token score, used
//  in PartitionActiveListFactory.
type Partitioner interface {
	// Partitions the given array of tokens in place, so that the highest scoring n token will be at the beginning of
	// the array, not in any order.
	Partition(tokens []*Token, size, n int)
}

type DefaultPartitioner struct {
}

// Partitions sub-array of tokens around the end token.
// Put all elements less or equal then pivot to the start of the array,
// shifting new pivot position
func (part DefaultPartitioner) endPointPartition(tokens []*Token, start, end int) int {
	pivot := tokens[end]
	pivotScore := pivot.GetScore()

	i := start
	j := end - 1

	for {

		for i < end && tokens[i].GetScore() >= pivotScore {
			i++
		}
		for j > i && tokens[j].GetScore() < pivotScore {
			j--
		}

		if j <= i {
			break
		}

		current := tokens[j]
		tokens[j] = tokens[i]
		tokens[i] = current

	}

	tokens[end] = tokens[i]
	tokens[i] = pivot

	return i
}

// Partitions sub-array of tokens around the x-th token by selecting the midpoint of the token array as the pivot.
// Partially solves issues with slow performance on already sorted arrays.
func (part DefaultPartitioner) midPointPartition(tokens []*Token, start, end int) int {
	middle := (start + end) >> 1
	temp := tokens[end]
	tokens[end] = tokens[middle]
	tokens[middle] = temp

	return part.endPointPartition(tokens, start, end)
}

func (part DefaultPartitioner) Partition(tokens []*Token, size, n int) int {
	if len(tokens) > n {
		return part.midPointSelect(tokens, 0, size-1, n, 0)
	} else {
		return part.findBest(tokens, size)
	}
}

// Simply find the best token and put it in the last slot
func (part DefaultPartitioner) findBest(tokens []*Token, size int) int {
	r := -1
	lowestScore := math.MaxFloat64
	for i := 0; i < len(tokens); i++ {
		currentScore := tokens[i].GetScore()
		if currentScore <= lowestScore {
			lowestScore = currentScore
			r = i // "r" is the returned index
		}
	}

	// exchange tokens[r] <=> last token,
	// where tokens[r] has the lowest score
	last := size - 1
	if last >= 0 {
		lastToken := tokens[last]
		tokens[last] = tokens[r]
		tokens[r] = lastToken
	}

	// return the last index
	return last
}

// Selects the token with the ith largest token score.
func (part DefaultPartitioner) midPointSelect(tokens []*Token, start, end, targetSize, depth int) int {
	if depth > DEFAULT_MAX_DEPTH {
		return part.simplePointSelect(tokens, start, end, targetSize)
	}

	if start == end {
		return start
	}

	partitionToken := part.midPointPartition(tokens, start, end)
	newSize := partitionToken - start + 1

	if targetSize == newSize {
		return partitionToken
	}

	if targetSize < newSize {
		return part.midPointSelect(tokens, start, partitionToken-1, targetSize, depth+1)
	}

	return part.midPointSelect(tokens, partitionToken+1, end, targetSize-newSize, depth+1)
}

// Fallback method to get the partition
func (part DefaultPartitioner) simplePointSelect(tokens []*Token, start, end, targetSize int) int {
	sort.Sort(scorer.ByScore(tokens[start : end+1]))
	return start + targetSize - 1
}
