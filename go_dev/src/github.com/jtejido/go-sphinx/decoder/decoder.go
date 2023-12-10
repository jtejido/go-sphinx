package decoder

import (
	"io"
	"math"

	"github.com/jtejido/go-sphinx/decoder/search"
	"github.com/jtejido/go-sphinx/result"
)

const (
	DEFAULT_FEATURE_BLOCK_SIZE     = math.MaxInt32
	DEFAULT_FIRE_NON_FINAL_RESULTS = false
	DEFAULT_AUTO_ALLOCATE          = false
)

type Decoder struct {
	BaseDecoder
	featureBlockSize int
}

func NewDefaultDecoder() *Decoder {
	// this doesn't autoallocate
	d := new(Decoder)
	d.searchManager = search.NewDefaultWordPruningBreadthFirstLookaheadSearchManager()
	d.fireNonFinalResults = DEFAULT_FIRE_NON_FINAL_RESULTS

	d.featureBlockSize = DEFAULT_FEATURE_BLOCK_SIZE

	return d
}

func NewDecoder(searchManager search.SearchManager, fireNonFinalResults, autoAllocate bool, featureBlockSize int) *Decoder {
	d := new(Decoder)
	d.searchManager = searchManager
	d.fireNonFinalResults = fireNonFinalResults

	if autoAllocate {
		d.searchManager.Allocate()
	}

	d.featureBlockSize = featureBlockSize

	return d

}

func (d *Decoder) Decode(referenceText io.Reader) (result *result.Result) {
	d.searchManager.StartRecognition()
	for {
		result = d.searchManager.Recognize(d.featureBlockSize)
		if result != nil {
			result.SetReferenceText(referenceText)
			d.fireResultListeners(result)
		}

		if result == nil || result.IsFinal() {
			break
		}
	}
	d.searchManager.StopRecognition()
	return
}
