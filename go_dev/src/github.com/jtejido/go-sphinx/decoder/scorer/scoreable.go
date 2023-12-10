package scorer

import (
	"github.com/jtejido/go-sphinx/frontend"
)

// utils for sorting Scoreable
type ByScore []Scoreable

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].GetScore() < a[j].GetScore() }

// Represents an entity that can be scored against a data
type Scoreable interface {
	frontend.Data
	/**
	 * Calculates a score against the given data. The score can be retrieved with get score
	 */
	CalculateScore(frontend.Data) float32

	/**
	 * Retrieves a previously calculated (and possibly normalized) score
	 */
	GetScore() float32

	/**
	 * Normalizes a previously calculated score
	 */
	NormalizeScore(maxScore float32) float64
}
