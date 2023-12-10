package scorer

import (
	"github.com/jtejido/go-sphinx/frontend"
)

// Provides a mechanism for scoring a set of HMM states
type AcousticScorer interface {

	// Allocates resources for this scorer
	Allocate()

	// Deallocates resources for this scorer
	Deallocate()

	// starts the scorer
	StartRecognition()

	// stops the scorer
	StopRecognition()

	// Scores the given set of states over previously stored acoustic data if any or a new one
	CalculateScores([]Scoreable) frontend.Data

	// Scores the given set of states over previously acoustic data from frontend
	// and stores latter in the queue
	CalculateScoresAndStoreData([]Scoreable) frontend.Data
}
