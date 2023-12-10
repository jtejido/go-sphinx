package api

import (
	"github.com/jtejido/go-sphinx/decoder/adaptation"
	"github.com/jtejido/go-sphinx/recognizer"
)

// Base struct for high-level speech recognizers.
type BaseSpeechRecognizer struct {
	context              *Context
	recognizer           *recognizer.Recognizer
	clusters             *adaptation.ClusteredDensityFileData
	speechSourceProvider *SpeechSourceProvider
}

// Returns result of the recognition.
//
// recognition result or nil if there is no result, e.g., because the microphone or input stream has been closed
func (br *BaseSpeechRecognizer) GetResult() *SpeechResult {
	result := br.recognizer.Recognize(nil)

	if result == nil {
		return nil
	}

	return NewSpeechResult(result)

}

func (br *BaseSpeechRecognizer) CreateStats(numClasses int) *adaptation.Stats {
	br.clusters = adaptation.NewClusteredDensityFileData(br.context.GetLoader(), numClasses)
	return adaptation.NewStats(br.context.GetLoader(), br.clusters)
}

func (br *BaseSpeechRecognizer) SetTransform(transform *adaptation.Transform) {
	if br.clusters != nil && transform != nil {
		br.context.GetLoader().Update(transform, br.clusters)
	}
}

func (br *BaseSpeechRecognizer) LoadTransform(path string, numClass int) error {
	br.clusters = adaptation.NewClusteredDensityFileData(br.context.GetLoader(), numClass)
	transform := adaptation.NewTransform(br.context.GetLoader(), numClass)
	if err := transform.Load(path); err != nil {
		return err
	}
	br.context.GetLoader().Update(transform, br.clusters)
	return nil
}
