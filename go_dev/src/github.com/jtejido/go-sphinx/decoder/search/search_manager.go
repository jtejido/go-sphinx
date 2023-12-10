package search

import (
	"github.com/jtejido/go-sphinx/result"
)

// Defines the interface for the SearchManager. The SearchManager's primary role is to execute the search for a given
// number of frames. The SearchManager will return interim results as the recognition proceeds and when recognition
// completes a final result will be returned.
type SearchManager interface {

	// Allocates the resources necessary for this search. This should be called once before an recognitions are
	// performed
	Allocate()

	// Deallocates resources necessary for this search. This should be called once after all recognitions are completed
	// at the search manager is no longer needed.
	Deallocate()

	// Prepares the SearchManager for recognition.  This method must be called before recognize is
	// called. Typically, start and stop are called bracketing an utterance.
	StartRecognition()

	// Performs post-recognition cleanup. This method should be called after recognize returns a final result.
	StopRecognition()

	// Performs recognition. Processes no more than the given number of frames before returning. This method returns a
	// partial result after nFrames have been processed, or a final result if recognition completes while processing
	// frames.  If a final result is returned, the actual number of frames processed can be retrieved from the result.
	// This method may block while waiting for frames to arrive.
	//
	// Accepts the maximum number of frames to process. A final result may be returned before all nFrames are
	// processed.
	//
	// Returns the recognition result, the result may be a partial or a final result; or return null if no frames are
	// arrived.
	Recognize(int) *result.Result
}
