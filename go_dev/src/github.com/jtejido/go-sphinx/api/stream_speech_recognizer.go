package api

import (
	"io"

	"github.com/jtejido/go-sphinx/recognizer"
	"github.com/jtejido/go-sphinx/util"
)

// Speech recognizer that works with audio resources.
type StreamSpeechRecognizer struct {
	BaseSpeechRecognizer
}

// Constructs new stream recognizer.
func NewDefaultStreamSpeechRecognizer(configuration *Configuration) *StreamSpeechRecognizer {
	ssr := new(StreamSpeechRecognizer)
	ssr.context = NewDefaultContext(configuration)
	// ssr.Recognizer = ssr.Context.GetInstance("recognizer").(recognizer.Recognizer)

	ssr.recognizer = recognizer.NewDefaultRecognizer() // default hell starts here
	ssr.speechSourceProvider = &SpeechSourceProvider{}
	return ssr
}

func (ssr *StreamSpeechRecognizer) StartRecognition(stream io.Reader) {
	ssr.recognizer.Allocate()
	ssr.context.SetSpeechSource(stream, util.INFINITE)
}

// Starts recognition process.
//
// Starts recognition process and optionally clears previous data.
func (ssr *StreamSpeechRecognizer) StartRecognitionLimit(stream io.Reader, timeFrame *util.TimeFrame) {
	ssr.recognizer.Allocate()
	ssr.context.SetSpeechSource(stream, timeFrame)
}

// Stops recognition process.
//
// Recognition process is paused until the next call to startRecognition.
func (ssr *StreamSpeechRecognizer) StopRecognition() {
	ssr.recognizer.Deallocate()
}
