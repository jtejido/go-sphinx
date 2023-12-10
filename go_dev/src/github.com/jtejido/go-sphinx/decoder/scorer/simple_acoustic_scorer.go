package scorer

import (
	fe "github.com/jtejido/go-sphinx/frontend"
	"github.com/jtejido/go-sphinx/frontend/util"
	"math"
)

// Implements some basic scorer functionality, including a simple default
// acoustic scoring implementation which scores within the current thread, that
// can be changed by overriding the {@link #doScoring} method.
type SimpleAcousticScorer struct {

	// Property the defines the frontend to retrieve features from for scoring
	frontEnd *fe.FrontEnd

	// An optional post-processor for computed scores that will normalize
	// scores. If not set, no normalization will applied and the token scores
	// will be returned unchanged.
	scoreNormalizer ScoreNormalizer
	storedData      []fe.Data
	seenEnd         bool
}

func NewDefaultSimpleAcousticScorer() *SimpleAcousticScorer {
	sas := new(SimpleAcousticScorer)
	sas.frontEnd = fe.NewDefaultFrontend()
	sas.storedData = make([]fe.Data, 0)

	return sas
}

func NewSimpleAcousticScorer(frontEnd *fe.FrontEnd, scoreNormalizer ScoreNormalizer) *SimpleAcousticScorer {
	sas := new(SimpleAcousticScorer)
	sas.frontEnd = frontEnd
	sas.scoreNormalizer = scoreNormalizer
	sas.storedData = make([]fe.Data, 0)

	return sas
}

func (sas *SimpleAcousticScorer) CalculateScores(scoreableList []Scoreable) fe.Data {
	var data fe.Data

	if len(sas.storedData) == 0 {

		for {
			data = sas.getNextData()
			_, ok := data.(fe.Signal)

			if !ok {
				break
			}

			_, ses_ok := data.(fe.SpeechEndSignal)

			if ses_ok {
				seenEnd = true
				break
			}

			_, des_ok := data.(fe.SpeechEndSignal)

			if des_ok {
				if sas.seenEnd {
					return nil
				} else {
					break
				}
			}
		}
		if data == nil {
			return nil
		}
	} else {
		data = sas.storedData[0]
		// remove head
		copy(sas.storedData[0:], sas.storedData[1:])
		sas.storedData[len(sas.storedData)-1] = nil
		sas.storedData = sas.storedData[:len(wpbflsm.ciScores)-1]
	}

	return sas.calculateScoresForData(scoreableList, data)
}

func (sas *SimpleAcousticScorer) CalculateScoresAndStoreData(scoreableList []Scoreable) fe.Data {
	var data fe.Data

	for {
		data = sas.getNextData()
		_, ok := data.(fe.Signal)

		if !ok {
			break
		}

		_, ses_ok := data.(fe.SpeechEndSignal)

		if ses_ok {
			seenEnd = true
			break
		}

		_, des_ok := data.(fe.SpeechEndSignal)

		if des_ok {
			if sas.seenEnd {
				return nil
			} else {
				break
			}
		}

	}
	if data == nil {
		return nil
	}

	sas.storedData = append(sas.storedData, data)

	return sas.calculateScoresForData(scoreableList, data)
}

func (sas *SimpleAcousticScorer) calculateScoresForData(scoreableList []Scoreable, data fe.Data) fe.Data {
	var dutil util.DataUtil
	_, ses_ok := data.(fe.SpeechEndSignal)
	_, des_ok := data.(fe.SpeechEndSignal)

	if ses_ok || des_ok {
		return data
	}

	if len(scoreableList) == 0 {
		return nil
	}

	ddat, ddat_ok := data.(fe.DoubleData)
	// convert the data to FloatData if not yet done
	if ddat_ok {
		data = dutil.DoubleData2FloatData(ddat)
	}

	bestToken := sas.doScoring(scoreableList, data)

	// apply optional score normalization
	// assume it's a token
	if sas.scoreNormalizer != nil {
		bestToken = sas.scoreNormalizer.Normalize(scoreableList, bestToken)
	}

	return bestToken
}

func (sas *SimpleAcousticScorer) getNextData() fe.Data {
	return sas.frontEnd.GetData()
}

func (sas *SimpleAcousticScorer) StartRecognition() {
	sas.storedData.Clear()
}

func (sas *SimpleAcousticScorer) StopRecognition() {}

func (sas *SimpleAcousticScorer) doScoring(scoreableList []Scoreable, data fe.Data) Scoreable {

	var best Scoreable
	bestScore := -Float.MaxFloat64

	for _, item := range scoreableList {
		item.CalculateScore(data)
		if item.GetScore() > bestScore {
			bestScore = item.GetScore()
			best = item
		}
	}

	return best
}

// Even if we don't do any meaningful allocation here, we implement the
// methods because most extending scorers do need them either.
func (sas *SimpleAcousticScorer) Allocate() {}

func (sas *SimpleAcousticScorer) Deallocate() {}
